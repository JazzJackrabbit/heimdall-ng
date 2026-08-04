package providers_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/JazzJackrabbit/heimdall-ng/models"
	"github.com/JazzJackrabbit/heimdall-ng/providers"
	"github.com/JazzJackrabbit/heimdall-ng/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// imageCaptureTransport records the outgoing image request and returns a canned
// gpt-image-1 response, so payload construction is testable offline.
type imageCaptureTransport struct {
	path        string
	contentType string
	body        []byte
	/** Overrides the default response body when set. */
	canned string
}

func (c *imageCaptureTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	c.path = r.URL.Path
	c.contentType = r.Header.Get("Content-Type")
	c.body, _ = io.ReadAll(r.Body)

	canned := c.canned
	if canned == "" {
		canned = `{"created": 1, "data": [{"b64_json": "aGVsbG8="}],
			"usage": {"total_tokens": 1583, "input_tokens": 55, "output_tokens": 1528}}`
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(canned)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

// parseMultipart splits a captured multipart body into its form fields and the
// uploaded file parts.
func parseMultipart(t *testing.T, contentType string, body []byte) (map[string]string, map[string][]byte) {
	t.Helper()

	_, params, err := mime.ParseMediaType(contentType)
	require.NoError(t, err, "content type should be parseable")

	reader := multipart.NewReader(strings.NewReader(string(body)), params["boundary"])
	fields := map[string]string{}
	files := map[string][]byte{}

	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err, "reading multipart part")

		content, err := io.ReadAll(part)
		require.NoError(t, err, "reading part content")

		if part.FileName() != "" {
			files[part.FileName()] = content
		} else {
			fields[part.FormName()] = string(content)
		}
	}

	return fields, files
}

// A gpt-image-1 model carrying reference images has to reach /images/edits.
// /images/generations accepts a prompt only, so routing an attachment-bearing
// request there drops the references without an error.
func TestGPTImageWithReferencesUsesEditEndpoint(t *testing.T) {
	t.Parallel()

	transport := &imageCaptureTransport{}
	client := http.Client{Transport: transport}
	openai := providers.NewOpenAI([]string{"test-key"})

	res, err := openai.CompleteResponse(context.Background(), request.Completion{
		Model: &models.GPTImage{
			ImageFile: []models.OpenaiImagePayload{
				{Url: "data:image/png;base64,aGVsbG8="},
			},
		},
		SystemMessage: "Redraw the attached screen in a new design language.",
		UserMessage:   "Make it feel calm and premium",
		Temperature:   1,
	}, client, nil)
	require.NoError(t, err, "CompleteResponse returned an unexpected error")
	assert.Equal(t, "aGVsbG8=", res.Content, "content should be the image base64")

	assert.Equal(t, "/v1/images/edits", transport.path, "reference images must route to the edit endpoint")
	assert.True(t, strings.HasPrefix(transport.contentType, "multipart/form-data"),
		"the edit endpoint takes multipart form data, got %q", transport.contentType)

	fields, files := parseMultipart(t, transport.contentType, transport.body)

	assert.Equal(t, models.ImageModelAlias, fields["model"], "model field")
	assert.Equal(t, "Redraw the attached screen in a new design language.\n\nMake it feel calm and premium",
		fields["prompt"], "the system message must be folded into the prompt")

	require.Len(t, files, 1, "one reference image should be uploaded")
	assert.Equal(t, []byte("hello"), files["reference_1.png"], "the decoded image bytes should be uploaded")
}

// Without reference images the call stays on the JSON generation endpoint.
func TestGPTImageWithoutReferencesUsesGenerationEndpoint(t *testing.T) {
	t.Parallel()

	transport := &imageCaptureTransport{}
	client := http.Client{Transport: transport}
	openai := providers.NewOpenAI([]string{"test-key"})

	_, err := openai.CompleteResponse(context.Background(), request.Completion{
		Model:       &models.GPTImage{},
		UserMessage: "A red panda sipping bubble tea",
		Temperature: 1,
	}, client, nil)
	require.NoError(t, err, "CompleteResponse returned an unexpected error")

	assert.Equal(t, "/v1/images/generations", transport.path, "a prompt-only request should use the generation endpoint")
	assert.Equal(t, "application/json", transport.contentType, "the generation endpoint takes JSON")

	var payload map[string]any
	require.NoError(t, json.Unmarshal(transport.body, &payload), "request body should be JSON")
	assert.Equal(t, "A red panda sipping bubble tea", payload["prompt"], "prompt field")
	assert.InDelta(t, 1, payload["n"], 0, "n should be numeric on the JSON endpoint")

	// Unset size must not be sent: the API defaults to "auto" and picks an
	// aspect ratio for the prompt. Pinning a square here composed every
	// portrait subject into 1024x1024 and cropped it.
	assert.NotContains(t, payload, "size", "size must be omitted when the caller did not set one")
}

// A caller that does ask for a size gets exactly that, on both transports.
func TestGPTImageSendsAnExplicitSize(t *testing.T) {
	t.Parallel()

	transport := &imageCaptureTransport{}
	client := http.Client{Transport: transport}
	openai := providers.NewOpenAI([]string{"test-key"})

	_, err := openai.CompleteResponse(context.Background(), request.Completion{
		Model:       &models.GPTImage{Size: models.GPTImageSize1024x1536},
		UserMessage: "A book cover",
		Temperature: 1,
	}, client, nil)
	require.NoError(t, err, "CompleteResponse returned an unexpected error")

	var payload map[string]any
	require.NoError(t, json.Unmarshal(transport.body, &payload), "request body should be JSON")
	assert.Equal(t, "1024x1536", payload["size"], "an explicit size should be sent through")

	editTransport := &imageCaptureTransport{}
	editClient := http.Client{Transport: editTransport}
	_, err = openai.CompleteResponse(context.Background(), request.Completion{
		Model: &models.GPTImage{
			Size:      models.GPTImageSize1536x1024,
			ImageFile: []models.OpenaiImagePayload{{Url: "data:image/png;base64,aGVsbG8="}},
		},
		UserMessage: "Widen this",
		Temperature: 1,
	}, editClient, nil)
	require.NoError(t, err, "edit CompleteResponse returned an unexpected error")

	fields, _ := parseMultipart(t, editTransport.contentType, editTransport.body)
	assert.Equal(t, "1536x1024", fields["size"], "the edit endpoint should carry the size too")
}

// Every reference image is uploaded, and the extension follows the mime type so
// the API accepts the part.
func TestGPTImageUploadsAllReferencesWithMimeExtensions(t *testing.T) {
	t.Parallel()

	transport := &imageCaptureTransport{}
	client := http.Client{Transport: transport}
	openai := providers.NewOpenAI([]string{"test-key"})

	_, err := openai.CompleteResponse(context.Background(), request.Completion{
		Model: &models.GPTImage{
			ImageFile: []models.OpenaiImagePayload{
				{Url: "data:image/png;base64,aGVsbG8="},
				{Url: "data:image/jpeg;base64,d29ybGQ="},
				{Url: "data:image/webp;base64,IQ=="},
			},
		},
		UserMessage: "Combine these",
		Temperature: 1,
	}, client, nil)
	require.NoError(t, err, "CompleteResponse returned an unexpected error")

	_, files := parseMultipart(t, transport.contentType, transport.body)

	require.Len(t, files, 3, "all reference images should be uploaded")
	assert.Equal(t, []byte("hello"), files["reference_1.png"], "png reference")
	assert.Equal(t, []byte("world"), files["reference_2.jpg"], "jpeg reference")
	assert.Equal(t, []byte("!"), files["reference_3.webp"], "webp reference")
}

// The picture is billed as output tokens, so dropping the usage block made
// every image call read as free to anything metering spend.
func TestGPTImageReportsUsage(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"generation", "edit"} {
		model := &models.GPTImage{}
		if name == "edit" {
			model.ImageFile = []models.OpenaiImagePayload{{Url: "data:image/png;base64,aGVsbG8="}}
		}

		transport := &imageCaptureTransport{}
		client := http.Client{Transport: transport}
		openai := providers.NewOpenAI([]string{"test-key"})

		res, err := openai.CompleteResponse(context.Background(), request.Completion{
			Model:       model,
			UserMessage: "A lemon",
			Temperature: 1,
		}, client, nil)
		require.NoError(t, err, "%s: CompleteResponse returned an unexpected error", name)

		assert.Equal(t, 55, res.Usage.PromptTokens, "%s: input tokens", name)
		assert.Equal(t, 1528, res.Usage.CompletionTokens, "%s: output tokens — the image itself", name)
		assert.Equal(t, 1583, res.Usage.TotalTokens, "%s: total tokens", name)
	}
}

// A response without a usage block (DALL·E, or an older shape) still totals
// correctly rather than reporting a total of zero against non-zero parts.
func TestGPTImageUsageTotalsWhenAbsent(t *testing.T) {
	t.Parallel()

	transport := &imageCaptureTransport{
		canned: `{"created": 1, "data": [{"b64_json": "aGVsbG8="}],
			"usage": {"input_tokens": 10, "output_tokens": 90}}`,
	}
	client := http.Client{Transport: transport}
	openai := providers.NewOpenAI([]string{"test-key"})

	res, err := openai.CompleteResponse(context.Background(), request.Completion{
		Model:       &models.GPTImage{},
		UserMessage: "A lemon",
		Temperature: 1,
	}, client, nil)
	require.NoError(t, err, "CompleteResponse returned an unexpected error")
	assert.Equal(t, 100, res.Usage.TotalTokens, "total should fall back to input+output")
}

// A non-http(s) URL must be refused rather than dereferenced.
func TestGPTImageRejectsUnsupportedURLScheme(t *testing.T) {
	t.Parallel()

	transport := &imageCaptureTransport{}
	client := http.Client{Transport: transport}
	openai := providers.NewOpenAI([]string{"test-key"})

	_, err := openai.CompleteResponse(context.Background(), request.Completion{
		Model: &models.GPTImage{
			ImageFile: []models.OpenaiImagePayload{{Url: "file:///etc/passwd"}},
		},
		UserMessage: "Anything",
		Temperature: 1,
	}, client, nil)
	require.Error(t, err, "a file:// reference should be refused")
	assert.Contains(t, err.Error(), "unsupported image URL scheme")
}

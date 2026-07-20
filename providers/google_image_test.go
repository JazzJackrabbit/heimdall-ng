package providers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/JazzJackrabbit/heimdall/models"
	"github.com/JazzJackrabbit/heimdall/providers"
	"github.com/JazzJackrabbit/heimdall/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureTransport records the outgoing request body and returns a canned
// Gemini image response, so payload construction is testable offline.
type captureTransport struct {
	body []byte
}

func (c *captureTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	c.body, _ = io.ReadAll(r.Body)
	canned := `{
		"candidates": [{"content": {"parts": [{"inlineData": {"mimeType": "image/png", "data": "aGVsbG8="}}]}}],
		"usageMetadata": {"promptTokenCount": 10, "candidatesTokenCount": 20, "totalTokenCount": 30}
	}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(canned)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

// Gemini image models do not apply systemInstruction, so the system message
// must be folded into the prompt text part — otherwise a caller's brief is
// silently dropped.
func TestGoogleImageModelsFoldSystemMessageIntoPrompt(t *testing.T) {
	t.Parallel()

	for _, m := range []models.Model{
		&models.Gemini25FlashImage{},
		&models.Gemini3ProImagePreview{},
	} {
		transport := &captureTransport{}
		client := http.Client{Transport: transport}
		google := providers.NewGoogle([]string{"test-key"})

		res, err := google.CompleteResponse(context.Background(), request.Completion{
			Model:         m,
			SystemMessage: "Draw the subject as a flat vector sticker.",
			UserMessage:   "A red panda sipping bubble tea",
			Temperature:   1,
		}, client, nil)
		require.NoError(t, err, "%T: CompleteResponse returned an unexpected error", m)
		assert.Equal(t, "aGVsbG8=", res.Content, "%T: content should be the image base64", m)

		var payload map[string]any
		require.NoError(t, json.Unmarshal(transport.body, &payload), "%T: request body should be JSON", m)

		assert.NotContains(t, payload, "systemInstruction", "%T: image models must not send systemInstruction", m)

		text := string(transport.body)
		assert.Contains(t, text, "Draw the subject as a flat vector sticker.\\n\\nA red panda sipping bubble tea",
			"%T: prompt text part should carry the folded system message", m)
	}
}

// A system-prompt-only sheet (image sample with empty input) must still send a
// non-empty prompt.
func TestGoogleImageModelsSystemMessageOnly(t *testing.T) {
	t.Parallel()

	transport := &captureTransport{}
	client := http.Client{Transport: transport}
	google := providers.NewGoogle([]string{"test-key"})

	_, err := google.CompleteResponse(context.Background(), request.Completion{
		Model:         &models.Gemini25FlashImage{},
		SystemMessage: "Draw a lighthouse at dusk.",
		Temperature:   1,
	}, client, nil)
	require.NoError(t, err)
	assert.Contains(t, string(transport.body), "Draw a lighthouse at dusk.")
}

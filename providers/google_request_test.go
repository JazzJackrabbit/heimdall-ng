package providers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JazzJackrabbit/heimdall-ng/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureGoogleRequest points the Google provider at a test server that
// records the request path and JSON body, then replies with a minimal
// streamed completion.
func captureGoogleRequest(t *testing.T) (*string, *map[string]any) {
	t.Helper()

	var path string
	var body map[string]any

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			path = r.URL.Path
			if raw, err := io.ReadAll(r.Body); err == nil {
				_ = json.Unmarshal(raw, &body)
			}

			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`data: {"candidates":[{"content":{"parts":[{"text":"Hello"}],"role":"model"},"finishReason":"STOP"}],"usageMetadata":{"promptTokenCount":1,"candidatesTokenCount":1,"totalTokenCount":2}}` + "\n\n"))
		},
	))
	t.Cleanup(server.Close)

	original := googleBaseURL
	googleBaseURL = server.URL + "/%s:streamGenerateContent?alt=sse&key=%s"
	t.Cleanup(func() { googleBaseURL = original })

	return &path, &body
}

func TestGemini37FlashRequest(t *testing.T) {
	path, body := captureGoogleRequest(t)

	res, _ := streamAndCollect(
		t,
		NewGoogle([]string{"test-key"}),
		models.Gemini37Flash{ThinkingLevel: models.MediumThinkingLevel},
	)

	assert.Equal(t, "Hello", res.Content)
	assert.Equal(t, "/gemini-3.7-flash:streamGenerateContent", *path)

	config, ok := (*body)["generationConfig"].(map[string]any)
	require.True(t, ok, "generationConfig missing from request body")
	thinking, ok := config["thinkingConfig"].(map[string]any)
	require.True(t, ok, "thinkingConfig missing from generationConfig")
	assert.Equal(t, "medium", thinking["thinkingLevel"])
}

func TestGemini38FlashRequest(t *testing.T) {
	path, body := captureGoogleRequest(t)

	res, _ := streamAndCollect(
		t,
		NewGoogle([]string{"test-key"}),
		models.Gemini38Flash{ThinkingLevel: models.MediumThinkingLevel},
	)

	assert.Equal(t, "Hello", res.Content)
	assert.Equal(t, "/gemini-3.8-flash:streamGenerateContent", *path)

	config, ok := (*body)["generationConfig"].(map[string]any)
	require.True(t, ok, "generationConfig missing from request body")
	thinking, ok := config["thinkingConfig"].(map[string]any)
	require.True(t, ok, "thinkingConfig missing from generationConfig")
	assert.Equal(t, "medium", thinking["thinkingLevel"])
}

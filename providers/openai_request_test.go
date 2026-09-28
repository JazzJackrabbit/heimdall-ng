package providers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JazzJackrabbit/heimdall-ng/models"
	"github.com/JazzJackrabbit/heimdall-ng/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureOpenAIRequest points the OpenAI provider at a test server that
// records the JSON request body and replies with a minimal streamed completion.
func captureOpenAIRequest(t *testing.T) *map[string]any {
	t.Helper()

	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if raw, err := io.ReadAll(r.Body); err == nil {
				_ = json.Unmarshal(raw, &body)
			}

			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`data: {"choices":[{"delta":{"content":"Hello"}}]}` + "\n\n" +
				`data: {"choices":[],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}` + "\n\n" +
				"data: [DONE]\n\n"))
		},
	))
	t.Cleanup(server.Close)

	original := openAIBaseURL
	openAIBaseURL = server.URL
	t.Cleanup(func() { openAIBaseURL = original })

	return &body
}

// GPT-6 models reject temperature and top_p at the default reasoning effort, so
// the provider must not forward them. Earlier models still receive temperature.
func TestOpenAISamplingParamsPerModel(t *testing.T) {
	cases := []struct {
		model           models.Model
		sendTemperature bool
	}{
		{models.GPT56Sol{}, true},
		{models.GPT6Astra{}, false},
		{models.GPT6Sol{}, false},
		{models.GPT6Luna{}, false},
	}

	for _, tc := range cases {
		body := captureOpenAIRequest(t)

		res, err := NewOpenAI([]string{"test-key"}).StreamResponse(
			context.Background(),
			http.Client{},
			request.Completion{
				Model:         tc.model,
				SystemMessage: "you are a helpful assistant.",
				UserMessage:   "Say hello.",
			},
			func(string) error { return nil },
			nil,
		)
		require.NoError(t, err, "%T: StreamResponse returned an unexpected error", tc.model)
		assert.Equal(t, "Hello", res.Content, "%T: content", tc.model)
		assert.Equal(t, tc.model.GetName(), (*body)["model"], "%T: model name", tc.model)

		if tc.sendTemperature {
			assert.Contains(t, *body, "temperature", "%T: temperature should be sent", tc.model)
		} else {
			assert.NotContains(t, *body, "temperature", "%T: temperature must be omitted", tc.model)
		}
		assert.NotContains(t, *body, "top_p", "%T: top_p must be omitted", tc.model)
	}
}

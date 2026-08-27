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

// captureAnthropicRequest points the Anthropic provider at a test server that
// records the JSON request body and replies with a minimal streamed message.
func captureAnthropicRequest(t *testing.T) *map[string]any {
	t.Helper()

	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if raw, err := io.ReadAll(r.Body); err == nil {
				_ = json.Unmarshal(raw, &body)
			}

			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("event: message_start\n" +
				`data: {"type":"message_start","message":{"id":"msg_1","usage":{"input_tokens":1,"output_tokens":1}}}` + "\n\n" +
				"event: content_block_delta\n" +
				`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}` + "\n\n" +
				"event: message_delta\n" +
				`data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":1}}` + "\n\n" +
				"event: message_stop\n" +
				`data: {"type":"message_stop"}` + "\n\n"))
		},
	))
	t.Cleanup(server.Close)

	original := anthropicBaseUrl
	anthropicBaseUrl = server.URL
	t.Cleanup(func() { anthropicBaseUrl = original })

	return &body
}

// Claude 4.7 and later reject temperature and top_p with a 400 when they are
// set, so the provider must not forward them. Earlier models still accept them.
func TestAnthropicSamplingParamsPerModel(t *testing.T) {
	cases := []struct {
		model        models.Model
		sendSampling bool
	}{
		{models.Claude46Sonnet{}, true},
		{models.Claude45Haiku{}, true},
		{models.Claude47Opus{}, false},
		{models.Claude48Opus{}, false},
		{models.Claude5Opus{}, false},
		{models.Claude5Sonnet{}, false},
		{models.ClaudeFable5{}, false},
	}

	for _, tc := range cases {
		body := captureAnthropicRequest(t)

		_, err := NewAnthropic([]string{"test-key"}).StreamResponse(
			context.Background(),
			http.Client{},
			request.Completion{
				Model:         tc.model,
				SystemMessage: "you are a helpful assistant.",
				UserMessage:   "Say hello.",
				Temperature:   0.7,
				TopP:          0.9,
			},
			func(string) error { return nil },
			nil,
		)
		require.NoError(t, err, "%T: StreamResponse returned an unexpected error", tc.model)

		if tc.sendSampling {
			assert.InDelta(t, 0.7, (*body)["temperature"], 0.001, "%T: temperature should be sent", tc.model)
			assert.InDelta(t, 0.9, (*body)["top_p"], 0.001, "%T: top_p should be sent", tc.model)
		} else {
			assert.NotContains(t, *body, "temperature", "%T: temperature must be omitted", tc.model)
			assert.NotContains(t, *body, "top_p", "%T: top_p must be omitted", tc.model)
		}
	}
}

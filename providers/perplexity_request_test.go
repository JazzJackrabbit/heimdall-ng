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

// capturePerplexityRequest points the Perplexity provider at a test server
// that records the request path and JSON body and replies with the given
// Agent API event stream.
func capturePerplexityRequest(t *testing.T, events string) (*string, *map[string]any) {
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
			_, _ = w.Write([]byte(events))
		},
	))
	t.Cleanup(server.Close)

	original := perplexityBaseUrl
	perplexityBaseUrl = server.URL + "/v1/agent"
	t.Cleanup(func() { perplexityBaseUrl = original })

	return &path, &body
}

const perplexityCompletedStream = "event: response.output_text.delta\n" +
	`data: {"type":"response.output_text.delta","delta":"Hi"}` + "\n\n" +
	"event: response.completed\n" +
	`data: {"type":"response.completed","response":{"model":"openai/gpt-5.6-luna","usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2}}}` + "\n\n"

// The Sonar types map to the Agent API presets Perplexity recommends, with the
// system message as instructions and the history replayed as message items.
func TestPerplexityAgentRequestShape(t *testing.T) {
	path, body := capturePerplexityRequest(t, perplexityCompletedStream)

	schema := map[string]any{
		"type":       "object",
		"properties": map[string]any{"answer": map[string]any{"type": "string"}},
	}
	res, err := NewPerplexity([]string{"test-key"}).CompleteResponse(
		context.Background(),
		request.Completion{
			Model:         models.SonarPro{StructuredOutput: schema},
			SystemMessage: "Answer briefly.",
			UserMessage:   "And now?",
			History: []request.Message{
				{Role: "user", Content: "What time is it in Oslo?"},
				{Role: "assistant", Content: "It is noon."},
			},
			Temperature: 0.5,
		},
		http.Client{},
		nil,
	)
	require.NoError(t, err, "CompleteResponse returned an unexpected error")
	assert.Equal(t, "Hi", res.Content)
	assert.Equal(t, "openai/gpt-5.6-luna", res.Model, "the model the preset ran should be reported")
	assert.Equal(t, "/v1/agent", *path)

	assert.Equal(t, "low", (*body)["preset"], "sonar-pro maps to the low preset")
	assert.NotContains(t, *body, "model", "presets choose the model")
	assert.Equal(t, "Answer briefly.", (*body)["instructions"])
	assert.Equal(t, true, (*body)["stream"])
	assert.InDelta(t, 0.5, (*body)["temperature"], 0.001)

	input, ok := (*body)["input"].([]any)
	require.True(t, ok, "input should be an array of message items")
	require.Len(t, input, 3, "history plus the current user message")
	for i, item := range input {
		message, ok := item.(map[string]any)
		require.True(t, ok, "input item %d should be an object", i)
		assert.Equal(t, "message", message["type"], "input item %d type", i)
	}
	last, _ := input[2].(map[string]any)
	assert.Equal(t, "user", last["role"])
	assert.Equal(t, "And now?", last["content"])

	format, ok := (*body)["response_format"].(map[string]any)
	require.True(t, ok, "response_format should be set for structured output")
	assert.Equal(t, "json_schema", format["type"])
	jsonSchema, ok := format["json_schema"].(map[string]any)
	require.True(t, ok, "json_schema should be an object")
	assert.Equal(t, "response", jsonSchema["name"])
	assert.NotEmpty(t, jsonSchema["schema"])
}

func TestPerplexityPresetPerModel(t *testing.T) {
	cases := []struct {
		model  models.Model
		preset string
	}{
		{models.Sonar{}, "fast"},
		{models.SonarPro{}, "low"},
		{models.SonarReasoningPro{}, "medium"},
	}

	for _, tc := range cases {
		_, body := capturePerplexityRequest(t, perplexityCompletedStream)

		_, err := NewPerplexity([]string{"test-key"}).CompleteResponse(
			context.Background(),
			request.Completion{
				Model:         tc.model,
				SystemMessage: "sys",
				UserMessage:   "hello",
			},
			http.Client{},
			nil,
		)
		require.NoError(t, err, "%T: CompleteResponse returned an unexpected error", tc.model)
		assert.Equal(t, tc.preset, (*body)["preset"], "%T: preset", tc.model)
		assert.NotContains(t, *body, "response_format", "%T: no schema means no response_format", tc.model)
	}
}

// A failed response ends the stream with an error rather than an empty
// completion.
func TestPerplexityFailedResponseIsAnError(t *testing.T) {
	capturePerplexityRequest(t, "event: response.failed\n"+
		`data: {"type":"response.failed","response":{"error":{"message":"preset unavailable"}}}`+"\n\n")

	_, err := NewPerplexity([]string{"test-key"}).CompleteResponse(
		context.Background(),
		request.Completion{
			Model:         models.Sonar{},
			SystemMessage: "sys",
			UserMessage:   "hello",
		},
		http.Client{},
		nil,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "preset unavailable")
}

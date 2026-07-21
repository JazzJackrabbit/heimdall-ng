package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JazzJackrabbit/heimdall-ng/models"
	"github.com/JazzJackrabbit/heimdall-ng/request"
	"github.com/JazzJackrabbit/heimdall-ng/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sseServer returns an httptest server that replies to any request with the
// given SSE body.
func sseServer(t *testing.T, body string) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(body))
		},
	))
	t.Cleanup(server.Close)

	return server
}

func streamAndCollect(
	t *testing.T,
	provider LLMProvider,
	model models.Model,
) (response.Completion, string) {
	t.Helper()

	var streamed strings.Builder
	res, err := provider.StreamResponse(
		context.Background(),
		http.Client{},
		request.Completion{
			Model:         model,
			SystemMessage: "you are a helpful assistant.",
			UserMessage:   "Say hello.",
		},
		func(chunk string) error {
			streamed.WriteString(chunk)
			return nil
		},
		nil,
	)
	require.NoError(t, err, "StreamResponse returned an unexpected error")

	return res, streamed.String()
}

func TestAnthropicStreamingUsageParsing(t *testing.T) {
	body := "event: message_start\n" +
		`data: {"type":"message_start","message":{"id":"msg_1","usage":{"input_tokens":12,"output_tokens":1}}}` + "\n\n" +
		"event: content_block_delta\n" +
		`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}` + "\n\n" +
		"event: content_block_delta\n" +
		`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"!"}}` + "\n\n" +
		"event: message_delta\n" +
		`data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":7}}` + "\n\n" +
		"event: message_stop\n" +
		`data: {"type":"message_stop"}` + "\n\n"

	server := sseServer(t, body)
	original := anthropicBaseUrl
	anthropicBaseUrl = server.URL
	t.Cleanup(func() { anthropicBaseUrl = original })

	res, streamed := streamAndCollect(
		t,
		NewAnthropic([]string{"test-key"}),
		models.Claude45Haiku{},
	)

	assert.Equal(t, "Hello!", res.Content)
	assert.Equal(t, "Hello!", streamed)
	assert.Equal(t, 12, res.Usage.PromptTokens)
	assert.Equal(t, 7, res.Usage.CompletionTokens)
	assert.Equal(t, 19, res.Usage.TotalTokens)
}

func TestOpenAIStreamingUsageParsing(t *testing.T) {
	body := `data: {"choices":[{"delta":{"content":"Hello"}}]}` + "\n\n" +
		`data: {"choices":[{"delta":{"content":"!"}}]}` + "\n\n" +
		`data: {"choices":[],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}` + "\n\n" +
		"data: [DONE]\n\n"

	server := sseServer(t, body)
	original := openAIBaseURL
	openAIBaseURL = server.URL
	t.Cleanup(func() { openAIBaseURL = original })

	res, streamed := streamAndCollect(
		t,
		NewOpenAI([]string{"test-key"}),
		models.GPT4OMini{},
	)

	assert.Equal(t, "Hello!", res.Content)
	assert.Equal(t, "Hello!", streamed)
	assert.Equal(t, 10, res.Usage.PromptTokens)
	assert.Equal(t, 5, res.Usage.CompletionTokens)
	assert.Equal(t, 15, res.Usage.TotalTokens)
}

func TestGrokStreamingUsageParsing(t *testing.T) {
	body := `data: {"choices":[{"delta":{"content":"Hello"}}]}` + "\n\n" +
		`data: {"choices":[],"usage":{"prompt_tokens":8,"completion_tokens":3,"total_tokens":11}}` + "\n\n" +
		"data: [DONE]\n\n"

	server := sseServer(t, body)
	original := grokBaseURL
	grokBaseURL = server.URL
	t.Cleanup(func() { grokBaseURL = original })

	res, streamed := streamAndCollect(
		t,
		NewGrok([]string{"test-key"}),
		models.Grok3{},
	)

	assert.Equal(t, "Hello", res.Content)
	assert.Equal(t, "Hello", streamed)
	assert.Equal(t, 8, res.Usage.PromptTokens)
	assert.Equal(t, 3, res.Usage.CompletionTokens)
	assert.Equal(t, 11, res.Usage.TotalTokens)
}

func TestOpenRouterStreamingUsageParsing(t *testing.T) {
	body := `data: {"choices":[{"delta":{"content":"Hello"}}]}` + "\n\n" +
		`data: {"choices":[],"usage":{"prompt_tokens":20,"completion_tokens":6,"total_tokens":26}}` + "\n\n" +
		"data: [DONE]\n\n"

	server := sseServer(t, body)
	original := openRouterBaseURL
	openRouterBaseURL = server.URL
	t.Cleanup(func() { openRouterBaseURL = original })

	res, streamed := streamAndCollect(
		t,
		NewOpenRouter([]string{"test-key"}),
		models.OpenRouterModel{ModelName: "openai/gpt-4o-mini"},
	)

	assert.Equal(t, "Hello", res.Content)
	assert.Equal(t, "Hello", streamed)
	assert.Equal(t, 20, res.Usage.PromptTokens)
	assert.Equal(t, 6, res.Usage.CompletionTokens)
	assert.Equal(t, 26, res.Usage.TotalTokens)
}

func TestGoogleStreamingUsageParsing(t *testing.T) {
	// The final chunk deliberately finishes with MAX_TOKENS to verify that
	// usage is captured regardless of finish reason.
	body := `data: {"candidates":[{"content":{"parts":[{"text":"Hello"}],"role":"model"}}]}` + "\n\n" +
		`data: {"candidates":[{"content":{"parts":[{"text":"!"}],"role":"model"},"finishReason":"MAX_TOKENS"}],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":5,"totalTokenCount":15}}` + "\n\n"

	server := sseServer(t, body)
	original := googleBaseURL
	googleBaseURL = server.URL + "/%s:streamGenerateContent?alt=sse&key=%s"
	t.Cleanup(func() { googleBaseURL = original })

	res, streamed := streamAndCollect(
		t,
		NewGoogle([]string{"test-key"}),
		models.Gemini25FlashPreview{},
	)

	assert.Equal(t, "Hello!", res.Content)
	assert.Equal(t, "Hello!", streamed)
	assert.Equal(t, 10, res.Usage.PromptTokens)
	assert.Equal(t, 5, res.Usage.CompletionTokens)
	assert.Equal(t, 15, res.Usage.TotalTokens)
}

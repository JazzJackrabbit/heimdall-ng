package providers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/JazzJackrabbit/heimdall-ng/models"
	"github.com/JazzJackrabbit/heimdall-ng/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// aiandServer points the ai& provider at a test server that records each
// request body and answers with the given status and body.
func aiandServer(t *testing.T, status int, body string) (*[]map[string]any, *atomic.Int32) {
	t.Helper()

	var bodies []map[string]any
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			assert.Equal(t, "/chat/completions", r.URL.Path)
			assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
			raw, _ := io.ReadAll(r.Body)
			var parsed map[string]any
			assert.NoError(t, json.Unmarshal(raw, &parsed))
			bodies = append(bodies, parsed)

			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
		},
	))
	t.Cleanup(server.Close)

	original := aiandBaseURL
	aiandBaseURL = server.URL
	t.Cleanup(func() { aiandBaseURL = original })

	return &bodies, &calls
}

func TestAIAndStreamingSeparatesReasoningFromContent(t *testing.T) {
	body := ": keep-alive\n\n" +
		`data: {"choices":[{"index":0,"delta":{"reasoning_content":"17 x 23","reasoning":"17 x 23"},"finish_reason":null}]}` + "\n\n" +
		`data: {"choices":[{"index":0,"delta":{"content":"39"},"finish_reason":null}]}` + "\n\n" +
		`data: {"choices":[{"index":0,"delta":{"content":"1"},"finish_reason":"stop"}]}` + "\n\n" +
		`data: {"choices":[],"usage":{"prompt_tokens":10,"completion_tokens":9,"total_tokens":19}}` + "\n\n" +
		"data: [DONE]\n\n"
	bodies, _ := aiandServer(t, http.StatusOK, body)

	res, streamed := streamAndCollect(t, NewAIAnd([]string{"test-key"}), models.AIAndGLM53{
		ReasoningEffort: models.AIAndReasoningLow,
		MaxOutputTokens: 2048,
	})

	assert.Equal(t, "391", res.Content)
	assert.Equal(t, "391", streamed, "reasoning must not reach the chunk handler")
	assert.Equal(t, "17 x 23", res.Thoughts)
	assert.Equal(t, models.AIAndGLM53Alias, res.Model)
	assert.Equal(t, 10, res.Usage.PromptTokens)
	assert.Equal(t, 9, res.Usage.CompletionTokens)
	assert.Equal(t, 19, res.Usage.TotalTokens)

	require.Len(t, *bodies, 1)
	sent := (*bodies)[0]
	assert.Equal(t, "zai-org/glm-5.3", sent["model"])
	assert.Equal(t, true, sent["stream"])
	assert.Equal(t, "low", sent["reasoning_effort"])
	assert.InDelta(t, 2048, sent["max_completion_tokens"], 0)
	assert.NotContains(t, sent, "temperature", "unset temperature is left to the model")

	messages := sent["messages"].([]any)
	require.Len(t, messages, 2)
	assert.Equal(t, map[string]any{"role": "system", "content": "you are a helpful assistant."}, messages[0])
	assert.Equal(t, map[string]any{"role": "user", "content": "Say hello."}, messages[1])
}

func TestAIAndVisionRequestCarriesImagesPdfsAndSchema(t *testing.T) {
	body := `data: {"choices":[{"index":0,"delta":{"content":"{}"},"finish_reason":"stop"}]}` + "\n\n" + "data: [DONE]\n\n"
	bodies, _ := aiandServer(t, http.StatusOK, body)

	schema := map[string]any{"name": "verdict", "strict": true, "schema": map[string]any{"type": "object"}}
	_, err := NewAIAnd([]string{"test-key"}).CompleteResponse(context.Background(), request.Completion{
		Model: models.AIAndKimiK3{
			ImageFile:        []models.AIAndImagePayload{{URL: "https://example.com/cat.png"}},
			PdfFiles:         []models.AIAndPdf{"JVBERi0x"},
			StructuredOutput: schema,
		},
		UserMessage: "Describe both.",
		Temperature: 0.7,
	}, http.Client{}, nil)
	require.NoError(t, err)

	require.Len(t, *bodies, 1)
	sent := (*bodies)[0]
	assert.InDelta(t, 0.7, sent["temperature"], 1e-6)
	assert.NotContains(t, sent, "reasoning_effort")
	assert.Equal(t, map[string]any{"type": "json_schema", "json_schema": map[string]any{
		"name": "verdict", "strict": true, "schema": map[string]any{"type": "object"},
	}}, sent["response_format"])

	messages := sent["messages"].([]any)
	require.Len(t, messages, 1, "no system message when none is given")
	parts := messages[0].(map[string]any)["content"].([]any)
	require.Len(t, parts, 3)
	assert.Equal(t, map[string]any{
		"type":      "image_url",
		"image_url": map[string]any{"url": "https://example.com/cat.png", "detail": "auto"},
	}, parts[0])
	assert.Equal(t, map[string]any{
		"type": "file",
		"file": map[string]any{"filename": "document_1.pdf", "file_data": "data:application/pdf;base64,JVBERi0x"},
	}, parts[1])
	assert.Equal(t, map[string]any{"type": "text", "text": "Describe both."}, parts[2])
}

func TestAIAndBadRequestIsNotRetried(t *testing.T) {
	body := `{"error":{"message":"Model 'zai-org/glm-5.3' cannot disable reasoning","type":"invalid_request_error","code":"invalid_value"}}`
	_, calls := aiandServer(t, http.StatusBadRequest, body)

	_, err := NewAIAnd([]string{"test-key", "test-key"}).CompleteResponse(context.Background(), request.Completion{
		Model:       models.AIAndGLM53{ReasoningEffort: models.AIAndReasoningNone},
		UserMessage: "Hello",
	}, http.Client{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot disable reasoning (invalid_value)")
	assert.Equal(t, int32(1), calls.Load(), "a 400 fails on every key, so it is sent once")
}

func TestAIAndReasoningCutOffByTokenLimitIsAnError(t *testing.T) {
	body := `data: {"choices":[{"index":0,"delta":{"reasoning_content":"thinking"},"finish_reason":"length"}]}` + "\n\n" + "data: [DONE]\n\n"
	aiandServer(t, http.StatusOK, body)

	_, err := NewAIAnd([]string{"test-key"}).CompleteResponse(context.Background(), request.Completion{
		Model:       models.AIAndDeepSeekV4Flash{MaxOutputTokens: 16},
		UserMessage: "Hello",
	}, http.Client{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "output token limit")
}

func TestAIAndRejectsOtherModels(t *testing.T) {
	_, calls := aiandServer(t, http.StatusOK, "")

	_, err := NewAIAnd([]string{"test-key"}).CompleteResponse(context.Background(), request.Completion{
		Model:       models.Grok43{},
		UserMessage: "Hello",
	}, http.Client{}, nil)
	require.Error(t, err)
	assert.Equal(t, int32(0), calls.Load())
}

func TestAIAndModelCatalog(t *testing.T) {
	seen := map[string]bool{}
	for _, m := range models.AIAndModels() {
		name := m.GetName()
		assert.False(t, seen[name], "duplicate model %s", name)
		seen[name] = true
		assert.Equal(t, models.AIAndProvider, m.GetProvider(), name)
		assert.Contains(t, name, "/", "%s keeps its lab prefix", name)
		assert.Positive(t, m.GetOutputCostPer1M(), name)
		assert.Positive(t, m.ContextWindow(), name)
		assert.Contains(t, m.ReasoningEfforts(), m.DefaultReasoningEffort(), name)
		assert.Contains(t, models.GetAll(), name)
	}
	assert.Len(t, seen, 13)
}

func TestAIAndLiveCompletion(t *testing.T) {
	apiKey := os.Getenv("AIAND_API_KEY")
	if apiKey == "" {
		t.Skip("AIAND_API_KEY not set")
	}

	res, err := NewAIAnd([]string{apiKey}).CompleteResponse(context.Background(), request.Completion{
		Model:         models.AIAndDeepSeekV4Flash{ReasoningEffort: models.AIAndReasoningNone},
		SystemMessage: "you are a helpful assistant.",
		UserMessage:   "Say hello in one sentence.",
	}, http.Client{Timeout: 2 * time.Minute}, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.Content)
	assert.Positive(t, res.Usage.TotalTokens)
}

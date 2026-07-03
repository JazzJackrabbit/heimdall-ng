//go:build perplexity

package providers

import (
	"testing"

	"github.com/JazzJackrabbit/heimdall/models"
	"github.com/stretchr/testify/assert"
)

func TestPerplexityStreamingUsageParsing(t *testing.T) {
	body := `data: {"choices":[{"delta":{"content":"Hello"}}],"usage":{"prompt_tokens":9,"completion_tokens":4,"total_tokens":13}}` + "\n\n" +
		"data: [DONE]\n\n"

	server := sseServer(t, body)
	original := perplexityBaseUrl
	perplexityBaseUrl = server.URL
	t.Cleanup(func() { perplexityBaseUrl = original })

	res, streamed := streamAndCollect(
		t,
		NewPerplexity([]string{"test-key"}),
		models.Sonar{},
	)

	assert.Equal(t, "Hello", res.Content)
	assert.Equal(t, "Hello", streamed)
	assert.Equal(t, 9, res.Usage.PromptTokens)
	assert.Equal(t, 4, res.Usage.CompletionTokens)
	assert.Equal(t, 13, res.Usage.TotalTokens)
}

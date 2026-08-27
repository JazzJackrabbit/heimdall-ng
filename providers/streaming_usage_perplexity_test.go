package providers

import (
	"testing"

	"github.com/JazzJackrabbit/heimdall-ng/models"
	"github.com/stretchr/testify/assert"
)

func TestPerplexityStreamingUsageParsing(t *testing.T) {
	body := "event: response.output_text.delta\n" +
		`data: {"type":"response.output_text.delta","delta":"Hello"}` + "\n\n" +
		"event: response.completed\n" +
		`data: {"type":"response.completed","response":{"id":"r1","model":"openai/gpt-5.6-luna","output_text":"Hello","usage":{"input_tokens":9,"output_tokens":4,"total_tokens":13}}}` + "\n\n"

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
	assert.Equal(t, "openai/gpt-5.6-luna", res.Model)
	assert.Equal(t, 9, res.Usage.PromptTokens)
	assert.Equal(t, 4, res.Usage.CompletionTokens)
	assert.Equal(t, 13, res.Usage.TotalTokens)
}

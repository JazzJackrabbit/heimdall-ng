package models

const PerplexityProvider = "perplexity"

// SonarReasoningPro runs through the Perplexity Agent API with the medium
// preset, the replacement Perplexity recommends for sonar-reasoning-pro.
type SonarReasoningPro struct {
	StructuredOutput map[string]any
}

func (s SonarReasoningPro) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000002
}

func (s SonarReasoningPro) GetName() string {
	return "sonar-reasoning-pro"
}

func (s SonarReasoningPro) GetProvider() string {
	return PerplexityProvider
}

var _ Model = new(SonarReasoningPro)

// Deprecated: sonar-reasoning was removed from the Perplexity API on December 15, 2025. Requests to this model will fail.
// Use SonarReasoningPro (sonar-reasoning-pro) as a replacement. Requests are sent with the medium preset.
type SonarReasoning struct {
	StructuredOutput map[string]any
}

func (s SonarReasoning) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000001
}

func (s SonarReasoning) GetName() string {
	return "sonar-reasoning"
}

func (s SonarReasoning) GetProvider() string {
	return PerplexityProvider
}

var _ Model = new(SonarReasoning)

// SonarPro runs through the Perplexity Agent API with the low preset, the
// replacement Perplexity recommends for sonar-pro.
type SonarPro struct {
	StructuredOutput map[string]any
}

func (s SonarPro) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000003
}

func (s SonarPro) GetName() string {
	return "sonar-pro"
}

func (s SonarPro) GetProvider() string {
	return PerplexityProvider
}

var _ Model = new(SonarPro)

// Sonar runs through the Perplexity Agent API with the fast preset, the
// replacement Perplexity recommends for sonar.
type Sonar struct {
	StructuredOutput map[string]any
}

func (s Sonar) EstimateCost(text string) float64 {
	return (float64(len(text)) / 4) * 0.000001
}

func (s Sonar) GetName() string {
	return "sonar"
}

func (s Sonar) GetProvider() string {
	return PerplexityProvider
}

var _ Model = new(Sonar)

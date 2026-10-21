package models

// AIAndProvider routes to ai& (https://docs.aiand.com), an OpenAI-compatible
// API serving open-weight models from Japan. Model ids keep the lab prefix
// ai& uses (for example "zai-org/glm-5.3").
const AIAndProvider = "aiand"

const (
	AIAndKimiK3Alias           = "moonshotai/kimi-k3"
	AIAndKimiK27CodeAlias      = "moonshotai/kimi-k2.7-code"
	AIAndGLM53Alias            = "zai-org/glm-5.3"
	AIAndGLM53FlashAlias       = "zai-org/glm-5.3-flash"
	AIAndGLM52Alias            = "zai-org/glm-5.2"
	AIAndDeepSeekV4ProAlias    = "deepseek-ai/deepseek-v4-pro"
	AIAndDeepSeekV41FlashAlias = "deepseek-ai/deepseek-v4.1-flash"
	AIAndDeepSeekV4FlashAlias  = "deepseek-ai/deepseek-v4-flash"
	AIAndQwen38_27BAlias       = "qwen/qwen3.8-27b"
	AIAndQwen36_27BAlias       = "qwen/qwen3.6-27b"
	AIAndGemma4_31BAlias       = "google/gemma-4-31b-it"
	AIAndGPTOSS120BAlias       = "openai/gpt-oss-120b"
	AIAndMotif3Alias           = "motif-technologies/motif-3"
)

// AIAndReasoningEffort is the reasoning_effort sent to ai&. Each model accepts
// its own subset, listed by its ReasoningEfforts method. A level a model does
// not list runs at its next listed level up, and "none" on a model that does
// not list it is rejected with a 400.
type AIAndReasoningEffort string

const (
	AIAndReasoningNone   AIAndReasoningEffort = "none"
	AIAndReasoningLow    AIAndReasoningEffort = "low"
	AIAndReasoningMedium AIAndReasoningEffort = "medium"
	AIAndReasoningHigh   AIAndReasoningEffort = "high"
	AIAndReasoningXHigh  AIAndReasoningEffort = "xhigh"
	AIAndReasoningMax    AIAndReasoningEffort = "max"
)

// AIAndImagePayload is an image for a vision model: a public https URL or a
// base64 data URL. ai& rejects http URLs, IP hosts and internal hostnames.
type AIAndImagePayload struct {
	URL    string
	Detail string
}

// AIAndPdf is a base64-encoded PDF for a model with the document capability.
// ai& accepts up to 100 pages per PDF and 16 MiB of new PDFs per request.
type AIAndPdf string

// AIAndParams is what the ai& provider sends for a model. Text-only models
// leave ImageFile and PdfFiles empty.
type AIAndParams struct {
	ImageFile        []AIAndImagePayload
	PdfFiles         []AIAndPdf
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	MaxOutputTokens  int
}

// AIAndModel is implemented by every ai& model type.
type AIAndModel interface {
	Model
	CostBreakdown
	GetAIAndParams() AIAndParams
	// ReasoningEfforts lists the reasoning_effort values the model accepts.
	ReasoningEfforts() []AIAndReasoningEffort
	// DefaultReasoningEffort is the level ai& suggests as a UI default. It is
	// not sent when ReasoningEffort is empty: the model's own default applies.
	DefaultReasoningEffort() AIAndReasoningEffort
	// ContextWindow is the maximum combined input and output tokens.
	ContextWindow() int
	// MaxOutputTokenLimit is the largest max_completion_tokens the model takes.
	MaxOutputTokenLimit() int
	// AcceptsImages reports the vision capability; such models also take PDFs.
	AcceptsImages() bool
}

// AIAndKimiK3 is moonshotai/kimi-k3: Multimodal Kimi model with 1M context and toggleable max-effort thinking for long-horizon agent work.
type AIAndKimiK3 struct {
	ImageFile        []AIAndImagePayload
	PdfFiles         []AIAndPdf
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndKimiK3) GetName() string { return AIAndKimiK3Alias }

func (AIAndKimiK3) GetProvider() string { return AIAndProvider }

func (m AIAndKimiK3) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndKimiK3) GetInputCostPer1M() float64 { return 3.0 }

func (AIAndKimiK3) GetOutputCostPer1M() float64 { return 12.5 }

func (AIAndKimiK3) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningLow, AIAndReasoningHigh, AIAndReasoningMax}
}

func (AIAndKimiK3) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningMax }

func (AIAndKimiK3) ContextWindow() int { return 1048576 }

func (AIAndKimiK3) MaxOutputTokenLimit() int { return 131072 }

func (AIAndKimiK3) AcceptsImages() bool { return true }

func (m AIAndKimiK3) GetAIAndParams() AIAndParams {
	return AIAndParams{
		ImageFile:        m.ImageFile,
		PdfFiles:         m.PdfFiles,
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndKimiK3)

// AIAndKimiK27Code is moonshotai/kimi-k2.7-code: Coding-focused Kimi model, stronger on long-horizon repo work with less overthinking.
type AIAndKimiK27Code struct {
	ImageFile        []AIAndImagePayload
	PdfFiles         []AIAndPdf
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndKimiK27Code) GetName() string { return AIAndKimiK27CodeAlias }

func (AIAndKimiK27Code) GetProvider() string { return AIAndProvider }

func (m AIAndKimiK27Code) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndKimiK27Code) GetInputCostPer1M() float64 { return 0.75 }

func (AIAndKimiK27Code) GetOutputCostPer1M() float64 { return 3.5 }

func (AIAndKimiK27Code) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningHigh}
}

func (AIAndKimiK27Code) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningHigh }

func (AIAndKimiK27Code) ContextWindow() int { return 262144 }

func (AIAndKimiK27Code) MaxOutputTokenLimit() int { return 262144 }

func (AIAndKimiK27Code) AcceptsImages() bool { return true }

func (m AIAndKimiK27Code) GetAIAndParams() AIAndParams {
	return AIAndParams{
		ImageFile:        m.ImageFile,
		PdfFiles:         m.PdfFiles,
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndKimiK27Code)

// AIAndGLM53 is zai-org/glm-5.3: Flagship GLM model for long-horizon coding, agents, and complex project delivery.
// Text input only.
type AIAndGLM53 struct {
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndGLM53) GetName() string { return AIAndGLM53Alias }

func (AIAndGLM53) GetProvider() string { return AIAndProvider }

func (m AIAndGLM53) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndGLM53) GetInputCostPer1M() float64 { return 1.0 }

func (AIAndGLM53) GetOutputCostPer1M() float64 { return 4.0 }

func (AIAndGLM53) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningLow, AIAndReasoningHigh, AIAndReasoningMax}
}

func (AIAndGLM53) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningMax }

func (AIAndGLM53) ContextWindow() int { return 1048576 }

func (AIAndGLM53) MaxOutputTokenLimit() int { return 131072 }

func (AIAndGLM53) AcceptsImages() bool { return false }

func (m AIAndGLM53) GetAIAndParams() AIAndParams {
	return AIAndParams{
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndGLM53)

// AIAndGLM53Flash is zai-org/glm-5.3-flash: Multimodal GLM model for coding and agents with a million-token context.
type AIAndGLM53Flash struct {
	ImageFile        []AIAndImagePayload
	PdfFiles         []AIAndPdf
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndGLM53Flash) GetName() string { return AIAndGLM53FlashAlias }

func (AIAndGLM53Flash) GetProvider() string { return AIAndProvider }

func (m AIAndGLM53Flash) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndGLM53Flash) GetInputCostPer1M() float64 { return 0.15 }

func (AIAndGLM53Flash) GetOutputCostPer1M() float64 { return 0.5 }

func (AIAndGLM53Flash) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningLow, AIAndReasoningHigh, AIAndReasoningMax}
}

func (AIAndGLM53Flash) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningMax }

func (AIAndGLM53Flash) ContextWindow() int { return 1048550 }

func (AIAndGLM53Flash) MaxOutputTokenLimit() int { return 131072 }

func (AIAndGLM53Flash) AcceptsImages() bool { return true }

func (m AIAndGLM53Flash) GetAIAndParams() AIAndParams {
	return AIAndParams{
		ImageFile:        m.ImageFile,
		PdfFiles:         m.PdfFiles,
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndGLM53Flash)

// AIAndGLM52 is zai-org/glm-5.2: Open flagship GLM for long-horizon coding agents and million-token context work.
// Text input only.
type AIAndGLM52 struct {
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndGLM52) GetName() string { return AIAndGLM52Alias }

func (AIAndGLM52) GetProvider() string { return AIAndProvider }

func (m AIAndGLM52) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndGLM52) GetInputCostPer1M() float64 { return 1.0 }

func (AIAndGLM52) GetOutputCostPer1M() float64 { return 4.0 }

func (AIAndGLM52) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningNone, AIAndReasoningHigh, AIAndReasoningMax}
}

func (AIAndGLM52) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningMax }

func (AIAndGLM52) ContextWindow() int { return 1048576 }

func (AIAndGLM52) MaxOutputTokenLimit() int { return 131072 }

func (AIAndGLM52) AcceptsImages() bool { return false }

func (m AIAndGLM52) GetAIAndParams() AIAndParams {
	return AIAndParams{
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndGLM52)

// AIAndDeepSeekV4Pro is deepseek-ai/deepseek-v4-pro: Open MoE flagship with million-token context for coding and long agent runs.
// Text input only.
type AIAndDeepSeekV4Pro struct {
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndDeepSeekV4Pro) GetName() string { return AIAndDeepSeekV4ProAlias }

func (AIAndDeepSeekV4Pro) GetProvider() string { return AIAndProvider }

func (m AIAndDeepSeekV4Pro) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndDeepSeekV4Pro) GetInputCostPer1M() float64 { return 1.0 }

func (AIAndDeepSeekV4Pro) GetOutputCostPer1M() float64 { return 2.5 }

func (AIAndDeepSeekV4Pro) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningNone, AIAndReasoningHigh, AIAndReasoningMax}
}

func (AIAndDeepSeekV4Pro) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningHigh }

func (AIAndDeepSeekV4Pro) ContextWindow() int { return 1048576 }

func (AIAndDeepSeekV4Pro) MaxOutputTokenLimit() int { return 384000 }

func (AIAndDeepSeekV4Pro) AcceptsImages() bool { return false }

func (m AIAndDeepSeekV4Pro) GetAIAndParams() AIAndParams {
	return AIAndParams{
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndDeepSeekV4Pro)

// AIAndDeepSeekV41Flash is deepseek-ai/deepseek-v4.1-flash: DeepSeek V4.1 Flash lane for economical reasoning, coding and long-context work.
type AIAndDeepSeekV41Flash struct {
	ImageFile        []AIAndImagePayload
	PdfFiles         []AIAndPdf
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndDeepSeekV41Flash) GetName() string { return AIAndDeepSeekV41FlashAlias }

func (AIAndDeepSeekV41Flash) GetProvider() string { return AIAndProvider }

func (m AIAndDeepSeekV41Flash) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndDeepSeekV41Flash) GetInputCostPer1M() float64 { return 0.3 }

func (AIAndDeepSeekV41Flash) GetOutputCostPer1M() float64 { return 0.6 }

func (AIAndDeepSeekV41Flash) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningNone, AIAndReasoningHigh, AIAndReasoningMax}
}

func (AIAndDeepSeekV41Flash) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningHigh }

func (AIAndDeepSeekV41Flash) ContextWindow() int { return 1048576 }

func (AIAndDeepSeekV41Flash) MaxOutputTokenLimit() int { return 384000 }

func (AIAndDeepSeekV41Flash) AcceptsImages() bool { return true }

func (m AIAndDeepSeekV41Flash) GetAIAndParams() AIAndParams {
	return AIAndParams{
		ImageFile:        m.ImageFile,
		PdfFiles:         m.PdfFiles,
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndDeepSeekV41Flash)

// AIAndDeepSeekV4Flash is deepseek-ai/deepseek-v4-flash: Fast DeepSeek V4 lane for economical reasoning, coding, and long-context work.
// Text input only.
type AIAndDeepSeekV4Flash struct {
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndDeepSeekV4Flash) GetName() string { return AIAndDeepSeekV4FlashAlias }

func (AIAndDeepSeekV4Flash) GetProvider() string { return AIAndProvider }

func (m AIAndDeepSeekV4Flash) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndDeepSeekV4Flash) GetInputCostPer1M() float64 { return 0.15 }

func (AIAndDeepSeekV4Flash) GetOutputCostPer1M() float64 { return 0.25 }

func (AIAndDeepSeekV4Flash) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningNone, AIAndReasoningHigh, AIAndReasoningMax}
}

func (AIAndDeepSeekV4Flash) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningHigh }

func (AIAndDeepSeekV4Flash) ContextWindow() int { return 1048576 }

func (AIAndDeepSeekV4Flash) MaxOutputTokenLimit() int { return 384000 }

func (AIAndDeepSeekV4Flash) AcceptsImages() bool { return false }

func (m AIAndDeepSeekV4Flash) GetAIAndParams() AIAndParams {
	return AIAndParams{
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndDeepSeekV4Flash)

// AIAndQwen38_27B is qwen/qwen3.8-27b: Dense 27B vision-language model for coding, agent tasks, and image and video understanding.
type AIAndQwen38_27B struct {
	ImageFile        []AIAndImagePayload
	PdfFiles         []AIAndPdf
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndQwen38_27B) GetName() string { return AIAndQwen38_27BAlias }

func (AIAndQwen38_27B) GetProvider() string { return AIAndProvider }

func (m AIAndQwen38_27B) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndQwen38_27B) GetInputCostPer1M() float64 { return 0.4 }

func (AIAndQwen38_27B) GetOutputCostPer1M() float64 { return 3.0 }

func (AIAndQwen38_27B) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningNone, AIAndReasoningLow, AIAndReasoningMedium, AIAndReasoningXHigh}
}

func (AIAndQwen38_27B) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningXHigh }

func (AIAndQwen38_27B) ContextWindow() int { return 262144 }

func (AIAndQwen38_27B) MaxOutputTokenLimit() int { return 32768 }

func (AIAndQwen38_27B) AcceptsImages() bool { return true }

func (m AIAndQwen38_27B) GetAIAndParams() AIAndParams {
	return AIAndParams{
		ImageFile:        m.ImageFile,
		PdfFiles:         m.PdfFiles,
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndQwen38_27B)

// AIAndQwen36_27B is qwen/qwen3.6-27b: Qwen vision-language model for visual reasoning, documents, and agent tasks.
type AIAndQwen36_27B struct {
	ImageFile        []AIAndImagePayload
	PdfFiles         []AIAndPdf
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndQwen36_27B) GetName() string { return AIAndQwen36_27BAlias }

func (AIAndQwen36_27B) GetProvider() string { return AIAndProvider }

func (m AIAndQwen36_27B) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndQwen36_27B) GetInputCostPer1M() float64 { return 0.32 }

func (AIAndQwen36_27B) GetOutputCostPer1M() float64 { return 3.2 }

func (AIAndQwen36_27B) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningNone, AIAndReasoningHigh}
}

func (AIAndQwen36_27B) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningHigh }

func (AIAndQwen36_27B) ContextWindow() int { return 262144 }

func (AIAndQwen36_27B) MaxOutputTokenLimit() int { return 65536 }

func (AIAndQwen36_27B) AcceptsImages() bool { return true }

func (m AIAndQwen36_27B) GetAIAndParams() AIAndParams {
	return AIAndParams{
		ImageFile:        m.ImageFile,
		PdfFiles:         m.PdfFiles,
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndQwen36_27B)

// AIAndGemma4_31B is google/gemma-4-31b-it: Largest Gemma 4 instruction model for open, self-hosted chat and reasoning.
type AIAndGemma4_31B struct {
	ImageFile        []AIAndImagePayload
	PdfFiles         []AIAndPdf
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndGemma4_31B) GetName() string { return AIAndGemma4_31BAlias }

func (AIAndGemma4_31B) GetProvider() string { return AIAndProvider }

func (m AIAndGemma4_31B) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndGemma4_31B) GetInputCostPer1M() float64 { return 0.2 }

func (AIAndGemma4_31B) GetOutputCostPer1M() float64 { return 0.5 }

func (AIAndGemma4_31B) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningNone, AIAndReasoningHigh}
}

func (AIAndGemma4_31B) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningNone }

func (AIAndGemma4_31B) ContextWindow() int { return 262144 }

func (AIAndGemma4_31B) MaxOutputTokenLimit() int { return 32768 }

func (AIAndGemma4_31B) AcceptsImages() bool { return true }

func (m AIAndGemma4_31B) GetAIAndParams() AIAndParams {
	return AIAndParams{
		ImageFile:        m.ImageFile,
		PdfFiles:         m.PdfFiles,
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndGemma4_31B)

// AIAndGPTOSS120B is openai/gpt-oss-120b: Open GPT reasoning model for self-hosted agents and controllable deployments.
// Text input only.
type AIAndGPTOSS120B struct {
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndGPTOSS120B) GetName() string { return AIAndGPTOSS120BAlias }

func (AIAndGPTOSS120B) GetProvider() string { return AIAndProvider }

func (m AIAndGPTOSS120B) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndGPTOSS120B) GetInputCostPer1M() float64 { return 0.15 }

func (AIAndGPTOSS120B) GetOutputCostPer1M() float64 { return 0.6 }

func (AIAndGPTOSS120B) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningLow, AIAndReasoningMedium, AIAndReasoningHigh}
}

func (AIAndGPTOSS120B) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningMedium }

func (AIAndGPTOSS120B) ContextWindow() int { return 131072 }

func (AIAndGPTOSS120B) MaxOutputTokenLimit() int { return 32768 }

func (AIAndGPTOSS120B) AcceptsImages() bool { return false }

func (m AIAndGPTOSS120B) GetAIAndParams() AIAndParams {
	return AIAndParams{
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndGPTOSS120B)

// AIAndMotif3 is motif-technologies/motif-3: Motif 3 is a large-scale, decoder-only Mixture-of-Experts (MoE) language model with 314 billion total parameters and 13.2 billion parameters activated per token..
// Text input only.
type AIAndMotif3 struct {
	StructuredOutput map[string]any
	ReasoningEffort  AIAndReasoningEffort
	// MaxOutputTokens caps reasoning and answer together. 0 leaves the cap to ai&.
	MaxOutputTokens int
}

func (AIAndMotif3) GetName() string { return AIAndMotif3Alias }

func (AIAndMotif3) GetProvider() string { return AIAndProvider }

func (m AIAndMotif3) EstimateCost(text string) float64 {
	return estimateAIAndCost(m, text)
}

func (AIAndMotif3) GetInputCostPer1M() float64 { return 0.5 }

func (AIAndMotif3) GetOutputCostPer1M() float64 { return 2.0 }

func (AIAndMotif3) ReasoningEfforts() []AIAndReasoningEffort {
	return []AIAndReasoningEffort{AIAndReasoningNone, AIAndReasoningHigh}
}

func (AIAndMotif3) DefaultReasoningEffort() AIAndReasoningEffort { return AIAndReasoningHigh }

func (AIAndMotif3) ContextWindow() int { return 262144 }

func (AIAndMotif3) MaxOutputTokenLimit() int { return 262144 }

func (AIAndMotif3) AcceptsImages() bool { return false }

func (m AIAndMotif3) GetAIAndParams() AIAndParams {
	return AIAndParams{
		StructuredOutput: m.StructuredOutput,
		ReasoningEffort:  m.ReasoningEffort,
		MaxOutputTokens:  m.MaxOutputTokens,
	}
}

var _ AIAndModel = new(AIAndMotif3)

// estimateAIAndCost approximates a request's cost from its text at four
// characters per token, averaging the input and output rates.
func estimateAIAndCost(m CostBreakdown, text string) float64 {
	averagePer1M := (m.GetInputCostPer1M() + m.GetOutputCostPer1M()) / 2
	return (float64(len(text)) / 4) * averagePer1M / 1e6
}

// AIAndModels returns every ai& model type, newest families first.
func AIAndModels() []AIAndModel {
	return []AIAndModel{
		AIAndKimiK3{},
		AIAndKimiK27Code{},
		AIAndGLM53{},
		AIAndGLM53Flash{},
		AIAndGLM52{},
		AIAndDeepSeekV4Pro{},
		AIAndDeepSeekV41Flash{},
		AIAndDeepSeekV4Flash{},
		AIAndQwen38_27B{},
		AIAndQwen36_27B{},
		AIAndGemma4_31B{},
		AIAndGPTOSS120B{},
		AIAndMotif3{},
	}
}

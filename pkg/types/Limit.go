package types

type Limit struct {
	Gpt35Limit     int  `json:"gpt35_limit"`
	Gpt4Limit      int  `json:"gpt4_limit"`
	Dalle3Limit    int  `json:"dalle3_limit"`
	ClaudeLimit    int  `json:"claude_limit"`
	ContextSupport bool `json:"context_support"`
}

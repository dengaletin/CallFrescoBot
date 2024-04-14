package types

type Usage struct {
	Gpt35         int `json:"gpt35"`
	Gpt4          int `json:"gpt4"`
	Dalle3        int `json:"dalle3"`
	Claude        int `json:"claude"`
	Gpt35Context  int `json:"gpt35_context"`
	Gpt4Context   int `json:"gpt4_context"`
	Dalle3Context int `json:"dalle3_context"`
	ClaudeContext int `json:"claude_context"`
}

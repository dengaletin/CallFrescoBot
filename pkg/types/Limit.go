package types

type Limit struct {
	Gpt4OMiniLimit int  `json:"gpt4OMiniLimit"`
	Gpt4OLimit     int  `json:"Gpt4OLimit"`
	Dalle3Limit    int  `json:"Dalle3Limit"`
	Gpt4O1Limit    int  `json:"Gpt4O1Limit"`
	ContextSupport bool `json:"ContextSupport"`
}

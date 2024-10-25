package types

type Usage struct {
	Gpt4OMini             int `json:"Gpt4OMini"`
	Gpt4O                 int `json:"Gpt4O"`
	Dalle3                int `json:"Dalle3"`
	Gpt4O1                int `json:"Gpt4O1"`
	Gpt4OMiniContext      int `json:"Gpt4OMiniContext"`
	Gpt4OContext          int `json:"Gpt4OContext"`
	Dalle3Context         int `json:"Dalle3Context"`
	Gpt4O1Context         int `json:"Gpt4o1Context"`
	Gpt4OMiniVoice        int `json:"Gpt4OMiniVoice"`
	Gpt4OVoice            int `json:"Gpt4OVoice"`
	Dalle3Voice           int `json:"Dalle3Voice"`
	Gpt4O1Voice           int `json:"Gpt4O1Voice"`
	Gpt4OMiniContextVoice int `json:"Gpt4OMiniContextVoice"`
	Gpt4OContextVoice     int `json:"Gpt4OContextVoice"`
	Dalle3ContextVoice    int `json:"Dalle3ContextVoice"`
	Gpt4O1ContextVoice    int `json:"Gpt4o1ContextVoice"`
}

package types

type Usage struct {
	Gpt4OMini        int `json:"Gpt4OMini"`
	Gpt4O            int `json:"Gpt4O"`
	Dalle3           int `json:"Dalle3"`
	Gpt4O1           int `json:"Gpt4O1"`
	Gpt4OMiniContext int `json:"Gpt4OMiniContext"`
	Gpt4OContext     int `json:"Gpt4OContext"`
	Dalle3Context    int `json:"Dalle3Context"`
	Gpt4O1Context    int `json:"Gpt4o1Context"`
}

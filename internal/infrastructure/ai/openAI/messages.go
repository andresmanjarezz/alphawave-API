package openai

// 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

const GPT_MODEL = "gpt-3.5-turbo"

type OpenAiAPI struct {
	token string
}

func NewOpenAiAPI(token string) *OpenAiAPI {
	return &OpenAiAPI{
		token: token,
	}
}

type messageInput struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messagesInput struct {
	Model    string       `json:"model"`
	Messages messageInput `json:"messages"`
}

// func (o *OpenAiAPI) NewMessage(message string) []string {
// 	client := &http.Client{}

// 	input := messagesInput{
// 		Model:    GPT_MODEL,
// 		Messages: []messageInput{},
// 	}
// 	return []string{}
// }

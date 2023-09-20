package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

const GPT_MODEL = "gpt-3.5-turbo"

type OpenAiAPI struct {
	token string
	url   string
}

func NewOpenAiAPI(token string, url string) *OpenAiAPI {
	return &OpenAiAPI{
		token: token,
		url:   url,
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messagesInput struct {
	Model    string          `json:"model"`
	Messages []types.Message `json:"messages"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Response struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	// Model   string             `json:"model"`
	Choices []messagesResponse `json:"choices"`
	Usage   Usage              `json:"usage"`
}

type messagesResponse struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type OutputMessage struct {
	Message Message
}

func (o *OpenAiAPI) NewMessage(messages []types.Message) (OutputMessage, error) {
	client := &http.Client{}

	input := messagesInput{
		Model:    GPT_MODEL,
		Messages: messages,
	}

	inputBytes, err := json.Marshal(input)

	if err != nil {
		return OutputMessage{}, err
	}

	req, err := http.NewRequest("POST", o.url, bytes.NewBuffer(inputBytes))

	if err != nil {
		return OutputMessage{}, errors.New("error creating request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", o.token))
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)

	if err != nil {
		return OutputMessage{}, errors.New("error response")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return OutputMessage{}, err
	}

	var output Response

	err = json.Unmarshal(body, &output)

	if err != nil {
		return OutputMessage{}, err
	}

	var content OutputMessage

	for _, messageItem := range output.Choices {
		content = OutputMessage{
			Message: messageItem.Message,
		}
	}

	return content, nil
}

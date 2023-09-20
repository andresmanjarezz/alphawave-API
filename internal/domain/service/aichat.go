package service

import (
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	openai "github.com/Coke15/AlphaWave-BackEnd/internal/infrastructure/ai/openAI"
)

type openAI interface {
	NewMessage(messages []types.Message) (openai.OutputMessage, error)
}

type AiChatService struct {
	openAI openAI
}

func NewAiChatService(openAI openAI) *AiChatService {
	return &AiChatService{
		openAI: openAI,
	}
}

func (s *AiChatService) NewMessage(messages []types.Message) (types.Message, error) {

	message, err := s.openAI.NewMessage(messages)

	if err != nil {
		return types.Message{}, err
	}

	return types.Message{
		Role:    message.Message.Role,
		Content: message.Message.Content,
	}, nil
}

package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ChatRequest struct {
	ApiKey       string
	Prompt       string
	SystemPrompt string
	Model        string
	FileContents string
}

type ChatResponse struct {
	Text string
}

func Chat(request ChatRequest) (ChatResponse, error) {
	client := openai.NewClient(option.WithAPIKey(request.ApiKey))
	ctx := context.Background()

	model := request.Model
	if model == "" {
		model = openai.ChatModelGPT4oMini
	}

	systemPrompt := request.SystemPrompt

	var messages []openai.ChatCompletionMessageParamUnion
	if systemPrompt != "" {
		messages = append(messages, openai.SystemMessage(systemPrompt))
	}

	messages = append(messages, openai.UserMessage(request.Prompt))

	if request.FileContents != "" {
		messages = append(messages, openai.UserMessage(fmt.Sprintf("```%s```", request.FileContents)))
	}

	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    model,
	}

	resp, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ChatResponse{}, err
	}

	return ChatResponse{
		Text: resp.Choices[0].Message.Content,
	}, nil
}

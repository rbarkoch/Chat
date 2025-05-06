package llm

import (
	"chat/clargs"
	"chat/config"
	"chat/utils"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

type ChatRequest struct {
	ApiKey       string
	Prompt       string
	SystemPrompt string
	Model        string
	WebSearch    bool
	FileContents []string
}

type ChatResponse struct {
	Text string
}

// Performs the chat request through the OpenAI API.
func Chat(request ChatRequest) (ChatResponse, error) {
	client := openai.NewClient(option.WithAPIKey(request.ApiKey))
	ctx := context.Background()

	// User prompt.
	var messages responses.ResponseInputParam
	messages = append(messages, createUserMessage(request.Prompt))

	// File contents.
	for _, fileContents := range request.FileContents {
		messages = append(messages, createUserMessage(fmt.Sprintf("```%s```", fileContents)))
	}

	// Web search.
	tools := []responses.ToolUnionParam{}
	if request.WebSearch {
		tools = append(tools, responses.ToolUnionParam{
			OfWebSearch: &responses.WebSearchToolParam{
				Type: responses.WebSearchToolTypeWebSearchPreview,
			},
		})
	}

	// Construct parameters.
	params := responses.ResponseNewParams{
		Model:        utils.StringOrDefault(request.Model, openai.ChatModelGPT4oMini),
		Instructions: param.Opt[string]{Value: request.SystemPrompt},
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: messages,
		},
		Tools: tools,
	}

	// Call API.
	client.Responses.New(ctx, params)
	resp, err := client.Responses.New(ctx, params)
	if err != nil {
		fmt.Printf("API error: %v\n", err)
		return ChatResponse{}, err
	}

	// Output response.
	return ChatResponse{
		Text: resp.OutputText(),
	}, nil
}

// Constructs a chat request from the command line arguments and configuration.
func ConstructChatRequest(commandLineArguments clargs.CommandLineArgs, configuration config.Configuration) (ChatRequest, error) {
	// Start by constructing the chat request from configuration.
	chatRequest := ChatRequest{}

	chatRequest.ApiKey = configuration.Get(config.ConfigKeyApiKey)
	chatRequest.Model = configuration.Get(config.ConfigKeyModel)
	chatRequest.SystemPrompt = configuration.Get(config.ConfigKeySystemPrompt)
	chatRequest.WebSearch, _ = strconv.ParseBool(configuration.Get(config.ConfigWebSearch))

	// Override with command line arguments if provided.
	chatRequest.ApiKey = utils.StringOrDefault(commandLineArguments.ApiKey, chatRequest.ApiKey)
	chatRequest.Model = utils.StringOrDefault(commandLineArguments.Model, chatRequest.Model)
	chatRequest.SystemPrompt = utils.StringOrDefault(commandLineArguments.SystemPrompt, chatRequest.SystemPrompt)
	chatRequest.WebSearch = commandLineArguments.WebSearch || chatRequest.WebSearch

	// Set the rest of the values.

	// Prompt
	chatRequest.Prompt = commandLineArguments.Prompt

	// Files.
	if commandLineArguments.File != "" {
		filePaths := strings.Split(commandLineArguments.File, ",")

		for _, filePath := range filePaths {
			filePath = strings.TrimSpace(filePath)
			fileContents, err := utils.ReadFile(filePath)
			if err != nil {
				return ChatRequest{}, fmt.Errorf("error reading file: %s", err)
			}

			chatRequest.FileContents = append(chatRequest.FileContents, fileContents)
		}
	}

	// Validate the chat request.
	if chatRequest.Prompt == "" {
		return ChatRequest{}, fmt.Errorf("prompt is required")
	}

	if chatRequest.ApiKey == "" {
		return ChatRequest{}, fmt.Errorf("API key is required")
	}

	return chatRequest, nil
}

// Helper method to create the OpenAI models for a user message.
func createUserMessage(text string) responses.ResponseInputItemUnionParam {
	return responses.ResponseInputItemUnionParam{
		OfMessage: &responses.EasyInputMessageParam{
			Role: responses.EasyInputMessageRoleUser,
			Content: responses.EasyInputMessageContentUnionParam{
				OfString: param.Opt[string]{Value: text},
			},
		},
	}
}

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	openai "github.com/sashabaranov/go-openai"
)

type Agent struct {
	client     *anthropic.Client
	getUserMsg func() (string, bool)
}

func NewAgent(client *anthropic.Client, getUserMsg func() (string, bool)) *Agent {
	return &Agent{
		client:     client,
		getUserMsg: getUserMsg,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	conversation := []anthropic.MessageParam{}

	fmt.Println("Chat with Claude (use 'ctrl-c' to quit)")

	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userIn, ok := a.getUserMsg()
		if !ok {
		}

		UserMsg := anthropic.NewUserMessage(anthropic.NewTextBlock(userIn))
		conversation = append(conversation, UserMsg)

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}

		conversation = append(conversation, message.ToParam())

		for _, content := range message.Content {
			switch content.Type {
			case "text":
				fmt.Printf("\u001b[93mClaude\u001b[0m: %s\n", content.Text)
			}
		}

	}

}

func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Haiku20241022,
		MaxTokens: int64(1024),
		Messages:  conversation,
	})

	return message, err
}

func startAnthropic() {
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)
	getUserMsg := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}

		return scanner.Text(), true
	}

	agent := NewAgent(&client, getUserMsg)

	ctx := context.Background()

	err := agent.Run(ctx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

}

func startChatGPT() {
	fmt.Println("Chat with ChatGPT (use 'ctrl-c' to quit)")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}

		client := openai.NewClient(os.Getenv("OPENAI_KEY"))
		res, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT4oLatest,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    "user",
						Content: msg,
					},
				},
			},
		)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
		fmt.Printf("\u001b[93mChatGPT\u001b[0m: %s\n", res.Choices[0].Message.Content)
	}

}

func main() {
	argWithoutProg := os.Args

	if len(argWithoutProg) > 1 {
		switch {
		case argWithoutProg[1] == "-anthropic":
			startAnthropic()
		case argWithoutProg[1] == "-chatgpt":
			startChatGPT()
		default:
			fmt.Println("Only can use this flags:")
			fmt.Println("-anthropic")
			fmt.Println("-chatgpt")
		}
	}

}

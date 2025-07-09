package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
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

	return nil
}

func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Haiku20241022,
		MaxTokens: int64(1024),
		Messages:  conversation,
	})

	return message, err
}

func main() {
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

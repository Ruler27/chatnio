package hunyuan

import (
	"chat/globals"
	"context"
	"fmt"
)

type ChatProps struct {
	Model    string
	Messages []globals.Message
}

func (c *ChatInstance) FormatMessages(messages []globals.Message) []globals.Message {
	var result []globals.Message
	for _, message := range messages {
		switch message.Role {
		case globals.System:
			result = append(result, globals.Message{Role: globals.User, Content: message.Content})
		case globals.Assistant, globals.User:
			bound := len(result) > 0 && result[len(result)-1].Role == message.Role
			if bound {
				result[len(result)-1].Content += message.Content
			} else {
				result = append(result, message)
			}
		default:
			result = append(result, message)
		}
	}

	return result
}

func (c *ChatInstance) CreateStreamChatRequest(props *ChatProps, callback globals.Hook) error {
	credential := NewCredential(c.GetSecretId(), c.GetSecretKey())
	client := NewInstance(c.GetAppId(), credential)
	channel, err := client.Chat(context.Background(), NewRequest(Stream, c.FormatMessages(props.Messages)))
	if err != nil {
		return fmt.Errorf("tencent hunyuan error: %+v", err)
	}

	for chunk := range channel {
		if chunk.Error.Code != 0 {
			fmt.Printf("tencent hunyuan error: %+v\n", chunk.Error)
			break
		}

		if err := callback(chunk.Choices[0].Delta.Content); err != nil {
			return err
		}
	}

	return nil
}

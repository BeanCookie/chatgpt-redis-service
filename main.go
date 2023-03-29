package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	redis "github.com/redis/go-redis/v9"
	openai "github.com/sashabaranov/go-openai"
)

const CHATGPT_CHANNEL = "chatGPT"

type Message struct {
	ChannelId string
	Request   string
}

func main() {
	client := openai.NewClient(os.Getenv("OPENAI_KEY"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"), // 没有密码，默认值
		DB:       0,                           // 默认DB 0
	})

	pubsub := rdb.Subscribe(context.Background(), CHATGPT_CHANNEL)
	defer pubsub.Close()

	for requestMsg := range pubsub.Channel() {
		fmt.Println(requestMsg.Channel, requestMsg.Payload)

		var request Message

		err := json.Unmarshal([]byte(requestMsg.Payload), &request)
		if err != nil {
			fmt.Printf("Unmarshalling error: %v\n", err)
			continue
		}
		// 使用完毕，记得关闭
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: request.Request,
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}

		chatGPTResponse := resp.Choices[0].Message.Content
		err = rdb.Publish(context.Background(), request.ChannelId, chatGPTResponse).Err()
		if err != nil {
			fmt.Printf("Publish error: %v\n", err)
			continue
		}
	}
}

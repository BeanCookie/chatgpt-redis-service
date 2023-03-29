## 介绍
因为网络原因国内无法直接访问OpenAI的，所以我这里使用了Redis的发布订阅与OpenAI通信

## 部署
- https://railway.app/

## 免费的Redis
- https://redis.com/try-free/

## 客户端代码

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/snowflake"
	redis "github.com/redis/go-redis/v9"
)

const CHATGPT_CHANNEL = "chatGPT"


type Message struct {
	ChannelId string
	Request   string
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})

	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println("Snowflake NewNode error:", err)
		return
	}

	channelId := node.Generate()

	// 定义一个Message结构体实例
	msg := Message{
		ChannelId: fmt.Sprintf("%d", channelId),
		Request:   "Hello",
	}

	// 将Message结构体序列化为JSON格式的字符串
	requestJson, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Json encode error:", err)
		return
	}

	err = rdb.Publish(context.Background(), CHATGPT_CHANNEL, requestJson).Err()
	if err != nil {
		fmt.Println("Publish error:", err)
		return
	}

	pubsub := rdb.Subscribe(context.Background(), msg.ChannelId)
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		fmt.Println(msg)
		break
	}
}

```
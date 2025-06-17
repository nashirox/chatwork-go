package chatwork_test

import (
	"context"
	"fmt"
	"log"

	"github.com/nashirox/chatwork-go"
)

func ExampleClient_SendMessage() {
	// クライアントの作成
	client := chatwork.New("YOUR_API_TOKEN")
	ctx := context.Background()

	// メッセージの送信
	resp, _, err := client.Messages.SendMessage(ctx, 123456, "Hello, ChatWork!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message sent with ID: %s\n", resp.MessageID)
}

func ExampleClient_GetRooms() {
	// クライアントの作成
	client := chatwork.New("YOUR_API_TOKEN")
	ctx := context.Background()

	// ルーム一覧の取得
	rooms, _, err := client.Rooms.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		fmt.Printf("Room: %s (ID: %d)\n", room.Name, room.RoomID)
	}
}

func ExampleClient_CreateTask() {
	// クライアントの作成
	client := chatwork.New("YOUR_API_TOKEN")
	ctx := context.Background()

	// タスクの作成
	params := &chatwork.TaskCreateParams{
		Body:  "プレゼン資料を作成する",
		ToIDs: []int{123456, 789012},
	}

	resp, _, err := client.Tasks.Create(ctx, 123456, params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created %d tasks\n", len(resp.TaskIDs))
}

func ExampleClient_GetMyInfo() {
	// クライアントの作成  
	client := chatwork.New("YOUR_API_TOKEN")
	ctx := context.Background()

	// 自分の情報を取得
	me, _, err := client.Me.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name: %s\n", me.Name)
	fmt.Printf("ChatWork ID: %s\n", me.ChatworkID)
}
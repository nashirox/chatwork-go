# chatwork-go

[![Go Reference](https://pkg.go.dev/badge/github.com/nashirox/chatwork-go.svg)](https://pkg.go.dev/github.com/nashirox/chatwork-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nashirox/chatwork-go)](https://goreportcard.com/report/github.com/nashirox/chatwork-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client library for the [ChatWork API](https://developer.chatwork.com/).

## Installation

```bash
go get github.com/nashirox/chatwork-go
```

## Quick Start

```go
import "github.com/nashirox/chatwork-go"

// Create a new client with your API token
client := chatwork.New("YOUR_API_TOKEN")

// Get your rooms
rooms, _, err := client.Rooms.List(context.Background())
if err != nil {
    log.Fatal(err)
}

// Send a message
resp, _, err := client.Messages.SendMessage(context.Background(), roomID, "Hello, ChatWork!")
```

## Usage Examples

### Authentication

The library uses API token authentication. You can obtain your API token from your ChatWork account settings.

```go
client := chatwork.New("YOUR_API_TOKEN")
```

### Rooms

```go
ctx := context.Background()

// List all rooms
rooms, _, err := client.Rooms.List(ctx)

// Get a specific room
room, _, err := client.Rooms.Get(ctx, roomID)

// Create a new room
params := &chatwork.RoomCreateParams{
    Name:             "Project Room",
    Description:      "Discussion about the project",
    MembersAdminIDs:  []int{123456},
    MembersMemberIDs: []int{789012, 345678},
}
room, _, err := client.Rooms.Create(ctx, params)

// Update room information
updateParams := &chatwork.RoomUpdateParams{
    Name:        "Updated Project Room",
    Description: "Updated description",
}
room, _, err = client.Rooms.Update(ctx, roomID, updateParams)

// Leave a room
_, err = client.Rooms.Leave(ctx, roomID)

// Delete a room (creator only)
_, err = client.Rooms.DeleteRoom(ctx, roomID)
```

### Messages

```go
// Send a simple message
resp, _, err := client.Messages.SendMessage(ctx, roomID, "Hello, World!")

// Send a message with mentions
accountIDs := []int{123456, 789012}
resp, _, err = client.Messages.SendTo(ctx, roomID, accountIDs, "Hello team!")

// Reply to a message
resp, _, err = client.Messages.Reply(ctx, roomID, messageID, "Thanks for the message!")

// Quote a message
resp, _, err = client.Messages.Quote(ctx, roomID, messageID, "I agree with this")

// Send an info message
resp, _, err = client.Messages.SendInfo(ctx, roomID, "Important Notice", "Meeting at 3 PM")

// Get messages in a room
messages, _, err := client.Messages.List(ctx, roomID, nil)

// Get a specific message
message, _, err := client.Messages.Get(ctx, roomID, messageID)

// Update a message
updateParams := &chatwork.MessageUpdateParams{
    Body: "Updated message content",
}
message, _, err = client.Messages.Update(ctx, roomID, messageID, updateParams)

// Delete a message
_, _, err = client.Messages.Delete(ctx, roomID, messageID)
```

### Tasks

```go
// Create a task
params := &chatwork.TaskCreateParams{
    Body:  "Complete the report",
    ToIDs: []int{123456},
}
taskResp, _, err := client.Tasks.Create(ctx, roomID, params)

// Create a task with deadline
deadline := time.Now().Add(24 * time.Hour).Unix()
taskResp, _, err = client.Tasks.CreateWithDeadline(ctx, roomID, "Submit by tomorrow", []int{123456}, deadline)

// Get task information
task, _, err := client.Tasks.Get(ctx, roomID, taskID)

// Complete a task
task, _, err = client.Tasks.Complete(ctx, roomID, taskID)

// Reopen a task
task, _, err = client.Tasks.Reopen(ctx, roomID, taskID)

// Get my tasks
myTasks, _, err := client.MyTasks.List(ctx, nil)

// Get only open tasks
openTasks, _, err := client.MyTasks.GetOpen(ctx)

// Get completed tasks
completedTasks, _, err := client.MyTasks.GetCompleted(ctx)
```

### User Information

```go
// Get my information
me, _, err := client.Me.Get(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Name: %s\n", me.Name)
fmt.Printf("ChatWork ID: %s\n", me.ChatworkID)

// Get my status
status, _, err := client.Me.GetStatus(ctx)
fmt.Printf("Unread messages: %d\n", status.UnreadNum)
fmt.Printf("Unread rooms: %d\n", status.UnreadRoomNum)
fmt.Printf("Unread mentions: %d\n", status.MentionNum)
fmt.Printf("Unread tasks: %d\n", status.MytaskNum)
```

### Contacts

```go
// Get all contacts
contacts, _, err := client.Contacts.List(ctx)

// Get incoming requests
requests, _, err := client.IncomingRequests.List(ctx)

// Approve a contact request
approval, _, err := client.IncomingRequests.Approve(ctx, requestID)

// Reject a contact request
_, err = client.IncomingRequests.Reject(ctx, requestID)
```

### Files

```go
// Get files in a room
files, _, err := client.Rooms.GetFiles(ctx, roomID, 0)

// Get files uploaded by a specific user
files, _, err = client.Rooms.GetFiles(ctx, roomID, accountID)

// Get file information with download URL
file, _, err := client.Rooms.GetFile(ctx, roomID, fileID, true)
```

### Error Handling

The library provides detailed error information through the `ErrorResponse` type:

```go
resp, _, err := client.Messages.SendMessage(ctx, roomID, "Hello")
if err != nil {
    if chatworkErr, ok := err.(*chatwork.ErrorResponse); ok {
        // ChatWork API error
        fmt.Printf("API Error: %v\n", chatworkErr.Errors)
        fmt.Printf("Status Code: %d\n", chatworkErr.Response.StatusCode)
    } else {
        // Other error (network, etc.)
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Custom HTTP Client

You can provide a custom HTTP client for advanced use cases:

```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    // Add custom transport, middleware, etc.
}

client := chatwork.New("YOUR_API_TOKEN", 
    chatwork.OptionHTTPClient(httpClient),
)
```

## API Coverage

- ✅ **Rooms**: List, Create, Get, Update, Delete, Leave
- ✅ **Room Members**: Get, Update
- ✅ **Messages**: List, Create, Get, Update, Delete
- ✅ **Message Features**: Mentions, Replies, Quotes, Info messages
- ✅ **Tasks**: Create, Get, Update status
- ✅ **My Tasks**: List, Filter by status
- ✅ **Me**: Get profile, Get status
- ✅ **Contacts**: List
- ✅ **Incoming Requests**: List, Approve, Reject
- ✅ **Files**: List, Get with download URL
- ✅ **Read Status**: Get unread count, Mark as read

## Rate Limiting

ChatWork API has rate limits. The library exposes rate limit information in the `Response` struct:

```go
messages, resp, err := client.Messages.List(ctx, roomID, nil)
if resp != nil {
    fmt.Printf("Rate Limit: %d\n", resp.RateLimit.Limit)
    fmt.Printf("Remaining: %d\n", resp.RateLimit.Remaining)
    fmt.Printf("Reset: %d\n", resp.RateLimit.Reset)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Testing

Run tests with:

```bash
go test ./...
```

For integration tests with the actual API, set your API token:

```bash
export CHATWORK_API_TOKEN="your_token"
go test -tags=integration ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [slack-go/slack](https://github.com/slack-go/slack)
- API documentation from [ChatWork Developer](https://developer.chatwork.com/)

## Support

- Create an [issue](https://github.com/nashirox/chatwork-go/issues) for bug reports or feature requests
- Check the [ChatWork API documentation](https://developer.chatwork.com/) for API details
- See [examples](examples/) directory for more usage examples
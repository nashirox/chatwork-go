package chatwork

import (
	"context"
	"fmt"
)

// MessagesService handles communication with the message related
// methods of the ChatWork API.
//
// ChatWork API docs: https://developer.chatwork.com/reference/messages
type MessagesService service

// MessageCreateParams represents the parameters for creating a new message.
type MessageCreateParams struct {
	Body       string `url:"body"`
	SelfUnread bool   `url:"self_unread,omitempty"`
}

// MessageUpdateParams represents the parameters for updating a message.
type MessageUpdateParams struct {
	Body string `url:"body"`
}

// MessageListParams represents the parameters for listing messages.
type MessageListParams struct {
	// Force retrieval of messages
	// 0: Get up to 100 most recent messages (default)
	// 1: Get messages older than the most recent 100
	Force int
}

// MessageCreatedResponse represents the response when a message is created.
type MessageCreatedResponse struct {
	// The ID of the created message
	MessageID string `json:"message_id"`
}

// List returns messages in the specified room.
//
// By default, returns up to 100 most recent messages.
// Use params.Force = 1 to retrieve older messages.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id-messages
func (s *MessagesService) List(ctx context.Context, roomID int, params *MessageListParams) ([]*Message, *Response, error) {
	u := fmt.Sprintf("rooms/%d/messages", roomID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add URL parameters
	if params != nil && params.Force == 1 {
		q := req.URL.Query()
		q.Add("force", "1")
		req.URL.RawQuery = q.Encode()
	}

	var messages []*Message
	resp, err := s.client.Do(ctx, req, &messages)
	if err != nil {
		return nil, resp, err
	}

	return messages, resp, nil
}

// Create posts a new message to the specified room.
//
// The message body supports ChatWork message notation for mentions, quotes, etc.
//
// ChatWork API docs: https://developer.chatwork.com/reference/post-rooms-room_id-messages
func (s *MessagesService) Create(ctx context.Context, roomID int, params *MessageCreateParams) (*MessageCreatedResponse, *Response, error) {
	u := fmt.Sprintf("rooms/%d/messages", roomID)
	req, err := s.client.NewFormRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	result := new(MessageCreatedResponse)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get returns information about a specific message.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id-messages-message_id
func (s *MessagesService) Get(ctx context.Context, roomID int, messageID string) (*Message, *Response, error) {
	u := fmt.Sprintf("rooms/%d/messages/%s", roomID, messageID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	message := new(Message)
	resp, err := s.client.Do(ctx, req, message)
	if err != nil {
		return nil, resp, err
	}

	return message, resp, nil
}

// Update updates the specified message.
//
// Only the message creator can update their own messages.
// Messages can only be updated for a limited time after creation.
//
// ChatWork API docs: https://developer.chatwork.com/reference/put-rooms-room_id-messages-message_id
func (s *MessagesService) Update(ctx context.Context, roomID int, messageID string, params *MessageUpdateParams) (*Message, *Response, error) {
	u := fmt.Sprintf("rooms/%d/messages/%s", roomID, messageID)
	req, err := s.client.NewFormRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	message := new(Message)
	resp, err := s.client.Do(ctx, req, message)
	if err != nil {
		return nil, resp, err
	}

	return message, resp, nil
}

// Delete deletes the specified message.
//
// Only the message creator can delete their own messages.
// Messages can only be deleted for a limited time after creation.
//
// ChatWork API docs: https://developer.chatwork.com/reference/delete-rooms-room_id-messages-message_id
func (s *MessagesService) Delete(ctx context.Context, roomID int, messageID string) (*Message, *Response, error) {
	u := fmt.Sprintf("rooms/%d/messages/%s", roomID, messageID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}

	message := new(Message)
	resp, err := s.client.Do(ctx, req, message)
	if err != nil {
		return nil, resp, err
	}

	return message, resp, nil
}

// SendMessage is a convenience method for sending a simple text message.
//
// This is equivalent to calling Create with a MessageCreateParams containing only the body.
func (s *MessagesService) SendMessage(ctx context.Context, roomID int, body string) (*MessageCreatedResponse, *Response, error) {
	params := &MessageCreateParams{
		Body: body,
	}
	return s.Create(ctx, roomID, params)
}

// SendTo sends a message with mentions to specified users.
//
// The message will include [To:accountID] tags for each specified user,
// which will trigger notifications for those users.
func (s *MessagesService) SendTo(ctx context.Context, roomID int, accountIDs []int, body string) (*MessageCreatedResponse, *Response, error) {
	// Build mention tags
	mentions := ""
	for _, id := range accountIDs {
		mentions += fmt.Sprintf("[To:%d] ", id)
	}
	
	params := &MessageCreateParams{
		Body: mentions + body,
	}
	return s.Create(ctx, roomID, params)
}

// Reply sends a reply to a specific message.
//
// This creates a threaded conversation by linking the new message to the original.
func (s *MessagesService) Reply(ctx context.Context, roomID int, messageID string, body string) (*MessageCreatedResponse, *Response, error) {
	params := &MessageCreateParams{
		Body: fmt.Sprintf("[rp aid=%s] %s", messageID, body),
	}
	return s.Create(ctx, roomID, params)
}

// Quote sends a message quoting another message.
//
// This fetches the original message and includes it in a quote block
// before the new message body.
func (s *MessagesService) Quote(ctx context.Context, roomID int, messageID string, body string) (*MessageCreatedResponse, *Response, error) {
	// First, fetch the message to quote
	message, _, err := s.Get(ctx, roomID, messageID)
	if err != nil {
		return nil, nil, err
	}

	// Build the message with quote format
	quotedBody := fmt.Sprintf("[qt][qtmeta aid=%d time=%d]%s[/qt]\n%s", 
		message.Account.AccountID, 
		message.SendTime,
		message.Body,
		body,
	)
	
	params := &MessageCreateParams{
		Body: quotedBody,
	}
	return s.Create(ctx, roomID, params)
}

// SendInfo sends an information message with a title.
//
// Information messages are displayed with special formatting to highlight
// important information.
func (s *MessagesService) SendInfo(ctx context.Context, roomID int, title, body string) (*MessageCreatedResponse, *Response, error) {
	infoBody := fmt.Sprintf("[info][title]%s[/title]%s[/info]", title, body)
	params := &MessageCreateParams{
		Body: infoBody,
	}
	return s.Create(ctx, roomID, params)
}

// GetUnreadCount returns the number of unread messages in a room.
//
// This is a convenience method that uses the Rooms service's GetMessagesUnreadCount.
func (s *MessagesService) GetUnreadCount(ctx context.Context, roomID int) (int, *Response, error) {
	// Use RoomsService's GetMessagesUnreadCount
	roomsService := (*RoomsService)(&s.client.common)
	result, resp, err := roomsService.GetMessagesUnreadCount(ctx, roomID)
	if err != nil {
		return 0, resp, err
	}
	
	if count, ok := result["unread_num"]; ok {
		return count, resp, nil
	}
	
	return 0, resp, nil
}

// MarkAsRead marks all messages up to the specified message as read.
//
// This is a convenience method that uses the Rooms service's MarkMessagesAsRead.
func (s *MessagesService) MarkAsRead(ctx context.Context, roomID int, messageID string) (*Response, error) {
	// Use RoomsService's MarkMessagesAsRead
	roomsService := (*RoomsService)(&s.client.common)
	_, resp, err := roomsService.MarkMessagesAsRead(ctx, roomID, messageID)
	return resp, err
}
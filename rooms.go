package chatwork

import (
	"context"
	"fmt"
	"strconv"
)

// RoomsService handles communication with the room related
// methods of the ChatWork API.
//
// ChatWork API docs: https://developer.chatwork.com/reference/rooms
type RoomsService service

// RoomCreateParams represents the parameters for creating a new room.
//
// Name is required. Other fields are optional.
// Members can be specified with different permission levels.
type RoomCreateParams struct {
	Name             string   `url:"name"`
	Description      string   `url:"description,omitempty"`
	IconPreset       string   `url:"icon_preset,omitempty"`
	MembersAdminIDs  []int    `url:"members_admin_ids,comma,omitempty"`
	MembersMemberIDs []int    `url:"members_member_ids,comma,omitempty"`
	MembersReadonlyIDs []int  `url:"members_readonly_ids,comma,omitempty"`
}

// RoomUpdateParams represents the parameters for updating a room.
//
// All fields are optional. Only fields with non-zero values will be updated.
type RoomUpdateParams struct {
	Name        string `url:"name,omitempty"`
	Description string `url:"description,omitempty"`
	IconPreset  string `url:"icon_preset,omitempty"`
}

// RoomMembersUpdateParams represents the parameters for updating room members.
//
// This replaces all members in the room with the specified member lists.
// Be careful to include all desired members, not just new ones.
type RoomMembersUpdateParams struct {
	MembersAdminIDs    []int `url:"members_admin_ids,comma,omitempty"`
	MembersMemberIDs   []int `url:"members_member_ids,comma,omitempty"`
	MembersReadonlyIDs []int `url:"members_readonly_ids,comma,omitempty"`
}

// List returns the list of all rooms the authenticated user participates in.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms
func (s *RoomsService) List(ctx context.Context) ([]*Room, *Response, error) {
	req, err := s.client.NewRequest("GET", "rooms", nil)
	if err != nil {
		return nil, nil, err
	}

	var rooms []*Room
	resp, err := s.client.Do(ctx, req, &rooms)
	if err != nil {
		return nil, resp, err
	}

	return rooms, resp, nil
}

// Create creates a new group chat room.
//
// The authenticated user will automatically become an admin of the created room.
//
// ChatWork API docs: https://developer.chatwork.com/reference/post-rooms
func (s *RoomsService) Create(ctx context.Context, params *RoomCreateParams) (*Room, *Response, error) {
	req, err := s.client.NewFormRequest("POST", "rooms", params)
	if err != nil {
		return nil, nil, err
	}

	room := new(Room)
	resp, err := s.client.Do(ctx, req, room)
	if err != nil {
		return nil, resp, err
	}

	return room, resp, nil
}

// Get returns information about the specified room.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id
func (s *RoomsService) Get(ctx context.Context, roomID int) (*Room, *Response, error) {
	u := fmt.Sprintf("rooms/%d", roomID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	room := new(Room)
	resp, err := s.client.Do(ctx, req, room)
	if err != nil {
		return nil, resp, err
	}

	return room, resp, nil
}

// Update updates the room information.
//
// Only room admins can update room information.
//
// ChatWork API docs: https://developer.chatwork.com/reference/put-rooms-room_id
func (s *RoomsService) Update(ctx context.Context, roomID int, params *RoomUpdateParams) (*Room, *Response, error) {
	u := fmt.Sprintf("rooms/%d", roomID)
	req, err := s.client.NewFormRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	room := new(Room)
	resp, err := s.client.Do(ctx, req, room)
	if err != nil {
		return nil, resp, err
	}

	return room, resp, nil
}

// Delete leaves or deletes a room based on the action type.
//
// actionType can be "leave" or "delete":
// - "leave": Leave the room (any member can do this)
// - "delete": Delete the room (only room creator can do this)
//
// ChatWork API docs: https://developer.chatwork.com/reference/delete-rooms-room_id
func (s *RoomsService) Delete(ctx context.Context, roomID int, actionType string) (*Response, error) {
	u := fmt.Sprintf("rooms/%d", roomID)
	
	params := struct {
		ActionType string `url:"action_type"`
	}{
		ActionType: actionType,
	}
	
	req, err := s.client.NewFormRequest("DELETE", u, params)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Leave leaves the specified room.
//
// This is a convenience method that calls Delete with actionType "leave".
func (s *RoomsService) Leave(ctx context.Context, roomID int) (*Response, error) {
	return s.Delete(ctx, roomID, "leave")
}

// DeleteRoom deletes the specified room.
//
// Only the room creator can delete a room.
// This is a convenience method that calls Delete with actionType "delete".
func (s *RoomsService) DeleteRoom(ctx context.Context, roomID int) (*Response, error) {
	return s.Delete(ctx, roomID, "delete")
}

// GetMembers returns the list of all members in the specified room.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id-members
func (s *RoomsService) GetMembers(ctx context.Context, roomID int) ([]*Member, *Response, error) {
	u := fmt.Sprintf("rooms/%d/members", roomID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var members []*Member
	resp, err := s.client.Do(ctx, req, &members)
	if err != nil {
		return nil, resp, err
	}

	return members, resp, nil
}

// UpdateMembers updates the members of a room.
//
// This replaces all members in the room. Be sure to include all desired members.
// Only room admins can update members.
//
// ChatWork API docs: https://developer.chatwork.com/reference/put-rooms-room_id-members
func (s *RoomsService) UpdateMembers(ctx context.Context, roomID int, params *RoomMembersUpdateParams) (*Member, *Response, error) {
	u := fmt.Sprintf("rooms/%d/members", roomID)
	req, err := s.client.NewFormRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	member := new(Member)
	resp, err := s.client.Do(ctx, req, member)
	if err != nil {
		return nil, resp, err
	}

	return member, resp, nil
}

// GetMessagesReadStatus returns the read/unread status of a message.
//
// The response is a map with "unread_num" and "mention_num" keys.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id-messages-read
func (s *RoomsService) GetMessagesReadStatus(ctx context.Context, roomID int, messageID string) (map[string]int, *Response, error) {
	u := fmt.Sprintf("rooms/%d/messages/read", roomID)
	req, err := s.client.NewRequest("GET", u+"?message_id="+messageID, nil)
	if err != nil {
		return nil, nil, err
	}

	var status map[string]int
	resp, err := s.client.Do(ctx, req, &status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, nil
}

// MarkMessagesAsRead marks messages as read up to the specified message.
//
// All messages up to and including the specified message will be marked as read.
//
// ChatWork API docs: https://developer.chatwork.com/reference/put-rooms-room_id-messages-read
func (s *RoomsService) MarkMessagesAsRead(ctx context.Context, roomID int, messageID string) (map[string]string, *Response, error) {
	u := fmt.Sprintf("rooms/%d/messages/read", roomID)
	
	params := struct {
		MessageID string `url:"message_id"`
	}{
		MessageID: messageID,
	}
	
	req, err := s.client.NewFormRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	var result map[string]string
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetMessagesUnreadCount returns the number of unread messages in a room.
//
// The response includes "unread_num" and "mention_num".
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id-messages-unread
func (s *RoomsService) GetMessagesUnreadCount(ctx context.Context, roomID int) (map[string]int, *Response, error) {
	u := fmt.Sprintf("rooms/%d/messages/unread", roomID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var count map[string]int
	resp, err := s.client.Do(ctx, req, &count)
	if err != nil {
		return nil, resp, err
	}

	return count, resp, nil
}

// GetFiles returns the list of files in a room.
//
// If accountID is specified (non-zero), only files uploaded by that user are returned.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id-files
func (s *RoomsService) GetFiles(ctx context.Context, roomID int, accountID int) ([]*File, *Response, error) {
	u := fmt.Sprintf("rooms/%d/files", roomID)
	if accountID > 0 {
		u += "?account_id=" + strconv.Itoa(accountID)
	}
	
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var files []*File
	resp, err := s.client.Do(ctx, req, &files)
	if err != nil {
		return nil, resp, err
	}

	return files, resp, nil
}

// GetFile returns information about a specific file.
//
// If createDownloadURL is true, a download URL will be included in the response.
// Download URLs are valid for a limited time.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id-files-file_id
func (s *RoomsService) GetFile(ctx context.Context, roomID int, fileID int, createDownloadURL bool) (*File, *Response, error) {
	u := fmt.Sprintf("rooms/%d/files/%d", roomID, fileID)
	if createDownloadURL {
		u += "?create_download_url=1"
	}
	
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	file := new(File)
	resp, err := s.client.Do(ctx, req, file)
	if err != nil {
		return nil, resp, err
	}

	return file, resp, nil
}

// GetTasks returns the list of tasks in a room.
//
// Tasks can be filtered by various parameters.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id-tasks
func (s *RoomsService) GetTasks(ctx context.Context, roomID int, params *TaskListParams) ([]*Task, *Response, error) {
	u := fmt.Sprintf("rooms/%d/tasks", roomID)
	
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add URL parameters
	if params != nil {
		q := req.URL.Query()
		if params.AccountID > 0 {
			q.Add("account_id", strconv.Itoa(params.AccountID))
		}
		if params.AssignedByAccountID > 0 {
			q.Add("assigned_by_account_id", strconv.Itoa(params.AssignedByAccountID))
		}
		if params.Status != "" {
			q.Add("status", params.Status)
		}
		req.URL.RawQuery = q.Encode()
	}

	var tasks []*Task
	resp, err := s.client.Do(ctx, req, &tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, nil
}

// TaskListParams represents optional parameters for listing tasks.
type TaskListParams struct {
	// Filter by the account ID of the task assignee
	AccountID int

	// Filter by the account ID of the task creator
	AssignedByAccountID int

	// Filter by task status: "open" or "done"
	Status string
}
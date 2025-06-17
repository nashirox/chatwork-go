package chatwork

import "time"

// Room represents a ChatWork room (chat room).
//
// Rooms are the primary organizational unit in ChatWork.
// They can be either group chats, direct messages, or task-specific rooms.
type Room struct {
	RoomID         int    `json:"room_id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Role           string `json:"role"`
	Sticky         bool   `json:"sticky"`
	UnreadNum      int    `json:"unread_num"`
	MentionNum     int    `json:"mention_num"`
	MytaskNum      int    `json:"mytask_num"`
	MessageNum     int    `json:"message_num"`
	FileNum        int    `json:"file_num"`
	TaskNum        int    `json:"task_num"`
	IconPath       string `json:"icon_path"`
	LastUpdateTime int64  `json:"last_update_time"`
	Description    string `json:"description,omitempty"`
}

// Message represents a message in a ChatWork room.
//
// Messages are the primary communication unit in ChatWork.
// They can contain text, mentions, quotes, and attachments.
type Message struct {
	MessageID  string `json:"message_id"`
	Account    User   `json:"account"`
	Body       string `json:"body"`
	SendTime   int64  `json:"send_time"`
	UpdateTime int64  `json:"update_time"`
}

// User represents a ChatWork user account.
//
// This type contains both basic user information (like name and avatar)
// and detailed profile information when available.
type User struct {
	AccountID        int    `json:"account_id"`
	RoomID           int    `json:"room_id,omitempty"`
	Name             string `json:"name"`
	AvatarImageURL   string `json:"avatar_image_url"`
	ChatworkID       string `json:"chatwork_id,omitempty"`
	OrganizationID   int    `json:"organization_id,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	Department       string `json:"department,omitempty"`
	Title            string `json:"title,omitempty"`
	URL              string `json:"url,omitempty"`
	Introduction     string `json:"introduction,omitempty"`
	Mail             string `json:"mail,omitempty"`
	TelOrganization  string `json:"tel_organization,omitempty"`
	TelExtension     string `json:"tel_extension,omitempty"`
	TelMobile        string `json:"tel_mobile,omitempty"`
	Skype            string `json:"skype,omitempty"`
	Facebook         string `json:"facebook,omitempty"`
	Twitter          string `json:"twitter,omitempty"`
}

// Task represents a task assigned in a ChatWork room.
//
// Tasks are used to track work items and responsibilities.
// They can have assignees, due dates, and completion status.
type Task struct {
	TaskID            int    `json:"task_id"`
	Account           User   `json:"account"`
	AssignedByAccount User   `json:"assigned_by_account"`
	MessageID         string `json:"message_id"`
	Body              string `json:"body"`
	LimitTime         int64  `json:"limit_time"`
	Status            string `json:"status"`
	LimitType         string `json:"limit_type"`
}

// MyTask represents a task assigned to the authenticated user.
//
// This type includes additional room information compared to the regular Task type,
// making it easier to see tasks across multiple rooms.
type MyTask struct {
	TaskID            int         `json:"task_id"`
	Room              TaskRoom    `json:"room"`
	AssignedByAccount TaskAccount `json:"assigned_by_account"`
	MessageID         string      `json:"message_id"`
	Body              string      `json:"body"`
	LimitTime         int64       `json:"limit_time"`
	Status            string      `json:"status"`
	LimitType         string      `json:"limit_type"`
}

// TaskRoom represents minimal room information associated with a task.
// This is used in MyTask responses to provide context about where the task exists.
type TaskRoom struct {
	RoomID   int    `json:"room_id"`
	Name     string `json:"name"`
	IconPath string `json:"icon_path"`
}

// TaskAccount represents minimal account information for task assignees.
// This is used in task responses to identify who assigned or is assigned to tasks.
type TaskAccount struct {
	AccountID      int    `json:"account_id"`
	Name           string `json:"name"`
	AvatarImageURL string `json:"avatar_image_url"`
}

// Contact represents a user in your ChatWork contacts list.
//
// Contacts are users you have connected with on ChatWork.
// This type contains their basic information and organization details.
type Contact struct {
	AccountID        int    `json:"account_id"`
	RoomID           int    `json:"room_id"`
	Name             string `json:"name"`
	ChatworkID       string `json:"chatwork_id"`
	OrganizationID   int    `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Department       string `json:"department"`
	AvatarImageURL   string `json:"avatar_image_url"`
}

// Me represents the authenticated user's detailed information.
//
// This type contains all available information about the current user,
// including private details like email addresses that are not visible for other users.
type Me struct {
	AccountID        int    `json:"account_id"`
	RoomID           int    `json:"room_id"`
	Name             string `json:"name"`
	ChatworkID       string `json:"chatwork_id"`
	OrganizationID   int    `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Department       string `json:"department"`
	Title            string `json:"title"`
	URL              string `json:"url"`
	Introduction     string `json:"introduction"`
	Mail             string `json:"mail"`
	TelOrganization  string `json:"tel_organization"`
	TelExtension     string `json:"tel_extension"`
	TelMobile        string `json:"tel_mobile"`
	Skype            string `json:"skype"`
	Facebook         string `json:"facebook"`
	Twitter          string `json:"twitter"`
	AvatarImageURL   string `json:"avatar_image_url"`
	LoginMail        string `json:"login_mail"`
}

// MyStatus represents the authenticated user's unread counts and status.
//
// This provides a quick overview of pending items that need attention,
// including unread messages, mentions, and assigned tasks.
type MyStatus struct {
	UnreadRoomNum  int `json:"unread_room_num"`
	MentionRoomNum int `json:"mention_room_num"`
	MytaskRoomNum  int `json:"mytask_room_num"`
	UnreadNum      int `json:"unread_num"`
	MentionNum     int `json:"mention_num"`
	MytaskNum      int `json:"mytask_num"`
}

// File represents a file uploaded to a ChatWork room.
//
// Files can be images, documents, or any other type of attachment.
// Download URLs are only included when specifically requested.
type File struct {
	FileID      int    `json:"file_id"`
	Account     User   `json:"account"`
	MessageID   string `json:"message_id"`
	Filename    string `json:"filename"`
	Filesize    int    `json:"filesize"`
	UploadTime  int64  `json:"upload_time"`
	DownloadURL string `json:"download_url,omitempty"`
}

// Member represents a member of a ChatWork room.
//
// This includes their role in the room (admin, member, or readonly)
// and their basic account information.
type Member struct {
	AccountID        int    `json:"account_id"`
	Role             string `json:"role"`
	Name             string `json:"name"`
	ChatworkID       string `json:"chatwork_id"`
	OrganizationID   int    `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Department       string `json:"department"`
	AvatarImageURL   string `json:"avatar_image_url"`
}

// IncomingRequest represents a pending contact request.
//
// These are requests from other users to connect with you on ChatWork.
// You can approve or reject these requests.
type IncomingRequest struct {
	RequestID        int    `json:"request_id"`
	AccountID        int    `json:"account_id"`
	Message          string `json:"message"`
	Name             string `json:"name"`
	ChatworkID       string `json:"chatwork_id"`
	OrganizationID   int    `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Department       string `json:"department"`
	AvatarImageURL   string `json:"avatar_image_url"`
}

// Timestamp represents a Unix timestamp used throughout the ChatWork API.
//
// ChatWork uses Unix timestamps (seconds since epoch) for all time values.
// This type provides convenient conversion methods to Go's time.Time.
type Timestamp int64

// Time converts the Unix timestamp to a time.Time value.
func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// String returns the string representation of the timestamp.
// This is equivalent to calling Time().String().
func (t Timestamp) String() string {
	return t.Time().String()
}

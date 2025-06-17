package chatwork

import (
	"context"
	"fmt"
	"strconv"
)

// TasksService handles communication with the task related
// methods of the ChatWork API.
//
// ChatWork API docs: https://developer.chatwork.com/reference/tasks
type TasksService service

// TaskCreateParams represents the parameters for creating a new task.
type TaskCreateParams struct {
	// Task description (required)
	Body string `url:"body"`

	// Account IDs to assign the task to (required)
	ToIDs []int `url:"to_ids,comma"`

	// Task deadline as Unix timestamp (optional)
	Limit int64 `url:"limit,omitempty"`

	// Type of deadline: "none", "date", or "time" (optional)
	LimitType string `url:"limit_type,omitempty"`
}

// TaskCreatedResponse represents the response when tasks are created.
type TaskCreatedResponse struct {
	// IDs of the created tasks (one for each assignee)
	TaskIDs []int `json:"task_ids"`
}

// Create creates new tasks in the specified room.
//
// A separate task is created for each account ID specified in ToIDs.
//
// ChatWork API docs: https://developer.chatwork.com/reference/post-rooms-room_id-tasks
func (s *TasksService) Create(ctx context.Context, roomID int, params *TaskCreateParams) (*TaskCreatedResponse, *Response, error) {
	u := fmt.Sprintf("rooms/%d/tasks", roomID)
	req, err := s.client.NewFormRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	result := new(TaskCreatedResponse)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get returns information about a specific task.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-rooms-room_id-tasks-task_id
func (s *TasksService) Get(ctx context.Context, roomID int, taskID int) (*Task, *Response, error) {
	u := fmt.Sprintf("rooms/%d/tasks/%d", roomID, taskID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	task := new(Task)
	resp, err := s.client.Do(ctx, req, task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}

// UpdateStatus updates the status of a task.
//
// Status can be "open" or "done".
//
// ChatWork API docs: https://developer.chatwork.com/reference/put-rooms-room_id-tasks-task_id-status
func (s *TasksService) UpdateStatus(ctx context.Context, roomID int, taskID int, status string) (*Task, *Response, error) {
	u := fmt.Sprintf("rooms/%d/tasks/%d/status", roomID, taskID)
	
	params := struct {
		Body string `url:"body"`
	}{
		Body: status,
	}
	
	req, err := s.client.NewFormRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	task := new(Task)
	resp, err := s.client.Do(ctx, req, task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}

// Complete marks a task as completed.
//
// This is a convenience method that calls UpdateStatus with status "done".
func (s *TasksService) Complete(ctx context.Context, roomID int, taskID int) (*Task, *Response, error) {
	return s.UpdateStatus(ctx, roomID, taskID, "done")
}

// Reopen marks a task as open (not completed).
//
// This is a convenience method that calls UpdateStatus with status "open".
func (s *TasksService) Reopen(ctx context.Context, roomID int, taskID int) (*Task, *Response, error) {
	return s.UpdateStatus(ctx, roomID, taskID, "open")
}

// CreateSimple is a convenience method for creating a task without a deadline.
func (s *TasksService) CreateSimple(ctx context.Context, roomID int, body string, toIDs []int) (*TaskCreatedResponse, *Response, error) {
	params := &TaskCreateParams{
		Body:  body,
		ToIDs: toIDs,
	}
	return s.Create(ctx, roomID, params)
}

// CreateWithDeadline is a convenience method for creating a task with a deadline.
//
// The deadline should be a Unix timestamp.
func (s *TasksService) CreateWithDeadline(ctx context.Context, roomID int, body string, toIDs []int, deadline int64) (*TaskCreatedResponse, *Response, error) {
	params := &TaskCreateParams{
		Body:      body,
		ToIDs:     toIDs,
		Limit:     deadline,
		LimitType: "time",
	}
	return s.Create(ctx, roomID, params)
}

// MyTasksService handles communication with the "my tasks" related
// methods of the ChatWork API.
//
// ChatWork API docs: https://developer.chatwork.com/reference/my-tasks
type MyTasksService service

// MyTaskListParams represents optional parameters for listing my tasks.
type MyTaskListParams struct {
	// Filter by the account ID of who assigned the task
	AssignedByAccountID int

	// Filter by task status: "open" or "done"
	Status string
}

// List returns all tasks assigned to the authenticated user.
//
// Tasks can be filtered by status and who assigned them.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-my-tasks
func (s *MyTasksService) List(ctx context.Context, params *MyTaskListParams) ([]*MyTask, *Response, error) {
	req, err := s.client.NewRequest("GET", "my/tasks", nil)
	if err != nil {
		return nil, nil, err
	}

	// Add URL parameters
	if params != nil {
		q := req.URL.Query()
		if params.AssignedByAccountID > 0 {
			q.Add("assigned_by_account_id", strconv.Itoa(params.AssignedByAccountID))
		}
		if params.Status != "" {
			q.Add("status", params.Status)
		}
		req.URL.RawQuery = q.Encode()
	}

	var tasks []*MyTask
	resp, err := s.client.Do(ctx, req, &tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, nil
}

// GetOpen returns all open (uncompleted) tasks assigned to the authenticated user.
//
// This is a convenience method that calls List with status "open".
func (s *MyTasksService) GetOpen(ctx context.Context) ([]*MyTask, *Response, error) {
	params := &MyTaskListParams{
		Status: "open",
	}
	return s.List(ctx, params)
}

// GetCompleted returns all completed tasks assigned to the authenticated user.
//
// This is a convenience method that calls List with status "done".
func (s *MyTasksService) GetCompleted(ctx context.Context) ([]*MyTask, *Response, error) {
	params := &MyTaskListParams{
		Status: "done",
	}
	return s.List(ctx, params)
}

// GetByRoom returns tasks assigned to the authenticated user in a specific room.
//
// This fetches all tasks and filters them by room ID locally.
func (s *MyTasksService) GetByRoom(ctx context.Context, roomID int) ([]*MyTask, *Response, error) {
	allTasks, resp, err := s.List(ctx, nil)
	if err != nil {
		return nil, resp, err
	}

	var tasks []*MyTask
	for _, task := range allTasks {
		if task.Room.RoomID == roomID {
			tasks = append(tasks, task)
		}
	}

	return tasks, resp, nil
}

// CompleteTask marks a task as completed.
//
// This is a convenience method that uses the Tasks service to complete a task.
func (s *MyTasksService) CompleteTask(ctx context.Context, roomID int, taskID int) (*Task, *Response, error) {
	tasksService := (*TasksService)(&s.client.common)
	return tasksService.Complete(ctx, roomID, taskID)
}

// ReopenTask marks a task as open (not completed).
//
// This is a convenience method that uses the Tasks service to reopen a task.
func (s *MyTasksService) ReopenTask(ctx context.Context, roomID int, taskID int) (*Task, *Response, error) {
	tasksService := (*TasksService)(&s.client.common)
	return tasksService.Reopen(ctx, roomID, taskID)
}
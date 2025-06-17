package chatwork

import (
	"context"
)

// MeService handles communication with the "me" related
// methods of the ChatWork API.
//
// ChatWork API docs: https://developer.chatwork.com/reference/me
type MeService service

// Get returns detailed information about the authenticated user.
//
// This includes private information like email addresses that are not
// visible when viewing other users' profiles.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-me
func (s *MeService) Get(ctx context.Context) (*Me, *Response, error) {
	req, err := s.client.NewRequest("GET", "me", nil)
	if err != nil {
		return nil, nil, err
	}

	me := new(Me)
	resp, err := s.client.Do(ctx, req, me)
	if err != nil {
		return nil, resp, err
	}

	return me, resp, nil
}

// GetStatus returns the authenticated user's unread counts.
//
// This includes counts for:
// - Unread messages
// - Unread rooms
// - Mentions
// - Tasks
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-my-status
func (s *MeService) GetStatus(ctx context.Context) (*MyStatus, *Response, error) {
	req, err := s.client.NewRequest("GET", "my/status", nil)
	if err != nil {
		return nil, nil, err
	}

	status := new(MyStatus)
	resp, err := s.client.Do(ctx, req, status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, nil
}
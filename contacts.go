package chatwork

import (
	"context"
	"strconv"
)

// ContactsService handles communication with the contacts related
// methods of the ChatWork API.
//
// ChatWork API docs: https://developer.chatwork.com/reference/contacts
type ContactsService service

// List returns all contacts of the authenticated user.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-contacts
func (s *ContactsService) List(ctx context.Context) ([]*Contact, *Response, error) {
	req, err := s.client.NewRequest("GET", "contacts", nil)
	if err != nil {
		return nil, nil, err
	}

	var contacts []*Contact
	resp, err := s.client.Do(ctx, req, &contacts)
	if err != nil {
		return nil, resp, err
	}

	return contacts, resp, nil
}

// IncomingRequestsService handles communication with the incoming requests related
// methods of the ChatWork API.
//
// ChatWork API docs: https://developer.chatwork.com/reference/incoming-requests
type IncomingRequestsService service

// List returns all pending contact requests.
//
// ChatWork API docs: https://developer.chatwork.com/reference/get-incoming_requests
func (s *IncomingRequestsService) List(ctx context.Context) ([]*IncomingRequest, *Response, error) {
	req, err := s.client.NewRequest("GET", "incoming_requests", nil)
	if err != nil {
		return nil, nil, err
	}

	var requests []*IncomingRequest
	resp, err := s.client.Do(ctx, req, &requests)
	if err != nil {
		return nil, resp, err
	}

	return requests, resp, nil
}

// IncomingRequestActionResponse represents the response when approving a contact request.
type IncomingRequestActionResponse struct {
	AccountID        int    `json:"account_id"`
	RoomID           int    `json:"room_id"`
	Name             string `json:"name"`
	ChatworkID       string `json:"chatwork_id"`
	OrganizationID   int    `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Department       string `json:"department"`
	AvatarImageURL   string `json:"avatar_image_url"`
}

// Approve approves a contact request.
//
// ChatWork API docs: https://developer.chatwork.com/reference/put-incoming_requests-request_id
func (s *IncomingRequestsService) Approve(ctx context.Context, requestID int) (*IncomingRequestActionResponse, *Response, error) {
	u := "incoming_requests/" + strconv.Itoa(requestID)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(IncomingRequestActionResponse)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Reject rejects a contact request.
//
// ChatWork API docs: https://developer.chatwork.com/reference/delete-incoming_requests-request_id
func (s *IncomingRequestsService) Reject(ctx context.Context, requestID int) (*Response, error) {
	u := "incoming_requests/" + strconv.Itoa(requestID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
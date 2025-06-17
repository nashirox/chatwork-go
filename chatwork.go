// Package chatwork provides a client library for the ChatWork API.
//
// The ChatWork API is a RESTful API for ChatWork, a group chat service.
// This library provides a simple and idiomatic way to interact with the API from Go applications.
//
// Usage
//
// Import the package:
//
//	import "github.com/nashirox/chatwork-go"
//
// Create a new client with your API token:
//
//	client := chatwork.New("YOUR_API_TOKEN")
//
// Use the client to interact with various endpoints:
//
//	// Get rooms
//	rooms, _, err := client.Rooms.List(context.Background())
//
//	// Send a message
//	resp, _, err := client.Messages.SendMessage(context.Background(), roomID, "Hello!")
//
// For more detailed examples, see the README and example files.
package chatwork

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://api.chatwork.com/v2"
	userAgent      = "chatwork-go"
)

// Client manages communication with the ChatWork API.
//
// The Client is the main entry point for interacting with the ChatWork API.
// It provides access to various API endpoints through service properties.
// The zero value is not usable; use New to create a configured client.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests. Defaults to the public ChatWork API.
	BaseURL *url.URL

	// User agent used when communicating with the ChatWork API.
	UserAgent string

	// API token for authentication.
	token string

	// Services used for talking to different parts of the ChatWork API.
	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services for each endpoint
	Rooms            *RoomsService
	Messages         *MessagesService
	Me               *MeService
	MyTasks          *MyTasksService
	Contacts         *ContactsService
	Tasks            *TasksService
	IncomingRequests *IncomingRequestsService
}

// service is the base type for all API endpoint services.
// It holds a reference to the client for making HTTP requests.
type service struct {
	client *Client
}

// New creates a new ChatWork API client with the provided API token.
//
// The API token can be obtained from your ChatWork account settings.
// You can provide additional options to customize the client behavior.
//
// Example:
//
//	client := chatwork.New("YOUR_API_TOKEN")
//
//	// With custom HTTP client
//	httpClient := &http.Client{Timeout: 30 * time.Second}
//	client := chatwork.New("YOUR_API_TOKEN", chatwork.OptionHTTPClient(httpClient))
func New(token string, options ...ClientOption) *Client {
	httpClient := &http.Client{}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
		token:     token,
	}

	c.common.client = c
	c.Rooms = (*RoomsService)(&c.common)
	c.Messages = (*MessagesService)(&c.common)
	c.Me = (*MeService)(&c.common)
	c.MyTasks = (*MyTasksService)(&c.common)
	c.Contacts = (*ContactsService)(&c.common)
	c.Tasks = (*TasksService)(&c.common)
	c.IncomingRequests = (*IncomingRequestsService)(&c.common)

	for _, option := range options {
		option(c)
	}

	return c
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// OptionHTTPClient sets a custom HTTP client to be used for API requests.
// This is useful for setting custom timeouts, transport settings, or middleware.
//
// Example:
//
//	httpClient := &http.Client{
//		Timeout: 30 * time.Second,
//	}
//	client := chatwork.New("token", chatwork.OptionHTTPClient(httpClient))
func OptionHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.client = httpClient
	}
}

// OptionDebug enables debug mode for the client.
// When enabled, the client will log detailed information about API requests and responses.
// This is useful for troubleshooting and development.
func OptionDebug(debug bool) ClientOption {
	return func(c *Client) {
		if debug {
			// TODO: Implement debug logging
		}
	}
}

// NewRequest creates a new API request with JSON body.
//
// The urlStr is relative to the BaseURL of the client.
// The body, if specified, is JSON encoded and included as the request body.
// The appropriate headers are automatically set.
//
// This method is primarily used internally by service methods,
// but can be used directly for making custom API requests.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-ChatWorkToken", c.token)

	return req, nil
}

// NewFormRequest creates a new API request with form-encoded body.
//
// This method is similar to NewRequest but encodes the body as form data
// instead of JSON. It's used for endpoints that expect form-encoded data.
//
// The body parameter should be a struct with url tags for encoding.
func (c *Client) NewFormRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.Reader
	if body != nil {
		form, err := query.Values(body)
		if err != nil {
			return nil, err
		}
		buf = strings.NewReader(form.Encode())
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-ChatWorkToken", c.token)

	return req, nil
}

// Do sends an API request and returns the API response.
//
// The API response is JSON decoded and stored in the value pointed to by v,
// or returned as an error if an API error has occurred. If v implements the
// io.Writer interface, the raw response body will be written to v, without
// attempting to first decode it.
//
// The provided context is used to cancel the request if needed.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

// Response is a ChatWork API response.
// This wraps the standard http.Response and provides convenient access to
// rate limit information and other ChatWork-specific response data.
type Response struct {
	*http.Response

	// Rate limit information parsed from headers
	RateLimit RateLimit
}

// RateLimit represents the rate limit information for the ChatWork API.
// ChatWork imposes rate limits on API requests to ensure fair usage.
type RateLimit struct {
	// The maximum number of requests that can be made in the rate limit window.
	Limit int

	// The number of requests remaining in the current rate limit window.
	Remaining int

	// The time at which the current rate limit window resets, in Unix timestamp.
	Reset int64
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	// TODO: Parse rate limit headers from response
	return response
}

// CheckResponse checks the API response for errors.
//
// ChatWork API returns non-2xx status codes to indicate errors.
// This function extracts error information from the response body
// and returns an appropriate error type.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}

// ErrorResponse represents an error response from the ChatWork API.
//
// ChatWork API returns error details in the response body when requests fail.
// This type captures those details and implements the error interface.
type ErrorResponse struct {
	// The HTTP response that caused this error
	Response *http.Response

	// Error messages returned by the API
	Errors []string `json:"errors"`
}

// Error returns a human-readable description of the API error.
// It includes the HTTP method, URL, status code, and error messages.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, strings.Join(r.Errors, ", "))
}

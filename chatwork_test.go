package chatwork

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const testToken = "test-token"

func TestNew(t *testing.T) {
	token := testToken
	client := New(token)

	if client.token != token {
		t.Errorf("Expected token %s, got %s", token, client.token)
	}

	if client.BaseURL.String() != defaultBaseURL {
		t.Errorf("Expected base URL %s, got %s", defaultBaseURL, client.BaseURL.String())
	}

	if client.UserAgent != userAgent {
		t.Errorf("Expected user agent %s, got %s", userAgent, client.UserAgent)
	}
}

func TestNewWithOptions(t *testing.T) {
	token := testToken
	customHTTPClient := &http.Client{}

	client := New(token, OptionHTTPClient(customHTTPClient))

	if client.client != customHTTPClient {
		t.Error("Custom HTTP client was not set")
	}
}

func TestNewRequest(t *testing.T) {
	client := New(testToken)

	req, err := client.NewRequest("GET", "test", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("Expected method GET, got %s", req.Method)
	}

	expectedURL := client.BaseURL.String() + "/test"
	if req.URL.String() != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
	}

	if req.Header.Get("X-ChatWorkToken") != testToken {
		t.Error("X-ChatWorkToken header not set correctly")
	}

	if req.Header.Get("User-Agent") != userAgent {
		t.Error("User-Agent header not set correctly")
	}
}

func TestCheckResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
	}{
		{
			name:       "success response",
			statusCode: http.StatusOK,
			body:       `{"success": true}`,
			wantErr:    false,
		},
		{
			name:       "error response",
			statusCode: http.StatusBadRequest,
			body:       `{"errors": ["Invalid request"]}`,
			wantErr:    true,
		},
		{
			name:       "no content",
			statusCode: http.StatusNoContent,
			body:       "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       http.NoBody,
				Request: &http.Request{
					Method: "GET",
					URL:    &url.URL{Path: "/test"},
				},
			}

			if tt.body != "" {
				resp.Body = httptest.NewRecorder().Result().Body
			}

			err := CheckResponse(resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Request: &http.Request{
			Method: "POST",
			URL:    &url.URL{Path: "/rooms"},
		},
	}

	err := &APIError{
		Response: resp,
		Errors:   []string{"Invalid parameters", "Room name is required"},
	}

	expected := "POST /rooms: 400 Invalid parameters, Room name is required"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

func TestTimestamp(t *testing.T) {
	ts := Timestamp(1609459200) // 2021-01-01 00:00:00 UTC

	timeStr := ts.String()
	if timeStr == "" {
		t.Error("Timestamp.String() returned empty string")
	}

	time := ts.Time()
	if time.Unix() != 1609459200 {
		t.Errorf("Expected Unix timestamp 1609459200, got %d", time.Unix())
	}
}

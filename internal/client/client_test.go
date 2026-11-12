package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestListFollowsCursorWithoutForwardingKeyOutsideCollection(t *testing.T) {
	const endpoint = "/api/0.1/organizations/test/security/code-projects/"
	for _, tt := range []struct {
		name, next string
		wantError  bool
	}{
		{"relative cursor", "?cursor=page2", false},
		{"absolute cursor", "https://argonix.test" + endpoint + "?cursor=page2", false},
		{"external URL", "https://another.test" + endpoint + "?cursor=page2", true},
		{"different collection", "https://argonix.test/api/0.1/something-else/", true},
		{"repeated link", "https://argonix.test" + endpoint, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			requests := 0
			c := Client{
				BaseURL: "https://argonix.test", APIKey: "fake", OrganizationID: "test",
				HTTPClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					requests++
					if req.URL.Host != "argonix.test" || req.URL.Path != endpoint {
						t.Fatalf("credentials sent to unexpected URL %s", req.URL)
					}
					if req.Header.Get("Authorization") != "Api-Key fake" {
						t.Fatal("missing API authentication")
					}
					body := fmt.Sprintf(`{"results":[{"name":"first"}],"next":%q}`, tt.next)
					if req.URL.Query().Get("cursor") == "page2" {
						body = `{"results":[{"name":"second"}],"next":null}`
					}
					return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, nil
				})},
			}
			var result []struct{ Name string }
			err := c.List(t.Context(), "/security/code-projects/", &result)
			if tt.wantError {
				if err == nil || requests != 1 {
					t.Fatalf("unsafe pagination got %v, requests %d", err, requests)
				}
				return
			}
			if err != nil || len(result) != 2 || result[0].Name != "first" || result[1].Name != "second" {
				t.Fatalf("expected both pages, got %+v: %v", result, err)
			}
		})
	}
}

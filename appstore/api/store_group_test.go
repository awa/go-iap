package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// newGroupTestClient points a StoreClient at a server that records the request
// and answers with status and body.
func newGroupTestClient(t *testing.T, status int, body string, gotReq **http.Request) *StoreClient {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*gotReq = r
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	return NewStoreClient(&StoreConfig{
		KeyContent: testPrivateKeyPEM(t),
		KeyID:      "TESTKEYID",
		BundleID:   "test.bundle.id",
		Issuer:     "test-issuer",
		HostDebug:  srv.URL,
	})
}

func TestGetCustomerGroups(t *testing.T) {
	t.Parallel()

	// Sample response from https://developer.apple.com/documentation/appstoreserverapi/get-customer-groups
	body := `{
		"groups": [
			{
				"groupId": "900000000000000000",
				"groupType": "ORGANIZATION",
				"roles": [
					{"productId": "example.productId.one", "role": "ADMIN"},
					{"productId": "example.productId.two", "role": "NONE"}
				]
			}
		]
	}`
	var req *http.Request
	cli := newGroupTestClient(t, http.StatusOK, body, &req)

	rsp, err := cli.GetCustomerGroups(context.Background(), "700000000000000000")
	if err != nil {
		t.Fatalf("GetCustomerGroups: %v", err)
	}

	if req.Method != http.MethodGet {
		t.Errorf("method = %s, want GET", req.Method)
	}
	if want := "/groups/v1/currentGroups/700000000000000000"; req.URL.Path != want {
		t.Errorf("path = %s, want %s", req.URL.Path, want)
	}
	want := &GetCustomerGroupsResponse{
		Groups: []GroupEntry{
			{
				GroupId:   "900000000000000000",
				GroupType: GroupTypeOrganization,
				Roles: []RoleEntry{
					{ProductId: "example.productId.one", Role: RoleAdmin},
					{ProductId: "example.productId.two", Role: RoleNone},
				},
			},
		},
	}
	if !reflect.DeepEqual(rsp, want) {
		t.Errorf("response = %+v, want %+v", rsp, want)
	}
}

func TestGetCustomerGroupsTransactionNotFound(t *testing.T) {
	t.Parallel()

	var req *http.Request
	cli := newGroupTestClient(t, http.StatusNotFound, `{"errorCode": 4040010, "errorMessage": "Transaction id not found."}`, &req)

	_, err := cli.GetCustomerGroups(context.Background(), "1")
	if !errors.Is(err, TransactionIdNotFoundError) {
		t.Errorf("error = %v, want TransactionIdNotFoundError", err)
	}
}

func TestGetGroupMembers(t *testing.T) {
	t.Parallel()

	// Sample response from https://developer.apple.com/documentation/appstoreserverapi/get-group-members
	body := `{
		"members": [{"appTransactionId": "700000000000000000"}],
		"paginationToken": "next-token",
		"hasMore": true
	}`

	tests := []struct {
		name            string
		paginationToken string
		limit           int
		wantQuery       string
	}{
		{name: "first page without limit", wantQuery: ""},
		{name: "first page with limit", limit: 50, wantQuery: "limit=50"},
		{name: "next page", paginationToken: "abc", limit: 50, wantQuery: "limit=50&paginationToken=abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var req *http.Request
			cli := newGroupTestClient(t, http.StatusOK, body, &req)

			rsp, err := cli.GetGroupMembers(context.Background(), "900000000000000000", tt.paginationToken, tt.limit)
			if err != nil {
				t.Fatalf("GetGroupMembers: %v", err)
			}

			if req.Method != http.MethodGet {
				t.Errorf("method = %s, want GET", req.Method)
			}
			if want := "/groups/v1/group/900000000000000000"; req.URL.Path != want {
				t.Errorf("path = %s, want %s", req.URL.Path, want)
			}
			if req.URL.RawQuery != tt.wantQuery {
				t.Errorf("query = %q, want %q", req.URL.RawQuery, tt.wantQuery)
			}
			want := &GetGroupMembersResponse{
				Members:         []GroupMemberEntry{{AppTransactionId: "700000000000000000"}},
				PaginationToken: "next-token",
				HasMore:         true,
			}
			if !reflect.DeepEqual(rsp, want) {
				t.Errorf("response = %+v, want %+v", rsp, want)
			}
		})
	}
}

func TestGetGroupMembersGroupNotFound(t *testing.T) {
	t.Parallel()

	var req *http.Request
	cli := newGroupTestClient(t, http.StatusNotFound, `{"errorCode": 4040022, "errorMessage": "Group not found."}`, &req)

	_, err := cli.GetGroupMembers(context.Background(), "1", "", 0)
	if !errors.Is(err, GroupNotFoundError) {
		t.Errorf("error = %v, want GroupNotFoundError", err)
	}
}

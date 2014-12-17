package yo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	// Test Yo Client.
	testClient *Client
	testToken  = "some_token"
	testUser   = "some_user"
)

func init() {
	testClient = NewClient(testToken)
	return
}

func TestYoAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(yoAllHandler))
	defer server.Close()
	YO_API = server.URL

	if err := testClient.YoAll(); err != nil {
		t.Fatal(err)
	}
}

func TestYoAllLink(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(yoAllHandler))
	defer server.Close()
	YO_API = server.URL

	if err := testClient.YoAllLink("http://google.com"); err != nil {
		t.Fatal(err)
	}
}

func yoAllHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusCreated)
}

func TestYoUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(yoUserHandler))
	defer server.Close()
	YO_API = server.URL

	if err := testClient.YoUser(testUser); err != nil {
		t.Fatal(err)
	}
}

func TestYoUserLink(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(yoUserHandler))
	defer server.Close()
	YO_API = server.URL

	if err := testClient.YoUserLink(testUser, "http://google.com"); err != nil {
		t.Fatal(err)
	}
}

func yoUserHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusCreated)
}

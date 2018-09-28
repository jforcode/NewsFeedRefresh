package newsApi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO test for errors. Testing for all the errors would be good but not necessary right now.
// probably can publish a test server for newsApi
// https://newsapi.org/docs/errors
func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Api-Key")
		if apiKey == "" {
			fmt.Fprintln(w, "Error")
		}
		fmt.Fprintln(w, "Done")
	}))
	defer server.Close()

	api := &Api{}
	api.Init(server.URL, "Test API Key")

	t.Run("valid api call", func(t *testing.T) {
		api.get("test", nil)
	})

	t.Run("api call without api key", func(t *testing.T) {

	})
}

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sad0vnikov/radish/config"
	"github.com/sad0vnikov/radish/http/responds"
	"github.com/sad0vnikov/radish/redis"
)

func TestGettingServersList(t *testing.T) {

	config.StubConfigLoader{}.Load()

	req := httptest.NewRequest("GET", "/v1/servers", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := GetServersList(w, r)
		responds.RespondJSON(w, resp)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned status code %v, expected %v", status, http.StatusOK)
	}

	expectedServers := map[string]redis.Server{}
	expectedServers["server1"] = redis.NewServer("server1", "127.0.0.1", 6379)
	expectedServers["server2"] = redis.NewServer("server2", "127.0.0.1", 6380)
	expectedServers["server3"] = redis.NewServer("server3", "127.0.0.1", 6381)

	expectedJSONBytes, err := json.Marshal(expectedServers)
	if err != nil {
		t.Fatal("error serialising expecting values to JSON")
	}

	expectedJSON := string(expectedJSONBytes)

	resultJSON := rr.Body.String()

	if expectedJSON != resultJSON {
		t.Errorf("got invalid json %v, expected %v", resultJSON, expectedJSON)
	}

}

func TestGettingKeys(t *testing.T) {
	config.StubConfigLoader{}.Load()

}

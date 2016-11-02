package auth0_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vouchedfor/auth0"
)

func TestApi_Post(t *testing.T) {
	body := struct {
		TestField string `json:"test_field"`
	}{
		TestField: "some data",
	}

	endpointUrl := "/test/url"

	apiServer := httptest.NewServer(http.HandlerFunc(mockApiHandler(t, endpointUrl)))
	defer apiServer.Close()

	api := auth0.Api{
		Url:   apiServer.URL,
		Token: "valid_token",
	}

	api.Post(endpointUrl, body)
}

func mockApiHandler(t *testing.T, endpointUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			if r.URL.Path == endpointUrl {
				requestData, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err.Error())
				}
				defer r.Body.Close()

				if r.Header["Authorization"][0] != "Bearer valid_token" {
					t.Error("Invalid authorization token header received")
				}

				if r.Header["Content-Type"][0] != "application/json; charset=utf-8" {
					t.Error("Invalid content-type header received")
				}

				if string(requestData) != "{\"test_field\":\"some data\"}" {
					t.Error("Invalid request body received")
				}

				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}
}

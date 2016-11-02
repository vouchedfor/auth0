package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Api struct {
	Url   string
	Token string
}

func (api *Api) Post(endpointUrl string, body interface{}) (*http.Response, error) {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, api.Url+endpointUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+api.Token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	return client.Do(req)
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("error code %d: %s", er.StatusCode, er.Message)
}

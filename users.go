package auth0

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CreateUser struct {
	Connection    string                 `json:"connection"`
	Email         string                 `json:"email"`
	Password      string                 `json:"password"`
	EmailVerified bool                   `json:"email_verified"`
	UserMetadata  map[string]interface{} `json:"user_metadata"`
	AppMetadata   map[string]interface{} `json:"app_metadata"`
}

type GetUser struct {
	Email         string                 `json:"email"`
	EmailVerified bool                   `json:"email_verified"`
	Username      string                 `json:"username"`
	PhoneNumber   string                 `json:"phone_number"`
	PhoneVerified bool                   `json:"phone_verified"`
	UserID        string                 `json:"user_id"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
	Identities    []Identity             `json:"identities"`
	AppMetadata   map[string]interface{} `json:"app_metadata"`
	UserMetadata  map[string]interface{} `json:"user_metadata"`
	Picture       string                 `json:"picture"`
	Name          string                 `json:"name"`
	Nickname      string                 `json:"nickname"`
	Multifactor   []string               `json:"multifactor"`
	LastIP        string                 `json:"last_ip"`
	LastLogin     string                 `json:"last_login"`
	LoginsCount   int                    `json:"logins_count"`
	Blocked       bool                   `json:"blocked"`
	GivenName     string                 `json:"given_name"`
	FamilyName    string                 `json:"family_name"`
}

type Identity struct {
	Connection string `json:"connection"`
	UserID     string `json:"user_id"`
	Provider   string `json:"provider"`
	IsSocial   bool   `json:"isSocial"`
}

func (api *Api) CreateUser(user CreateUser) *ErrorResponse {
	result, err := api.Post("/api/v2/users", user)
	if err != nil {
		return &ErrorResponse{Message: err.Error()}
	}

	defer result.Body.Close()
	responseData, err := ioutil.ReadAll(result.Body)
	if err != nil {
		panic(err.Error())
	}

	if result.StatusCode != http.StatusCreated {
		errorResponse := ErrorResponse{}
		err = json.Unmarshal(responseData, &errorResponse)
		if err != nil {
			panic(err.Error())
		}

		return &errorResponse
	}

	return nil
}

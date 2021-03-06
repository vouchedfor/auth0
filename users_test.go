package auth0_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vouchedfor/auth0"
)

type apiErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
	ErrorCode  string `json:"errorCode"`
}

func TestApi_CreateUser(t *testing.T) {
	appMetadata := map[string]interface{}{
		"userId":   342,
		"userType": "client",
	}
	user := auth0.CreateUserRequestData{
		Connection:   "test",
		Email:        "test_email@gmail.com",
		Password:     "test_password",
		AppMetadata:  appMetadata,
		UserMetadata: make(map[string]interface{}),
	}

	apiServer := httptest.NewServer(http.HandlerFunc(mockUserHandler(user, auth0.UpdateUserRequestData{})))
	defer apiServer.Close()

	api := auth0.Api{
		Url:   apiServer.URL,
		Token: "valid_token",
	}

	if err := api.CreateUser(user); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestApi_CreateUserEmailAlreadyExists(t *testing.T) {
	appMetadata := map[string]interface{}{
		"userId":   342,
		"userType": "client",
	}
	user := auth0.CreateUserRequestData{
		Email:        "mail_exists@test.com",
		Password:     "test_password",
		AppMetadata:  appMetadata,
		UserMetadata: make(map[string]interface{}),
	}

	apiServer := httptest.NewServer(http.HandlerFunc(mockUserHandler(user, auth0.UpdateUserRequestData{})))
	defer apiServer.Close()

	api := auth0.Api{
		Url:               apiServer.URL,
		Token:             "valid_token",
		DefaultConnection: "Username-Password-Authentication",
	}

	expectedError := api.CreateUser(user)
	if expectedError == nil {
		t.Error("should return response error")
	}
	if expectedError.StatusCode != 400 {
		t.Error("response error should have 400 code")
	}
	if expectedError.Message != "The user already exists" {
		t.Error("response error should have user already exists message")
	}
}

func TestApi_UpdateUser(t *testing.T) {
	user := auth0.UpdateUserRequestData{
		Connection:    "test",
		ID:            "test_id",
		Email:         "updated_email@gmail.com",
		EmailVerified: true,
		Password:      "updated_password",
	}

	apiServer := httptest.NewServer(http.HandlerFunc(mockUserHandler(auth0.CreateUserRequestData{}, user)))
	defer apiServer.Close()

	api := auth0.Api{
		Url:   apiServer.URL,
		Token: "valid_token",
	}

	if err := api.UpdateUser(user); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestApi_SendVerificationEmail(t *testing.T) {
	requestData := auth0.SendVerificationEmailRequestData{
		UserId: "auth0|testststststs",
	}

	apiServer := httptest.NewServer(http.HandlerFunc(mockUserHandler(auth0.CreateUserRequestData{}, auth0.UpdateUserRequestData{})))
	defer apiServer.Close()

	api := auth0.Api{
		Url:   apiServer.URL,
		Token: "valid_token",
	}

	if err := api.SendVerificationEmail(requestData); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func mockUserHandler(createUser auth0.CreateUserRequestData, updateUser auth0.UpdateUserRequestData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			switch r.URL.Path {
			case "/api/v2/users":
				if createUser.Email == "mail_exists@test.com" {
					w.WriteHeader(http.StatusBadRequest)
					apiResponse := apiErrorResponse{
						StatusCode: 400,
						Error:      "Bad Request",
						Message:    "The user already exists",
						ErrorCode:  "auth0_idp_error",
					}
					if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
						log.Println("Mock API: Failed to encode output body")
					}

					return
				}

				w.WriteHeader(http.StatusCreated)
				if err := json.NewEncoder(w).Encode(getCreateUserSampleResponse(createUser)); err != nil {
					log.Println("Mock API: Failed to encode output body")
				}
				break
			case "/api/v2/jobs/post_verification_email":
				w.WriteHeader(http.StatusCreated)
				break
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		case http.MethodPatch:
			if r.URL.Path == "/api/v2/users/test_id" {
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(getUpdateUserSampleResponse(updateUser)); err != nil {
					log.Println("Mock API: Failed to encode output body")
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func getCreateUserSampleResponse(createUser auth0.CreateUserRequestData) auth0.GetUser {
	return auth0.GetUser{
		Email:         createUser.Email,
		EmailVerified: createUser.EmailVerified,
		Username:      "",
		PhoneNumber:   "",
		PhoneVerified: false,
		UserID:        "usr_5457edea1b8f33391a000004",
		CreatedAt:     "",
		UpdatedAt:     "",
		Identities: []auth0.Identity{
			{
				Connection: createUser.Connection,
				UserID:     "5457edea1b8f22891a000004",
				Provider:   "auth0",
				IsSocial:   false,
			},
		},
		AppMetadata:  createUser.AppMetadata,
		UserMetadata: make(map[string]interface{}),
		Picture:      "",
		Name:         "",
		Nickname:     "",
		Multifactor:  nil,
		LastIP:       "",
		LastLogin:    "",
		LoginsCount:  0,
		Blocked:      false,
		GivenName:    "",
		FamilyName:   "",
	}
}

func getUpdateUserSampleResponse(updateUser auth0.UpdateUserRequestData) auth0.GetUser {
	return auth0.GetUser{
		Email:         updateUser.Email,
		EmailVerified: false,
		Username:      "",
		PhoneNumber:   "",
		PhoneVerified: false,
		UserID:        "usr_5457edea1b8f33391a000004",
		CreatedAt:     "",
		UpdatedAt:     "",
		Identities: []auth0.Identity{
			{
				Connection: updateUser.Connection,
				UserID:     "5457edea1b8f22891a000004",
				Provider:   "auth0",
				IsSocial:   false,
			},
		},
		AppMetadata:  make(map[string]interface{}),
		UserMetadata: make(map[string]interface{}),
		Picture:      "",
		Name:         "",
		Nickname:     "",
		Multifactor:  nil,
		LastIP:       "",
		LastLogin:    "",
		LoginsCount:  0,
		Blocked:      false,
		GivenName:    "",
		FamilyName:   "",
	}
}

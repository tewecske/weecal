package handlers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"weecal/internal/hash/passwordhash"
	"weecal/internal/store/session"
	"weecal/internal/store/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePostLogin(t *testing.T) {
	testUser := user.User{ID: 1, Email: "user1@test.com", Password: "Ultr4S3cur3"}
	testSession := session.Session{SessionID: "testSessionId1"}
	testCases := []struct {
		testName           string
		userStore          *user.TestUserStore
		sessionStore       *session.TestSessionStore
		hash               *passwordhash.TestPasswordHash
		expectedStatusCode int
		expectedCookie     *http.Cookie
	}{
		{
			testName: "success",
			userStore: &user.TestUserStore{
				User: &testUser,
			},
			sessionStore: &session.TestSessionStore{
				Session: &session.Session{
					SessionID: testSession.SessionID,
				},
			},
			hash:               &passwordhash.TestPasswordHash{Match: true},
			expectedStatusCode: http.StatusOK,
			expectedCookie: &http.Cookie{
				Name:     "session",
				Value:    base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", testSession.SessionID, testUser.ID))),
				HttpOnly: true,
			},
		},
		{
			testName: "fail - user not found",
			userStore: &user.TestUserStore{
				User: &user.User{
					Email:    testUser.Email,
					Password: testUser.Password,
				},
				Error: errors.New("User not found"),
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName: "fail - invalid password",
			userStore: &user.TestUserStore{
				User: &user.User{
					Email:    testUser.Email,
					Password: testUser.Password,
				},
			},
			hash:               &passwordhash.TestPasswordHash{Match: false},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			body := bytes.NewBufferString("email=" + tc.userStore.User.Email + "&password=" + tc.userStore.User.Password)
			req := httptest.NewRequest(http.MethodPost, "/login", body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()

			HandlePostLogin(tc.userStore, tc.sessionStore, tc.hash, "session")(rr, req)

			if rr.Result().StatusCode != tc.expectedStatusCode {
				t.Errorf("Expected status %v but got %v", tc.expectedStatusCode, rr.Result().StatusCode)
			}

			cookies := rr.Result().Cookies()

			if tc.expectedCookie != nil {
				sessionCookie := cookies[0]
				if sessionCookie.Name != tc.expectedCookie.Name {
					t.Errorf("Expected cookie named %v but got %v", tc.expectedCookie.Name, sessionCookie.Name)
				}
				if sessionCookie.Value != tc.expectedCookie.Value {
					t.Errorf("Expected session cookie value %v but got %v", tc.expectedCookie.Value, sessionCookie.Value)
				}
				if sessionCookie.HttpOnly != tc.expectedCookie.HttpOnly {
					t.Errorf("Expected session cookie HttpOnly to be %v but got %v", tc.expectedCookie.HttpOnly, sessionCookie.HttpOnly)
				}
			} else if len(cookies) != 0 {
				t.Errorf("Expected no cookies but got %v", len(cookies))
			}
		})
	}

}

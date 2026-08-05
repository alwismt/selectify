package auth_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"alwis.dev/selectify/internal/httpx/request"
	"alwis.dev/selectify/internal/testkit/test"
)

func TestController_UserRegister(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		email := fmt.Sprintf("testuser%d@example.com", rand.Int())
		req := request.UserRegisterRequest{
			Email:           email,
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Test",
			LastName:        "User",
			Phone:           fmt.Sprintf("+370123456%d", rand.Int()%10000),
		}

		resp := test.DoPost(t, ts.S.URL, "/api/v1/auth/register", req)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		cookies := resp.Cookies()
		require.NotEmpty(t, cookies)
		sessionCookie := findCookie(cookies, "slf")
		require.NotNil(t, sessionCookie)
		require.NotEmpty(t, sessionCookie.Value)

		t.Cleanup(func() {
			var userID uint
			err := ts.DB.RwDb.Get(&userID, "SELECT id FROM users WHERE email = $1", email)
			if err == nil {
				_, _ = ts.DB.RwDb.Exec("DELETE FROM user_session WHERE user_id = $1", userID)
				_, _ = ts.DB.RwDb.Exec("DELETE FROM user_role WHERE user_id = $1", userID)
				_, _ = ts.DB.RwDb.Exec("DELETE FROM users WHERE id = $1", userID)
			}
		})
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		email := fmt.Sprintf("duplicate%d@example.com", rand.Int())
		req := request.UserRegisterRequest{
			Email:           email,
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Test",
			LastName:        "User",
			Phone:           fmt.Sprintf("+370123456%d", rand.Int()%10000),
		}

		resp := test.DoPost(t, ts.S.URL, "/api/v1/auth/register", req)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		t.Cleanup(func() {
			var userID uint
			err := ts.DB.RwDb.Get(&userID, "SELECT id FROM users WHERE email = $1", email)
			if err == nil {
				_, _ = ts.DB.RwDb.Exec("DELETE FROM user_session WHERE user_id = $1", userID)
				_, _ = ts.DB.RwDb.Exec("DELETE FROM user_role WHERE user_id = $1", userID)
				_, _ = ts.DB.RwDb.Exec("DELETE FROM users WHERE id = $1", userID)
			}
		})

		resp2 := test.DoPost(t, ts.S.URL, "/api/v1/auth/register", req)
		require.Equal(t, http.StatusBadRequest, resp2.StatusCode)
	})

	t.Run("PasswordMismatch", func(t *testing.T) {
		req := request.UserRegisterRequest{
			Email:           fmt.Sprintf("test%d@example.com", rand.Int()),
			Password:        "password123",
			ConfirmPassword: "password456",
			FirstName:       "Test",
			LastName:        "User",
			Phone:           fmt.Sprintf("+370123456%d", rand.Int()%10000),
		}

		resp := test.DoPost(t, ts.S.URL, "/api/v1/auth/register", req)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		req := request.UserRegisterRequest{
			Email:           "invalid-email",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Test",
			LastName:        "User",
			Phone:           fmt.Sprintf("+370123456%d", rand.Int()%10000),
		}

		resp := test.DoPost(t, ts.S.URL, "/api/v1/auth/register", req)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingFields", func(t *testing.T) {
		req := request.UserRegisterRequest{
			Email: fmt.Sprintf("test%d@example.com", rand.Int()),
		}

		resp := test.DoPost(t, ts.S.URL, "/api/v1/auth/register", req)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestController_UserLogin(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testUserEmail := "travis@alwis.dev"
		password := "passVVord"

		var userID uint
		err := ts.DB.RwDb.Get(&userID, "SELECT id FROM users WHERE email = $1", testUserEmail)
		require.NoError(t, err, "Test user should exist in database")
		require.NotZero(t, userID)

		loginReq := request.LoginRequest{
			Email:    testUserEmail,
			Password: password,
		}

		resp := test.DoPost(t, ts.S.URL, "/api/v1/auth/login", loginReq)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		cookies := resp.Cookies()
		require.NotEmpty(t, cookies)
		sessionCookie := findCookie(cookies, "slf")
		require.NotNil(t, sessionCookie)
		require.NotEmpty(t, sessionCookie.Value)

		t.Cleanup(func() {
			_, _ = ts.DB.RwDb.Exec("DELETE FROM user_session WHERE user_id = $1", userID)
		})
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		req := request.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "wrongpassword",
		}

		resp := test.DoPost(t, ts.S.URL, "/api/v1/auth/login", req)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		req := request.LoginRequest{
			Email:    "invalid-email",
			Password: "password123",
		}

		resp := test.DoPost(t, ts.S.URL, "/api/v1/auth/login", req)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingFields", func(t *testing.T) {
		req := request.LoginRequest{
			Email: fmt.Sprintf("test%d@example.com", rand.Int()),
		}

		resp := test.DoPost(t, ts.S.URL, "/api/v1/auth/login", req)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

package auth

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetLoggedInUser(t *testing.T) {
	t.Run("user in context", func(t *testing.T) {
		claims := jwt.MapClaims{
			"user_id":    "user-1",
			"email":      "user1@example.com",
			"type":       "admin",
			"session_id": "session-123",
		}
		token := &jwt.Token{
			Claims: claims,
		}

		ctx := context.WithValue(context.Background(), ContextKeyUser, token)
		user := GetLoggedInUser(ctx)

		assert.Equal(t, "user-1", user.ID)
		assert.Equal(t, "user1@example.com", user.Email)
		assert.Equal(t, "admin", user.Type)
		assert.Equal(t, "session-123", user.SessionID)
	})

	t.Run("no user in context", func(t *testing.T) {
		user := GetLoggedInUser(context.Background())
		assert.Equal(t, UserAuth{}, user)
	})

	t.Run("user in context with missing fields", func(t *testing.T) {
		claims := jwt.MapClaims{
			"user_id": "user-2",
			// other fields missing
		}
		token := &jwt.Token{
			Claims: claims,
		}

		ctx := context.WithValue(context.Background(), ContextKeyUser, token)
		user := GetLoggedInUser(ctx)

		assert.Equal(t, "user-2", user.ID)
		assert.Empty(t, user.Email)
		assert.Empty(t, user.Type)
		assert.Empty(t, user.SessionID)
	})
}

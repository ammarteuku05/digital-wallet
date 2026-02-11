package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type UserAuth struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	SessionID string `json:"session_id"`
	Type      string `json:"type"`
}

type ContextUser string

const ContextKeyUser ContextUser = "user"

func GetLoggedInUser(ctx context.Context) UserAuth {

	v := ctx.Value(ContextKeyUser)
	if v == nil {
		return UserAuth{}
	}

	loggedInUser := v.(*jwt.Token)

	claims := loggedInUser.Claims.(jwt.MapClaims)

	var id string
	if val, ok := claims["user_id"].(string); ok {
		id = val
	}

	var email string
	if val, ok := claims["email"].(string); ok {
		email = val
	}

	var userType string
	if val, ok := claims["type"].(string); ok {
		userType = val
	}

	var sessionID string
	if val, ok := claims["session_id"].(string); ok {
		sessionID = val
	}

	return UserAuth{
		ID:        id,
		Email:     email,
		Type:      userType,
		SessionID: sessionID,
	}

}

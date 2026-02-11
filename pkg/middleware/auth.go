package middleware

// func AuthMiddleware(di *di.Container) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			token, err := auth.VerifyTokenFromRequest(c, di.Config.JWT.SigningKey)
// 			if err != nil {
// 				return response.NewUnauthorizedError(response.ErrUnauthorizedType.Message)
// 			}

// 			claims := token.Claims.(jwt.MapClaims)

// 			var tokenType string
// 			if val, ok := claims["type"].(string); ok {
// 				tokenType = val
// 			}

// 			if tokenType != "access" {
// 				return response.NewUnauthorizedError(response.ErrUnauthorizedType.Message)
// 			}

// 			validToken, err := di.RepoRegistry.GetSessionRepository().IsValidSession(context.Background(), claims["user_id"].(string), claims["session_id"].(string))
// 			if err != nil {
// 				return response.NewUnauthorizedError("Invalid token")
// 			}

// 			if !validToken {
// 				return response.NewUnauthorizedError("the client's session has expired")
// 			}
// 			ctx := context.WithValue(c.Request().Context(), auth.ContextKeyUser, token)
// 			r := c.Request().WithContext(ctx)
// 			c.SetRequest(r)

// 			if ctx.Value(auth.ContextKeyUser) != nil {
// 				return next(c)
// 			}

// 			return response.NewUnauthorizedError(response.ErrUnauthorizedType.Message)
// 		}
// 	}
// }

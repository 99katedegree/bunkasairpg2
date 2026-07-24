package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const UserIDKey = "userID"

// echoContextKeyType はコンテキストキーの型（衝突防止）
type echoContextKeyType struct{}

// EchoContextKey は echo.Context を context.Context に格納するためのキー
var EchoContextKey = echoContextKeyType{}

// InjectEchoContext は echo.Context を context.Context に埋め込む StrictMiddlewareFunc を返す
// oapi-codegen の StrictServerInterface を使う場合、context.Context から echo.Context を取得するために使用する
func InjectEchoContext() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), EchoContextKey, c)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// GetEchoContext は context.Context から echo.Context を取り出す
func GetEchoContext(ctx context.Context) (echo.Context, bool) {
	c, ok := ctx.Value(EchoContextKey).(echo.Context)
	return c, ok
}

func Auth(jwtSecret string) echo.MiddlewareFunc {
	secret := []byte(jwtSecret)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string][]string{"errors": {"UNAUTHORIZED"}})
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return secret, nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string][]string{"errors": {"UNAUTHORIZED"}})
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string][]string{"errors": {"UNAUTHORIZED"}})
			}
			sub, ok := claims["sub"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string][]string{"errors": {"UNAUTHORIZED"}})
			}
			id, err := uuid.Parse(sub)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string][]string{"errors": {"UNAUTHORIZED"}})
			}
			c.Set(UserIDKey, id)
			return next(c)
		}
	}
}

// GetUserID は echo.Context から userID を取り出す
func GetUserID(c echo.Context) (uuid.UUID, bool) {
	id, ok := c.Get(UserIDKey).(uuid.UUID)
	return id, ok
}

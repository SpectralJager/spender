package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	secret = []byte("secret")
)

type ctxKey string

const (
	userIDKey ctxKey = "userid"
)

type AuthClaims struct {
	UserID string `json:"userid"`
	jwt.RegisteredClaims
}

func JWTAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tokenStr := ctx.Request().Header.Get("X-Api-Token")
		log.Println("JWT token: ", tokenStr)
		if len(tokenStr) == 0 {
			return fmt.Errorf("unauthorized")
		}

		token, err := parseJWTToken(tokenStr)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("unauthorized")
		}

		claims, ok := token.Claims.(*AuthClaims)
		if !ok {
			log.Panicln("unexpected claims")
			return fmt.Errorf("unauthorized")
		}

		req := ctx.Request()
		ctx.SetRequest(
			req.WithContext(
				context.WithValue(
					req.Context(),
					userIDKey,
					claims.UserID,
				),
			),
		)

		return next(ctx)
	}
}

func parseJWTToken(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
}

func NewJWTTokenString(userID string) (string, error) {
	claims := &AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func GetUserIDFromRequest(req *http.Request) string {
	if userID, ok := req.Context().Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"shopping-cart-backend/config"
	"shopping-cart-backend/internal/domain"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

var identityKey = "id"

type auth struct {
	acntRepository domain.AccountRepository
	role           domain.Role
}

func NewAuthMiddleware(acntRepository domain.AccountRepository, role domain.Role) auth {
	return auth{acntRepository: acntRepository, role: role}
}

func (ath *auth) GetInstance() *jwt.HertzJWTMiddleware {

	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       config.App.Server.JWTRealm,
		Key:         []byte(config.App.Server.JWTSecret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*domain.Account); ok {
				return jwt.MapClaims{
					identityKey:  v.Email,
					"role":       string(v.Role),
					"account_id": v.ID,
					"is_active":  v.Active,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &domain.Account{
				Email:  claims[identityKey].(string),
				Role:   domain.Role(claims["role"].(string)),
				ID:     claims["account_id"].(string),
				Active: claims["is_active"].(bool),
			}
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			type login struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			// TODO: add payload validation here
			var loginVals login
			if err := c.BindAndValidate(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginVals.Email
			password := loginVals.Password

			acnt, err := ath.acntRepository.Get(email)
			if err != nil {
				return nil, errors.New("account does not exist")
			}

			if !acnt.Active {
				return nil, errors.New("account is not active anymore")
			}

			if email == strings.TrimSpace(acnt.Email) && password == strings.TrimSpace(acnt.Password) && acnt.Active {
				return &domain.Account{
					Email:  email,
					Active: acnt.Active,
					ID:     acnt.ID,
					Role:   acnt.Role,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			claims := jwt.ExtractClaims(ctx, c)
			email := claims[identityKey].(string)
			role := claims["role"].(string)
			active := claims["is_active"].(bool)

			// TODO: Remove after testing
			fmt.Println("email: ", email)
			fmt.Println("role: ", role)
			fmt.Println("ath role: ", ath.role)
			fmt.Println("is active? ", active)

			if v, ok := data.(*domain.Account); ok && v.Email == email && v.Role == domain.Role(role) &&
				domain.Role(strings.TrimSpace(string(v.Role))) == ath.role && v.Active && active {
				return true
			}
			return false
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer". If you want empty value, use WithoutDefaultTokenHeadName.
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}

package middleware

import (
	"context"
	"errors"
	"log"
	"shopping-cart-backend/internal/domain"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

var identityKey = "id"

type auth struct {
	acntRepository domain.AccountRepository
}

func NewAuthMiddleware(acntRepository domain.AccountRepository) auth {
	return auth{acntRepository: acntRepository}
}

func (ath *auth) GetInstance() *jwt.HertzJWTMiddleware {

	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "test zone",          // read from config
		Key:         []byte("secret key"), // read from config
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey, // read from config
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*domain.Account); ok {
				return jwt.MapClaims{
					identityKey: v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &domain.Account{
				Name: claims[identityKey].(string),
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
			// TODO: Remove commented code after testing

			// reqData := c.Request.Body()
			// acc := domain.Account{}
			// if err := json.Unmarshal(reqData, &acc); err != nil {
			// 	fmt.Printf("error encountered while unmarshalling JSON: %v\n", err)
			// 	c.JSON(400, map[string]string{
			// 		"error": "error encountered while unmarshalling JSON",
			// 	})
			// 	return false
			// }
			// fmt.Println("data: ", acc)

			// TODO: Work on authorization in the correct manner
			if _, ok := data.(*domain.Account); ok {
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

package api

import (
	"bytes"
	"context"
	"net/http"
	"shopping-cart-backend/config"
	"shopping-cart-backend/internal/repository"
	"shopping-cart-backend/internal/service"
	"shopping-cart-backend/pkg/database"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestRegisterRequest_Valid(t *testing.T) {
	config.Load()

	h := server.Default()
	databaseCfg := database.Config{
		Name:     config.App.Database.Name,
		User:     config.App.Database.User,
		Password: config.App.Database.Password,
		Host:     config.App.Database.Host,
		Port:     config.App.Database.Port,
	}

	db, err := database.NewFromEnv(context.Background(), &databaseCfg)
	if err != nil {
		hlog.Fatal(err)
	}
	acntRepo := repository.NewAccountRepository(db)
	acntService := service.NewAccountService(acntRepo)
	regnHandler := NewRegistrationHandler(acntService)

	h.POST("/register", regnHandler.HandleRegistration)
	data := bytes.NewBufferString("{\"name\":\"mayank\", \"role\":\"user\", \"email\":\"hi@hi.com\", \"password\":\"hellothere\"}")
	w := ut.PerformRequest(h.Engine, http.MethodPost, "/register",
		&ut.Body{Body: data, Len: data.Len()},
		ut.Header{Key: "Content-Type", Value: "application/json"})
	resp := w.Result()
	assert.DeepEqual(t, 201, resp.StatusCode())
	assert.DeepEqual(t, "{\"status\":\"success\",\"message\":\"account created successfully.\"}", string(resp.Body()))
	db.Execute("DELETE FROM accounts WHERE email='hi@hi.com'")
}

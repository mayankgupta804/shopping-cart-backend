package api

import (
	"context"
	"encoding/json"
	"fmt"
	"shopping-cart-backend/internal/serializer"
	"shopping-cart-backend/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"gopkg.in/dealancer/validate.v2"
)

func NewRegistrationHandler(acntService service.AccountService) *RegistrationHandler {
	return &RegistrationHandler{acntService: acntService}
}

type RegistrationHandler struct {
	acntService service.AccountService
}

func (regnHandler RegistrationHandler) HandleRegistration(c context.Context, ctx *app.RequestContext) {
	var err error
	reqData := ctx.Request.Body()
	createAccountReq := serializer.CreateAccountRequest{}
	if err = json.Unmarshal(reqData, &createAccountReq); err != nil {
		fmt.Printf("error encountered while unmarshalling JSON: %v\n", err)
		ctx.JSON(400, serializer.Error{Error: "JSON body seems to be malformed"})
		return
	}

	if err = validate.Validate(createAccountReq); err != nil {
		// maybe log validation errors to get a hang of how many times
		// users are unable to register. this would give us the idea
		// about how we can make the API friendly to use
		ctx.JSON(400, serializer.Error{Error: err.Error()})
		return
	}

	err = regnHandler.acntService.Create(createAccountReq)

	if err == service.ErrAccountAlreadyExists {
		ctx.JSON(400, serializer.Error{Error: err.Error()})
		return
	} else if err != nil {
		// Report issue to sentry and raise an alert
		fmt.Printf("internal server error: %v\n", err)
		ctx.JSON(500, serializer.Error{Error: "internal server error"})
		return
	}

	// log all successful account creations for analytics purpose
	// this would help us with users that are signing up on a daily/weekly/monthly basis
	ctx.JSON(201, serializer.CreateAccountResponse{Status: "success", Message: "account created successfully."})
}

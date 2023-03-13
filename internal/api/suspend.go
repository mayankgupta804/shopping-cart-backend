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

func NewSuspendHandler(acntService service.AccountService) *SuspendHandler {
	return &SuspendHandler{acntService: acntService}
}

type SuspendHandler struct {
	acntService service.AccountService
}

func (susHandler SuspendHandler) HandleAccountSuspension(c context.Context, ctx *app.RequestContext) {
	var err error
	reqData := ctx.Request.Body()
	suspendAccountReq := serializer.SuspendAccountRequest{}
	if err = json.Unmarshal(reqData, &suspendAccountReq); err != nil {
		fmt.Printf("error encountered while unmarshalling JSON: %v\n", err)
		ctx.JSON(400, serializer.Error{Error: "JSON body seems to be malformed"})
		return
	}

	//TODO: Handle case where the account is the current account.

	if err = validate.Validate(suspendAccountReq); err != nil {
		// maybe log validation errors to get a hang of how many times
		// users are unable to register. this would give us the idea
		// about how we can make the API friendly to use
		ctx.JSON(400, serializer.Error{Error: err.Error()})
		return
	}

	err = susHandler.acntService.Suspend(suspendAccountReq)

	if err == service.ErrAccountDoesNotExist {
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
	ctx.JSON(200, serializer.SuspendAccountResponse{Status: "success", Message: "account suspended successfully."})
}

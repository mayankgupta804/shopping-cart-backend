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

func NewItemHandler(itemServ service.ItemService) *ItemHandler {
	return &ItemHandler{itemServ: itemServ}
}

type ItemHandler struct {
	itemServ service.ItemService
}

func (itHandler ItemHandler) HandleAddItem(c context.Context, ctx *app.RequestContext) {
	var err error
	reqData := ctx.Request.Body()
	createItemReq := serializer.CreateItemRequest{}
	if err = json.Unmarshal(reqData, &createItemReq); err != nil {
		fmt.Printf("error encountered while unmarshalling JSON: %v\n", err)
		ctx.JSON(400, map[string]string{
			"error": "error encountered while unmarshalling JSON",
		})
		return
	}

	if err = validate.Validate(createItemReq); err != nil {
		// maybe log validation errors to get a hang of how many times
		// users are unable to register. this would give us the idea
		// about how we can make the API friendly to use
		ctx.JSON(400, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err = itHandler.itemServ.Add(createItemReq); err != nil {
		// Report issue to sentry and raise an alert
		fmt.Printf("internal server error: %v\n", err)
		ctx.JSON(500, map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, serializer.CreateItemResponse{Status: "success", Message: "item added successfully"})
}

func (itHandler ItemHandler) HandleGetItem(c context.Context, ctx *app.RequestContext) {
	items, err := itHandler.itemServ.List()
	if err != nil {
		// Report issue to sentry and raise an alert
		fmt.Printf("internal server error: %v\n", err)
		ctx.JSON(500, map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, items)
}

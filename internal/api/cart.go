package api

import (
	"context"
	"encoding/json"
	"fmt"
	"shopping-cart-backend/internal/serializer"
	"shopping-cart-backend/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"gopkg.in/dealancer/validate.v2"
)

func NewCartHandler(cartServ service.CartService) *CartHandler {
	return &CartHandler{cartServ: cartServ}
}

type CartHandler struct {
	cartServ service.CartService
}

func (crtHandler CartHandler) HandleAddToCart(c context.Context, ctx *app.RequestContext) {
	var err error
	reqData := ctx.Request.Body()
	addToCartReq := serializer.AddToCartRequest{}
	if err = json.Unmarshal(reqData, &addToCartReq); err != nil {
		fmt.Printf("error encountered while unmarshalling JSON: %v\n", err)
		ctx.JSON(400, serializer.Error{Error: "error encountered while unmarshalling JSON"})
		return
	}

	if err = validate.Validate(addToCartReq); err != nil {
		ctx.JSON(400, serializer.Error{Error: err.Error()})
		return
	}

	claims := jwt.ExtractClaims(c, ctx)
	accountID := claims["account_id"].(string)

	if err = crtHandler.cartServ.Add(accountID, addToCartReq); err != nil {
		// Report issue to sentry and raise an alert
		fmt.Printf("internal server error: %v\n", err)
		ctx.JSON(500, serializer.Error{Error: err.Error()})
		return
	}

	ctx.JSON(201, serializer.AddToCartResponse{Status: "success", Message: "item added successfully"})
}

func (crtHandler CartHandler) HandleRemoveFromCart(c context.Context, ctx *app.RequestContext) {
	var err error
	itemId, ok := ctx.GetQuery("item_id")
	if !ok {
		ctx.JSON(400, serializer.Error{Error: "item_id must not be empty"})
		return
	}
	removeFromCartReq := serializer.RemoveFromCartRequest{ItemID: itemId}
	if err = validate.Validate(removeFromCartReq); err != nil {
		ctx.JSON(400, serializer.Error{Error: err.Error()})
		return
	}

	claims := jwt.ExtractClaims(c, ctx)
	accountID := claims["account_id"].(string)

	if err = crtHandler.cartServ.Remove(accountID, removeFromCartReq); err != nil {
		// Report issue to sentry and raise an alert
		fmt.Printf("internal server error: %v\n", err)
		ctx.JSON(500, serializer.Error{Error: err.Error()})
		return
	}

	ctx.JSON(204, serializer.RemoveFromCartResponse{Status: "success", Message: "item removed successfully"})
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"shopping-cart-backend/internal/serializer"
	"shopping-cart-backend/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
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
		ctx.JSON(400, map[string]string{
			"error": "error encountered while unmarshalling JSON",
		})
		return
	}

	// if err = validate.Validate(addToCartReq); err != nil {
	// 	ctx.JSON(400, map[string]string{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	claims := jwt.ExtractClaims(c, ctx)
	accountID := claims["account_id"].(string)

	// TODO: Remove after testing
	fmt.Println("accountID: ", accountID)

	if err = crtHandler.cartServ.Add(accountID, addToCartReq); err != nil {
		// Report issue to sentry and raise an alert
		fmt.Printf("internal server error: %v\n", err)
		ctx.JSON(500, map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, serializer.CreateItemResponse{Status: "success", Message: "item added successfully"})
}

func (crtHandler CartHandler) HandleRemoveFromCart(c context.Context, ctx *app.RequestContext) {
	itemId, ok := ctx.GetQuery("item_id")
	if !ok {
		ctx.JSON(400, map[string]string{
			"error": "item_id must not be empty.",
		})
		return
	}
	removeFromCartReq := serializer.RemoveFromCartRequest{ItemID: itemId}

	claims := jwt.ExtractClaims(c, ctx)
	accountID := claims["account_id"].(string)

	// TODO: Remove after testing
	fmt.Println("accountID: ", accountID)

	if err := crtHandler.cartServ.Remove(accountID, removeFromCartReq); err != nil {
		// Report issue to sentry and raise an alert
		fmt.Printf("internal server error: %v\n", err)
		ctx.JSON(500, map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(204, serializer.CreateItemResponse{Status: "success", Message: "item removed successfully"})
}

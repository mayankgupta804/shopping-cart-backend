package api

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestPingRequest(t *testing.T) {
	h := server.Default()
	h.GET("/ping", PingHandler)
	w := ut.PerformRequest(h.Engine, "GET", "/ping", nil)
	resp := w.Result()
	assert.DeepEqual(t, 200, resp.StatusCode())
	assert.DeepEqual(t, "{\"ping\":\"pong\"}", string(resp.Body()))
}

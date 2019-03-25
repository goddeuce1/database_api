package middlewares

import (
	"github.com/valyala/fasthttp"
)

//SetHeaders - sets headers for response
func SetHeaders(ctx *fasthttp.RequestCtx, code int) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
}

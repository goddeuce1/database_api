package api

import (
	"encoding/json"

	mw "../middlewares"
	"../models"
	"github.com/valyala/fasthttp"
)

//PostIDDetailsPost - updates post info
func PostIDDetailsPost(ctx *fasthttp.RequestCtx) {
	postUpdate := models.PostUpdate{}
	err := json.Unmarshal(ctx.PostBody(), &postUpdate)

	if err != nil {
		mw.SetHeaders(ctx, fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
		return
	}

	postID := ctx.UserValue("id").(string)
	response, error := mw.PostIDDetailsPostMiddleware(postUpdate.Message, postID)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := json.Marshal(response)
		ctx.Write(result)

	} else if error == models.ErrPostNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := json.Marshal(error)
		ctx.Write(result)
	}

}

//PostIDDetailsGet - get post by its id
func PostIDDetailsGet(ctx *fasthttp.RequestCtx) {
	postID := ctx.UserValue("id").(string)
	related := string(ctx.FormValue("related"))
	response, error := mw.PostIDDetailsGetMiddleware(postID, related)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := json.Marshal(response)
		ctx.Write(result)

	} else if error == models.ErrPostNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := json.Marshal(error)
		ctx.Write(result)
	}
}

package api

import (
	"encoding/json"
	mw "park_base/park_db/middlewares"
	"park_base/park_db/models"

	"github.com/valyala/fasthttp"
)

//PostIDDetailsPost - updates post info
func PostIDDetailsPost(ctx *fasthttp.RequestCtx) {
	postUpdate := models.PostUpdate{}
	postUpdate.UnmarshalJSON(ctx.PostBody())

	postID := ctx.UserValue("id").(string)
	response, error := mw.PostIDDetailsPostMiddleware(postUpdate.Message, postID)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := response.MarshalJSON()
		ctx.Write(result)

	} else if error == models.ErrPostNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := error.MarshalJSON()
		ctx.Write(result)
	}

	return
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
		result, _ := error.MarshalJSON()
		ctx.Write(result)
	}

	return
}

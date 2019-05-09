package api

import (
	mw "park_base/park_db/middlewares"
	"park_base/park_db/models"

	"github.com/valyala/fasthttp"
)

//ForumCreate - creates forum
func ForumCreate(ctx *fasthttp.RequestCtx) {
	forum := models.Forum{}
	forum.UnmarshalJSON(ctx.PostBody())

	response, error := mw.ForumCreateMiddleware(&forum)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusCreated)
		result, _ := response.MarshalJSON()
		ctx.Write(result)

	} else if error == models.ErrForumAlreadyExists {
		mw.SetHeaders(ctx, fasthttp.StatusConflict)
		result, _ := response.MarshalJSON()
		ctx.Write(result)

	} else if error == models.ErrForumOwnerNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := error.MarshalJSON()
		ctx.Write(result)
	}

	return

}

//ForumSlugCreate - create thread
func ForumSlugCreate(ctx *fasthttp.RequestCtx) {
	thread := models.Thread{}
	thread.UnmarshalJSON(ctx.PostBody())

	slug := ctx.UserValue("slug").(string)
	response, error := mw.ForumSlugCreateMiddleware(&thread, slug)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusCreated)
		if len(response) == 1 {
			result, _ := response[0].MarshalJSON()
			ctx.Write(result)
		} else {
			result, _ := response.MarshalJSON()
			ctx.Write(result)
		}

	} else if error == models.ErrThreadAlreadyExists {
		mw.SetHeaders(ctx, fasthttp.StatusConflict)
		if len(response) == 1 {
			result, _ := response[0].MarshalJSON()
			ctx.Write(result)
		} else {
			result, _ := response.MarshalJSON()
			ctx.Write(result)
		}

	} else if error == models.ErrForumOrAuthorNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := error.MarshalJSON()
		ctx.Write(result)
	}

	return
}

//ForumSlugDetails - returns forum details by slug
func ForumSlugDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	response, err := mw.ForumSlugDetailsMiddleware(slug)

	if err == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := response.MarshalJSON()
		ctx.Write(result)

	} else if err == models.ErrForumNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := err.MarshalJSON()
		ctx.Write(result)
	}

	return
}

//ForumSlugThreads - returns threads from forum by slug
func ForumSlugThreads(ctx *fasthttp.RequestCtx) {
	limit := string(ctx.FormValue("limit"))
	since := string(ctx.FormValue("since"))
	desc := string(ctx.FormValue("desc"))

	slug := ctx.UserValue("slug").(string)

	response, err := mw.ForumSlugThreadsMiddleware(limit, since, desc, slug)

	if err == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := response.MarshalJSON()
		ctx.Write(result)

	} else if err == models.ErrForumNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := err.MarshalJSON()
		ctx.Write(result)
	}

	return
}

//ForumSlugUsers - returns users limit/desc/since
func ForumSlugUsers(ctx *fasthttp.RequestCtx) {
	limit := string(ctx.FormValue("limit"))
	since := string(ctx.FormValue("since"))
	desc := string(ctx.FormValue("desc"))

	slug := ctx.UserValue("slug").(string)

	response, err := mw.ForumSlugUsersMiddleware(limit, since, desc, slug)

	if err == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := response.MarshalJSON()
		ctx.Write(result)

	} else if err == models.ErrForumNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := err.MarshalJSON()
		ctx.Write(result)
	}

	return
}

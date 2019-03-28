package api

import (
	"encoding/json"
	"strconv"

	mw "../middlewares"
	"../models"
	"github.com/valyala/fasthttp"
)

//ThreadCreate - creates thread
func ThreadCreate(ctx *fasthttp.RequestCtx) {
	posts := models.Posts{}
	err := json.Unmarshal(ctx.PostBody(), &posts)

	if err != nil {
		mw.SetHeaders(ctx, fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
		return
	}

	thread := ctx.UserValue("slug_or_id").(string)

	if threadID, error := strconv.Atoi(thread); error == nil {

		for _, index := range posts {
			index.Thread = threadID
		}

	}

	response, error := mw.ThreadCreateMiddleware(posts, thread)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusCreated)
		result, _ := json.Marshal(response)
		ctx.Write(result)

	} else if error == models.ErrThreadAlreadyExists || error == models.ErrParentNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusConflict)
		result, _ := json.Marshal(response)
		ctx.Write(result)

	} else if error == models.ErrThreadAlreadyExists ||
		error == models.ErrUserNotFound || error == models.ErrThreadNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := json.Marshal(error)
		ctx.Write(result)
	}

}

//ThreadVote - sets +-1 rating to thread
func ThreadVote(ctx *fasthttp.RequestCtx) {
	vote := models.Vote{}
	err := json.Unmarshal(ctx.PostBody(), &vote)

	if err != nil {
		mw.SetHeaders(ctx, fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
		return
	}

	threadSlug := ctx.UserValue("slug_or_id").(string)

	response, error := mw.ThreadSlugVoteMiddleware(&vote, threadSlug)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := json.Marshal(response)
		ctx.Write(result)

	} else if error == models.ErrThreadNotFound || error == models.ErrUserNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := json.Marshal(error)
		ctx.Write(result)

	}

}

//ThreadDetailsGet - get info about thread by slug/id
func ThreadDetailsGet(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug_or_id").(string)

	response, err := mw.ThreadDetailsGetMiddleware(slug)

	if err == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := json.Marshal(response)
		ctx.Write(result)

	} else if err == models.ErrThreadNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := json.Marshal(err)
		ctx.Write(result)
	}
}

//ThreadDetailsPost - updates thread info
func ThreadDetailsPost(ctx *fasthttp.RequestCtx) {
	threadUpdate := models.ThreadUpdate{}
	err := json.Unmarshal(ctx.PostBody(), &threadUpdate)

	if err != nil {
		mw.SetHeaders(ctx, fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
		return
	}

	threadSlug := ctx.UserValue("slug_or_id").(string)

	response, error := mw.ThreadDetailsPostMiddleware(&threadUpdate, threadSlug)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := json.Marshal(response)
		ctx.Write(result)

	} else if error == models.ErrThreadNotFound || error == models.ErrUserNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := json.Marshal(error)
		ctx.Write(result)

	}
}

//ThreadPosts - returns thread posts
func ThreadPosts(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug_or_id").(string)
	limit := string(ctx.FormValue("limit"))
	since := string(ctx.FormValue("since"))
	sort := string(ctx.FormValue("sort"))
	desc := string(ctx.FormValue("desc"))

	response, error := mw.ThreadPostsMiddleware(slug, limit, since, sort, desc)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := json.Marshal(response)
		ctx.Write(result)

	} else if error == models.ErrThreadNotFound || error == models.ErrUserNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := json.Marshal(error)
		ctx.Write(result)

	}
}

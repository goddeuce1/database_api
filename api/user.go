package api

import (
	"encoding/json"

	mw "../middlewares"
	"../models"
	"github.com/valyala/fasthttp"
)

//UserCreate - creates user
func UserCreate(ctx *fasthttp.RequestCtx) {
	user := models.User{}
	err := json.Unmarshal(ctx.PostBody(), &user)
	user.Nickname = ctx.UserValue("nickname").(string)

	if err != nil {
		mw.SetHeaders(ctx, fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
		return
	}
	response, error := mw.UserCreateMiddleware(&user)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusCreated)
		result, _ := json.Marshal(user)
		ctx.Write(result)

	} else if error == models.ErrUserAlreadyExists {
		mw.SetHeaders(ctx, fasthttp.StatusConflict)
		result, _ := json.Marshal(response)
		ctx.Write(result)
	}
}

//UserProfileGet - returns desired user
func UserProfileGet(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	response, err := mw.UserProfileGetMiddleware(nickname)

	if err == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := json.Marshal(response)
		ctx.Write(result)

	} else if err == models.ErrUserNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := json.Marshal(err)
		ctx.Write(result)
	}
}

//UserProfilePost - returns new user settings
func UserProfilePost(ctx *fasthttp.RequestCtx) {
	user := models.User{}
	err := json.Unmarshal(ctx.PostBody(), &user)

	user.Nickname = ctx.UserValue("nickname").(string)

	if err != nil {
		mw.SetHeaders(ctx, fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
		return
	}

	returnUser, error := mw.UserProfilePostMiddleware(&user)

	if error == nil {
		mw.SetHeaders(ctx, fasthttp.StatusOK)
		result, _ := json.Marshal(returnUser)
		ctx.Write(result)

	} else if error == models.ErrUserNotFound {
		mw.SetHeaders(ctx, fasthttp.StatusNotFound)
		result, _ := json.Marshal(error)
		ctx.Write(result)
	} else if error == models.ErrSettingsConflict {
		mw.SetHeaders(ctx, fasthttp.StatusConflict)
		result, _ := json.Marshal(error)
		ctx.Write(result)
	}
}

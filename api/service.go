package api

import (
	"encoding/json"
	"fmt"

	"../database"
	mw "../middlewares"
	"../models"
	ops "../sqlops"
	"github.com/valyala/fasthttp"
)

//ServiceStatus - returns current status of database
func ServiceStatus(ctx *fasthttp.RequestCtx) {
	status := models.Status{}
	err := database.App.DB.QueryRow(ops.ServiceStatus).Scan(&status.Forum, &status.Post, &status.Thread, &status.User)

	if err != nil {
		mw.SetHeaders(ctx, fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
		return
	}

	mw.SetHeaders(ctx, fasthttp.StatusOK)
	result, _ := json.Marshal(status)
	ctx.Write(result)
}

//ServiceClear - clear everything in database
func ServiceClear(ctx *fasthttp.RequestCtx) {
	_, err := database.App.DB.Exec("TRUNCATE users, forums, threads, posts, votes")
	if err != nil {
		fmt.Println(err)
	}
	mw.SetHeaders(ctx, fasthttp.StatusOK)
}

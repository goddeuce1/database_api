package api

import (
	"park_base/park_db/database"
	"park_base/park_db/models"
	ops "park_base/park_db/sqlops"

	mw "park_base/park_db/middlewares"

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
	result, _ := status.MarshalJSON()
	ctx.Write(result)
}

//ServiceClear - clear everything in database
func ServiceClear(ctx *fasthttp.RequestCtx) {
	database.App.DB.Exec("TRUNCATE users, forums, threads, posts, votes RESTART IDENTITY CASCADE")
	mw.SetHeaders(ctx, fasthttp.StatusOK)
}

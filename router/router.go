package router

import (
	"../api"
	"github.com/buaazp/fasthttprouter"
)

//Router - export router
var Router = fasthttprouter.New()

func init() {
	Router.POST("/api/forum/:slug", api.ForumCreate)
	Router.POST("/api/forum/:slug/create", api.ForumSlugCreate)

	Router.GET("/api/forum/:slug/details", api.ForumSlugDetails)
	Router.GET("/api/forum/:slug/threads", api.ForumSlugThreads)
	Router.GET("/api/forum/:slug/users", api.ForumSlugUsers)
	Router.GET("/api/post/:id/details", api.PostIDDetailsGet)

	Router.POST("/api/post/:id/details", api.PostIDDetailsPost)
	Router.POST("/api/service/clear", api.ServiceClear)

	Router.GET("/api/service/status", api.ServiceStatus)

	Router.POST("/api/thread/:slug_or_id/create", api.ThreadCreate)

	Router.GET("/api/thread/:slug_or_id/details", api.ThreadDetailsGet)

	Router.POST("/api/thread/:slug_or_id/details", api.ThreadDetailsPost)

	Router.GET("/api/thread/:slug_or_id/posts", api.ThreadPosts)

	Router.POST("/api/thread/:slug_or_id/vote", api.ThreadVote)
	Router.POST("/api/user/:nickname/create", api.UserCreate)

	Router.GET("/api/user/:nickname/profile", api.UserProfileGet)

	Router.POST("/api/user/:nickname/profile", api.UserProfilePost)
}

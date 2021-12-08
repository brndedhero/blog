package router

import (
	"github.com/brndedhero/blog/controllers"
	"github.com/brndedhero/blog/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.PrometheusMiddleware)
	r.Path("/metrics").Handler(promhttp.Handler())

	r.HandleFunc("/", controllers.HomeHandler)
	r.HandleFunc("/posts", controllers.AllBlogPostsHandler)
	r.HandleFunc("/posts/new", controllers.NewBlogPostHandler)
	r.HandleFunc("/posts/{id:[\\d]+}", controllers.BlogPostHandler)

	return r
}

package controllers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/brndedhero/blog/helpers"
	"github.com/brndedhero/blog/models"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"message": "Welcome to the Index Page"}`)
}

func AllBlogPostsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "GET":
		data, err := models.GetAllBlogPosts()
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
	default:
		message, err := helpers.PrepareString(405, nil)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, message)
	}
}

func NewBlogPostHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		data, err := models.CreateBlogPost(req.FormValue("title"), req.FormValue("body"))
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, data)
	default:
		message, err := helpers.PrepareString(405, nil)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, message)
	}
}

func BlogPostHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(req)
	id, _ := strconv.ParseUint(params["id"], 10, 64)

	switch req.Method {
	case "GET":
		data, err := models.GetBlogPost(id)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
	case "POST":
		if err := req.ParseForm(); err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		data, err := models.UpdateBlogPost(id, req.FormValue("title"), req.FormValue("body"))
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, data)
	case "DELETE":
		data, err := models.DeleteBlogPost(id)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
	default:
		message, err := helpers.PrepareString(405, nil)
		if err != nil {
			message := helpers.PrepareErrorString(500, err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, message)
	}
}

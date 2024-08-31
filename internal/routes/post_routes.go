package routes

import (
	"github.com/AdblkA/blogging/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterPostRoutes(r *mux.Router, postHandler *handler.PostHandler) {
	r.HandleFunc("/post", postHandler.GetAllPosts).Methods("GET")
	r.HandleFunc("/post", postHandler.CreatePost).Methods("POST")
	r.HandleFunc("/post/{id}", postHandler.GetPostById).Methods("GET")
	r.HandleFunc("/post/{id}", postHandler.DeletePost).Methods("DELETE")
	r.HandleFunc("/post/{id}", postHandler.UpdatePost).Methods("PUT")
}

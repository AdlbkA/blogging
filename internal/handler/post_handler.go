package handler

import (
	"encoding/json"
	"github.com/AdblkA/blogging/internal/models"
	"github.com/AdblkA/blogging/internal/repository"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type PostHandler struct {
	Repo *repository.PostRepository
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := &Response{}
	defer json.NewEncoder(w).Encode(resp)

	repo := repository.PostRepository{Collection: h.Repo.Collection}
	res, err := repo.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()
		log.Println("error: ", err)
	}

	resp.Data = res
	w.WriteHeader(http.StatusOK)
}

func (h *PostHandler) GetPostById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &Response{}
	defer json.NewEncoder(w).Encode(resp)

	repo := repository.PostRepository{Collection: h.Repo.Collection}
	postId := mux.Vars(r)["id"]

	objId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: invalid object id format")
		resp.Error = "invalid object id format"
		return
	}

	res, err := repo.GetByID(objId.Hex())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()
		log.Println("error: ", err)
		return
	}

	resp.Data = res
	w.WriteHeader(http.StatusOK)

}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()

	res, err := h.Repo.Create(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := &Response{}
	defer json.NewEncoder(w).Encode(resp)

	postID := mux.Vars(r)["id"]

	if postID == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = "invalid post id"
		return
	}

	objId, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()
		log.Println("error: invalid object id format")
		return
	}

	var post models.Post

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()
		log.Println("error: invalid post body")
		return
	}

	post.UpdatedAt = time.Now()
	post.ID = objId

	repo := repository.PostRepository{Collection: h.Repo.Collection}
	count, err := repo.Update(post.ID, &post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()
		log.Println("error: ", err)
		return
	}

	resp.Data = count
	w.WriteHeader(http.StatusOK)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &Response{}
	defer json.NewEncoder(w).Encode(resp)

	postID := mux.Vars(r)["id"]

	repo := repository.PostRepository{Collection: h.Repo.Collection}
	count, err := repo.Delete(postID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()
		log.Println("error: ", err)
		return
	}

	resp.Data = count
	w.WriteHeader(http.StatusOK)
}

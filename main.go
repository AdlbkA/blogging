package main

import (
	"github.com/AdblkA/blogging/internal/db"
	"github.com/AdblkA/blogging/internal/handler"
	"github.com/AdblkA/blogging/internal/repository"
	"github.com/AdblkA/blogging/internal/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if err = db.InitDb(os.Getenv("MONGOCONN")); err != nil {
		log.Fatal(err)
	}
	defer db.CloseDb()

	postRepo := &repository.PostRepository{Collection: db.DB.Database(os.Getenv("DBNAME")).Collection(os.Getenv("DBCOLLECTION"))}

	postHandler := &handler.PostHandler{Repo: postRepo}

	r := mux.NewRouter()

	routes.RegisterPostRoutes(r, postHandler)

	log.Fatal(http.ListenAndServe(":8080", r))

}

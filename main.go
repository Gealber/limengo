package main

import (
	"log"

	draftsController "github.com/Gealber/limengo/controllers/drafts"
	draftsRepo "github.com/Gealber/limengo/repositories/drafts"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// I'm going to store data with badger in memory
	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// repository (there's no actual DB, just in memory)
	repo := draftsRepo.New(db)
	// controller
	ctr := draftsController.New(repo)

	r := gin.Default()
	// CORS middleware with an insecure way to set it up
	r.Use(cors.Default())

	api := r.Group("/api")
	{
		api.GET("/front/drafts", ctr.List)
		api.POST("/front/drafts", ctr.Post)
		api.GET("/front/drafts/:context", ctr.Get)
		api.PUT("/front/drafts/:context", ctr.Put)
		api.DELETE("/front/drafts/:context", ctr.Delete)
	}

	r.Run()
}

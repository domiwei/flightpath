package main

import (
	"flightpath/api"
	"flightpath/service/db"
	"flightpath/service/localcache"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.NewMySql("root@tcp(db:3306)/main")
	lc := localcache.NewLocalCache()

	router := gin.Default()
	v1 := router.Group("/v1")
	api.NewHandler(v1.Group("/calculate"), db, lc)
	router.Run(":8080")

	log.Println("system shutdown")
}

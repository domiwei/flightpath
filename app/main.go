package main

import (
	"flightpath/api"
	"flightpath/service/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.NewMySql("root@/main")
	/*
		_, err := db.Exec("insert into `paths` (`person_id`, `path`) values (?,?)", "ggg", "1,2,3,4")
		if err != nil {
			panic(err)
		}
		rows, err := db.Query("select path from `paths` where `person_id`=?", "ggg")
		if err != nil {
			panic(err)
		}
		var r string
		if !rows.Next() {
			panic("fuck")
		}
		if err := rows.Scan(&r); err != nil {
			panic(err)
		}
		fmt.Println(r)
	*/

	router := gin.Default()
	v1 := router.Group("/v1")
	api.NewHandler(v1.Group("/calculate"), db)
	router.Run()
}

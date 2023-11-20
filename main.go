package main

import (
	. "flapypan-blog-go/app"
	. "flapypan-blog-go/conf"
	"flapypan-blog-go/db"
	"flapypan-blog-go/router"
	"log"
)

func main() {
	LoadEnv()
	db.ConnectDB()
	defer db.DB.Close()
	app := CreateApp()
	router.Setup(app)
	log.Fatal(app.Listen(":" + Env("PORT")))
}

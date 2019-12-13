package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nandawinata/entry-task/pkg/handler"
	"github.com/nandawinata/entry-task/pkg/helper/json"
	mw "github.com/nandawinata/entry-task/pkg/middleware"
)

const (
	port = ":8080"
)

func main() {
	router := httprouter.New()
	router.POST("/register", json.ResponseJson(handler.Register))
	router.POST("/login", json.ResponseJson(handler.Login))
	router.GET("/profile", json.ResponseJson(mw.ValidateJwt(handler.Profile)))
	router.GET("/profile/:username", json.ResponseJson(mw.ValidateJwt(handler.ProfileByUsername)))
	router.POST("/update/profile", json.ResponseJson(mw.ValidateJwt(handler.UpdateProfile)))
	router.POST("/update/photo", json.ResponseJson(mw.ValidateJwt(handler.UpdatePhoto)))

	fmt.Printf("Starting at port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

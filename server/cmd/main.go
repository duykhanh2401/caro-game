package main

import (
	"fmt"

	"github.com/duykhanh2401/caro-game/internal/module/ws"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	hub := ws.NewHub()
	hub.Defaults()
	go hub.Run()
	router.Use(static.Serve("/", static.LocalFile("../client/dist", false)))
	router.GET("/ws/caro", hub.Handler())
	if err := router.Run(":2401"); err != nil {
		fmt.Println(err)
	}
}

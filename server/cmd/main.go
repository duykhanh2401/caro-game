package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	if err := router.Run(":2401"); err != nil {
		fmt.Println(err)
	}
}

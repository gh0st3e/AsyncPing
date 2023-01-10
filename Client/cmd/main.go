package main

import (
	"Client/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.New()

	handler.Mount(server)

	server.Run(":8085")
}

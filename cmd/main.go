package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"runtime"
	"youtube_converter/handlers"
)

const DownloadFolder = "./downloads/"

func main() {
	if err := os.MkdirAll(DownloadFolder, 0755); err != nil {
		log.Fatalf("Failed to create download folder: %v", err)
	}

	maxProcess := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcess)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	router.Use(handlers.CORS())

	router.GET("/convert", handlers.ConvertHandler)
	router.GET("/download/:filename", handlers.DownloadFileHandler)
	router.GET("/search", handlers.SearchHandler)
	router.GET("/ws", handlers.HandleConnections)

	go handlers.HandleMessages()

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

package routes

import (
	"net/http"
	"uploadfilego/handlers"
)

func SetupRoutes() {
	http.HandleFunc("/upload", handlers.UploadFile)
}

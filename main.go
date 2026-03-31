package main

import (
	"fmt"
	"log"
	"net/http"
	"uploadfilego/routes"
)

func main() {
	routes.SetupRoutes()

	port := ":8081"
	fmt.Println("Server berjalan di port", port)
	fmt.Println("Buka browser dan akses: http://localhost:8081/upload")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

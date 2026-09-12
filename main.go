package main

import (
	"Mav_backend/api"
	"log"
)

func main() {
	r := api.SetupRouter()
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
	log.Println("Server started on port 8080")
}

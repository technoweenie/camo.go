package main

import (
	".."
	"fmt"
)

func main() {
	server := camo.Server(8080)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Could not start server:", err)
	}
}

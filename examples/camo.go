package main

import (
	".."
	"fmt"
)

func main() {
	server := camo.Server()

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Could not start server:", err)
	}
}

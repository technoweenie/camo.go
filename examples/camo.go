package main

import (
	".."
	"fmt"
)

func main() {
	server := camo.NewServer(8080)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Could not start server:", err)
	}
}

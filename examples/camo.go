package main

import (
	".."
	"fmt"
)

func main() {
	server := camo.NewServer()
	server.SetDigestKey("monkey")

	if err := server.ListenAndServe(8080); err != nil {
		fmt.Println("Could not start server:", err)
	}
}

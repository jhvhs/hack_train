package main

import (
	"fmt"
	"hack_train/server"
)

func main() {
	if err := server.NewServer().Run(); err != nil  {
		fmt.Printf("Shutting down: %s", err.Error())
	}
}

package main

import (
	"aggron/internal/config"
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	config.LoadEnvVariables()
}
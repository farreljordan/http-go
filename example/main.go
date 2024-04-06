package main

import (
	"fmt"
	"github.com/farreljordan/http-go"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	err := http.Listen()
	if err != nil {
		os.Exit(1)
	}
}

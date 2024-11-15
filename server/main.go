package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	fmt.Printf("\n🎮 Snake Game Server\n")
	fmt.Printf("🌐 Server running at: http://localhost:%s\n\n", "8080") //change this
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/yk-saito/proglog/internal/server"
)

func main() {
	// 8080ポートでHTTPサーバーを起動
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}

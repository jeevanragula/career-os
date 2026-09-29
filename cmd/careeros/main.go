package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jeevanragula/career-os/internal/httpapi"
)

func main() {
	addr := os.Getenv("CAREEROS_HTTP_ADDR")
	if addr == "" { addr = ":8080" }
	server := &http.Server{Addr: addr, Handler: httpapi.NewHandler()}
	log.Printf("CareerOS listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal(err) }
}

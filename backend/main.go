package main

import (
	"embed"
	"io/fs"
	"log"
	"strconv"

	"example.com/hydro-gate-monitor-service/api"
	"example.com/hydro-gate-monitor-service/config"
	"example.com/hydro-gate-monitor-service/store"
)

//go:embed web/*
var webFiles embed.FS

func main() {
	webFS, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	port := config.Port()
	log.Printf("hydro gate service listening on :%d", port)
	log.Fatal(serveAddress(":"+strconv.Itoa(port), api.NewRouter(store.New(), webFS)))
}

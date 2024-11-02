//go:build !dev
// +build !dev

package project

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/public/*
var publicFS embed.FS

func Public() http.Handler {
	fsys, err := fs.Sub(publicFS, "web/public")
	if err != nil {
		panic(err)
	}

	return http.StripPrefix("/public/", http.FileServer(http.FS(fsys)))
}

const IsDevelopment = false

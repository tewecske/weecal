//go:build !dev
// +build !dev

package project

import (
	"embed"
	"net/http"
)

//go:embed web/public
var publicFS embed.FS

func Public() http.Handler {
	return http.FileServerFS(publicFS)
}

const IsDevelopment = false

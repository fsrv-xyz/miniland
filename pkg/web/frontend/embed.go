package frontend

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var distFileSystem embed.FS

func DistFileSystem() fs.FS {
	contentStatic, _ := fs.Sub(fs.FS(distFileSystem), "dist")
	return contentStatic
}

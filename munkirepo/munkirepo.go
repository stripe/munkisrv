package munkirepo

import "embed"

//go:embed catalogs
//go:embed client_resources
//go:embed icons/*
//go:embed manifests
var Repo embed.FS

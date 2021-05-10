package resources

import "embed"

//go:embed json/*.json cookie/*.cookie
var FS embed.FS

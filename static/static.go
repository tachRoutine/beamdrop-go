package static

import "embed"

//go:embed all:frontend*
var FrontendFiles embed.FS

package templates

import "embed"

//go:embed index.html results.html
var FS embed.FS

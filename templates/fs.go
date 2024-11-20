package templates

import "embed"

//go:embed *
var FS embed.FS

// Initialising the file system by telling the compiler to embed the files in it.

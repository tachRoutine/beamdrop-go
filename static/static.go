package static

import (
	"embed"
	"fmt"
)

//go:embed all:frontend*
var FrontendFiles embed.FS

func init(){
	fmt.Println("Static package initialized")
	fmt.Println("Embedded files:", FrontendFiles)
}

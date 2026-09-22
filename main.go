// Command sysproxy provides the legacy standalone system-proxy interface.
//
// Deprecated: KokoroBox integrations must use KokoroBox Service instead.
package main

import (
	"fmt"
	"os"

	"github.com/UruhaLushia/sysproxy-go/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

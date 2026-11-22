// Command catty builds one Java application into one standalone Go executable.
package main

import (
	"os"

	"catty/internal/build"
)

func main() {
	os.Exit(build.Run(os.Args[1:]))
}

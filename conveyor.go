// Package main is a simple wrapper of the real conveyor entrypoint package.
//
// This package should NOT be extended or modified in any way; to modify the
// conveyor binary, work in the `gitlab.com/<USER>/conveyor/cmd` package.
//
package main

import (
	conveyor "github.com/junland/conveyor/cmd"
)

func main() {
	conveyor.Run()
}

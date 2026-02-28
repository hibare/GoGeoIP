package main

import (
	"github.com/hibare/Waypoint/cmd"
	"github.com/hibare/Waypoint/cmd/common"
)

func main() {
	common.Banner()
	cmd.Execute()
}

package main

import (
	"github.com/Louisrca/bloatfish/cmd"
	"github.com/Louisrca/bloatfish/internal/server"
)

func main() {

	server.StartServer()

	cmd.Execute()
}

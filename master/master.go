package main

import (
	"flag"
	"fmt"
	"terminal-bot/src/client"
	"terminal-bot/src/filestorage"
	"terminal-bot/src/server"
)

func main() {
	mode := flag.String("mode", "server", "choose mode")

	flag.Parse()

	switch *mode {
	case "server":
		{
			server.Run()
		}
	case "client":
		{
			client.Run()
		}
	case "filestorage":
		{
			filestorage.Run()
		}

	default:
		{
			fmt.Printf("Unknown mode\n")
			return
		}
	}
}

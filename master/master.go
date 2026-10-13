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
			server.Run() // FIXME: REVIEW, REWORK NEDDED
		}
	case "client":
		{
			client.Run() // FIXME: REVIEW, REWORK NEDDED
		}
	case "filestorage":
		{
			filestorage.Run() // FIXME: REVIEW, REWORK NEDDED
		}

	default:
		{
			fmt.Printf("Unknown mode\n")
			return
		}
	}
}

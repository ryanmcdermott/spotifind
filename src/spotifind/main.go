package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
)

// Globals to hold user-defined artists to search for and connect.
var origin string
var destination string

func main() {
	app := cli.NewApp()
	app.Name = "spotifind"
	app.Version = "0.1.0"
	app.Authors = []cli.Author{cli.Author{Name: "Ryan McDermott", Email: "ryan.mcdermott@ryansworks.com"}}
	app.Usage = "Find the degree of connection between two artists using Spotify's related artist feature."

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "origin, o", Value: "Led Zeppelin", Usage: "Artist to start from."},
		cli.StringFlag{Name: "destination, d", Value: "Taylor Swift", Usage: "Artist to end on."},
	}

	app.Action = func(c *cli.Context) {
		if len(c.String("origin")) == 0 {
			log.Fatal("Error, please provide both a starting and ending artist.")
		} else {
			origin = c.String("origin")
			destination = c.String("destination")
			FindFirstArtist(origin, destination)
		}
	}

	app.Run(os.Args)
}

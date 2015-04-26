package main

import (
	"encoding/json"
	"fmt"
	"github.com/rapito/go-spotify/spotify"
	"log"
	"os"
	"strings"
)

// Global variable to hold Spotify API request struct
var spot spotify.Spotify

// The channel that is written to when we find a link between an origin and
// destination artist. This halts the program, and walks the tree of the connected artists
// and displays it to the user.
var GlobalQuitChan = make(chan Artist)

// Artist type definition.
// Related holds a slice of children that are connected to Artist by similarity.
type Artist struct {
	Name   string
	Id     string
	Parent *Artist
}

// Structure to use while parsing JSON from server.
type Artists struct {
	Href  string
	Items []Artist
}

// Generic struct with a slice we pass around to all functions within the app.
// Holds each actual Artist struct, which are individually used as pointers between
// all the functions within the app.
type VisitedArtists struct {
	Artists []Artist
}

// Dummy struct to hold the "Search For Artist" response from Spotify's API.
type JsonObject struct {
	Artists Artists
}

// Dummy struct to hold an "Artist's Related Artists", as delivered by Spotify's API.
type RelatedArtists struct {
	Artists []Artist
}

// Don't try to search through an artist we have already seen before. Search VisitedArtists.artists whenever
// we are adding a new artist -- O(n).
func IsArtistVisited(id string, visited *VisitedArtists) bool {
	for _, artist := range visited.Artists {
		if artist.Id == id {
			return true
		}
	}

	return false
}

// Append an Artist struct to a VisitedArtists.Artists slice.
func (visited *VisitedArtists) AddArtist(artist Artist) []Artist {
	visited.Artists = append(visited.Artists, artist)
	return visited.Artists
}

// Look through Spotify's API response of related artists and build each Artist struct to add
// to VisitedArtists struct and send each struct to our workers to begin looking for more related Artists.
func FindRelatedArtists(artist *Artist, visited *VisitedArtists) {
	results, err := spot.Get("artists/%s/related-artists", nil, artist.Id)

	if err != nil {
		log.Fatal(err)
	}

	var response RelatedArtists
	json.Unmarshal(results, &response)

	for _, relatedArtist := range response.Artists {
		if !IsArtistVisited(relatedArtist.Id, visited) {
			relatedArtist.Parent = artist
			visited.AddArtist(relatedArtist)

			visitedArtist := visited.Artists[len(visited.Artists)-1]
			work := WorkRequest{Artist: &visitedArtist, Visited: visited}

			if strings.ToLower(visitedArtist.Name) == strings.ToLower(destination) {
				GlobalQuitChan <- visitedArtist
				break
			} else {
				WorkQueue <- work
			}
		}
	}
}

// Walk the tree that we created of an artist back up through each parent until we get
// to the origin artist we started with.
func WalkRelatedArtists(artist *Artist, tree []string) []string {
	tree = append(tree, artist.Name)
	if artist.Parent != nil {
		tree = WalkRelatedArtists(artist.Parent, tree)
	}

	return tree
}

// The main app function and listener that waits for an artist match and then
// walks the tree of parent artists to display them to the user.
func Spotifind() {
	for {
		select {
		case connectedArtist := <-GlobalQuitChan:
			// Walk the tree and determine who the destination artist is connected to
			tree := WalkRelatedArtists(&connectedArtist, nil)

			fmt.Printf("The artist, %s, is connected to the artist, %s, by this chain: \n", connectedArtist.Name, origin)
			for _, artist := range tree {
				fmt.Println(artist)
				fmt.Println("\n")
			}

			return
		}
	}
}

// Initialize the Spotify API and begin the workers to start searching for our destination artist.
func FindFirstArtist(origin string, destination string) {
	client_id := os.Getenv("SPOTIFY_CLIENT_ID")
	client_secret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if len(client_id) == 0 || len(client_secret) == 0 {
		log.Fatal("Unable to connect to Spotify API. No SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET has been set in your environment variables.")
	}

	spot = spotify.New(client_id, client_secret)

	origin = strings.Replace(strings.ToLower(origin), " ", "+", -1)

	result, err := spot.Get("search/?q=%s&type=artist", nil, origin)

	if err != nil {
		fmt.Println(err)
		return
	}

	var response JsonObject
	json.Unmarshal(result, &response)

	firstArtist := response.Artists.Items[0]

	var visited = VisitedArtists{}
	visited.AddArtist(firstArtist)

	// Start 4 workers to handle Spotify API calls.
	StartDispatcher(4)

	work := WorkRequest{Artist: &visited.Artists[0], Visited: &visited}
	WorkQueue <- work

	Spotifind()
}

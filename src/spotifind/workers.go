package main

import (
	"fmt"
)

// Simple structure to hold a reference to an artist.
type WorkRequest struct {
	Artist  *Artist
	Visited *VisitedArtists
}

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

var WorkerQueue chan chan WorkRequest

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
}

// Create the Dispatcher that starts our workers
func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				// fmt.Println("Received work request")
				go func() {
					worker := <-WorkerQueue

					// fmt.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
	}

	return worker
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w Worker) Start() {
	go func() {
		for {
			// Add work into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				// fmt.Printf("worker%d: Received work request", w.ID)
				fmt.Printf("worker %d: %s\n", w.ID, work.Artist.Name)
				FindRelatedArtists(work.Artist, work.Visited)

			case <-GlobalQuitChan:
				// We have been asked to stop.
				fmt.Printf("worker %d stopping\n", w.ID)
				return
			}
		}
	}()
}

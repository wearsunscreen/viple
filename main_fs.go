package main

// Description: This is the main file for the file server. It creates a file server
// that serves the contents of the specified directory and starts the server on port 8080.
//
// To use,
// - put in a separate directory,
// - uncomment code below
// - rename main_fs.go to main.go,
// - rename main_fs function to main,
// - rename the package to main
// - set the directory to serve in the http.Dir() function.

// package main_fs

// import (
// 	"log"
// 	"net/http"
// )

// func main_fs() {
// 	// Create a file server that serves the contents of the specified directory.
// 	fileServer := http.FileServer(http.Dir("../viple"))

// 	// Register the file server as the handler for all requests.
// 	http.Handle("/", fileServer)

// 	// Start the server on port 8080.
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

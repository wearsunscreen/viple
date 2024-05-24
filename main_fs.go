// Description: This is a simple file server which is needed to run the page with
// wasm locally, the browser will not load wasm files.
// To use,
// - put in a separate directory,
// - uncomment code below
// - rename main_fs.go to main.go,
// - rename main_fs function to main,
// - set the directory to serve in the http.Dir() function.
// - run go mod init <any_module_name>
// - run go mod tidy
// - run the program.
package main

import (
	"log"
	"net/http"
)

func main_fs() {
	// Create a file server that serves the contents of the specified directory.
	fileServer := http.FileServer(http.Dir("../viple"))

	// Register the file server as the handler for all requests.
	http.Handle("/", fileServer)

	// Start the server on port 8080.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

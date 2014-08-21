// This package implements a simple HTTP server providing a REST API to a node handler.
//
// It provides four methods:
//
// GET /node/ Retrieves all the nodes.
// POST /node/ Creates a new node given a title.
// GET /node/{nodeID} Retrieves the node with the given id.
// PUT /node/{nodeID} Updates the node with the given id.
//
// Every method below gives more information about every API call, its parameters, and its results.
package server

import (
	"encoding/json"
	"fmt"
	"github.com/tracyde/hiera-gui/node"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var nodes = node.NewNodeManager()

const PathPrefix = "/node/"

func RegisterHandlers() {
	r := mux.NewRouter()
	r.HandleFunc(PathPrefix, errorHandler(ListNodes)).Methods("GET")
	r.HandleFunc(PathPrefix, errorHandler(NewNode)).Methods("POST")
	r.HandleFunc(PathPrefix+"{name}", errorHandler(GetNode)).Methods("GET")
	r.HandleFunc(PathPrefix+"{name}", errorHandler(UpdateNode)).Methods("PUT")
	http.Handle(PathPrefix, r)
}

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

// errorHandler wraps a function returning an error by handling the error and returning a http.Handler.
// If the error is of the one of the types defined above, it is handled as described for every type.
// If the error is of another type, it is considered as an internal error and its message is logged.
func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case notFound:
			http.Error(w, "node not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}

// ListNode handles GET requests on /node.
// There's no parameters and it returns an object with a Nodes field containing a list of nodes.
//
// Example:
//
// req: GET /node/
// res: 200 {"Nodes": [
// {"ID": 1, "Title": "Learn Go", "Done": false},
// {"ID": 2, "Title": "Buy bread", "Done": true}
// ]}
func ListNodes(w http.ResponseWriter, r *http.Request) error {
	res := struct{ Nodes []*node.Node }{nodes.All()}
	return json.NewEncoder(w).Encode(res)
}

// NewNode handles POST requests on /node.
// The request body must contain a JSON object with a Title field.
// The status code of the response is used to indicate any error.
//
// Examples:
//
// req: POST /node/ {"Title": ""}
// res: 400 empty title
//
// req: POST /node/ {"Title": "Buy bread"}
// res: 200
func NewNode(w http.ResponseWriter, r *http.Request) error {
	req := struct{ Name string; Role string }{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequest{err}
	}
	t, err := node.NewNode(req.Name, req.Role)
	if err != nil {
		return badRequest{err}
	}
	return nodes.Save(t)
}

// parseName obtains the id variable from the given request url,
// parses the obtained text and returns the result.
func parseName(r *http.Request) (string, error) {
	txt, ok := mux.Vars(r)["name"]
	if !ok {
		return "", fmt.Errorf("node name not found")
	}
	return txt, nil
}

// GetNode handles GET requsts to /node/{nodeID}.
// There's no parameters and it returns a JSON encoded node.
//
// Examples:
//
// req: GET /node/1
// res: 200 {"ID": 1, "Title": "Buy bread", "Done": true}
//
// req: GET /node/42
// res: 404 node not found
func GetNode(w http.ResponseWriter, r *http.Request) error {
	name, err := parseName(r)
	log.Println("Node is ", name)
	if err != nil {
		return badRequest{err}
	}
	n, ok := nodes.Find(name)
	log.Println("Found", ok)
	if !ok {
		return notFound{}
	}
	return json.NewEncoder(w).Encode(n)
}

// UpdateNode handles PUT requests to /node/{nodeID}.
// The request body must contain a JSON encoded node.
//
// Example:
//
// req: PUT /node/1 {"ID": 1, "Title": "Learn Go", "Done": true}
// res: 200
//
// req: PUT /node/2 {"ID": 2, "Title": "Learn Go", "Done": true}
// res: 400 inconsistent node IDs
func UpdateNode(w http.ResponseWriter, r *http.Request) error {
	name, err := parseName(r)
	if err != nil {
		return badRequest{err}
	}
	var n node.Node
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		return badRequest{err}
	}
	if n.Name != name {
		return badRequest{fmt.Errorf("inconsistent node names")}
	}
	if _, ok := nodes.Find(name); !ok {
		return notFound{}
	}
	return nodes.Save(&n)
}

package hello

import (
	"fmt"
	"net/http"
)

// Subscribe ...
func Subscribe(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "someone"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

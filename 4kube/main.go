package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const countLimit = 5

var count int
var version *int
var errHand *bool

func handlerHostName(w http.ResponseWriter, r *http.Request) {
	if *errHand {
		count++
		if count > countLimit {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}
	}

	name, err := os.Hostname()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	fmt.Fprintf(w, name)
	fmt.Fprintf(w, "\nversion ")
	fmt.Fprintf(w, strconv.Itoa(*version))
}

func main() {
	version = flag.Int("v", 1, "version as an int")
	errHand = flag.Bool("err", false, "err as a bool")
	flag.Parse()

	http.HandleFunc("/", handlerHostName)
	http.ListenAndServe("localhost:8080", nil)
}

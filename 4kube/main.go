package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

	//	if err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		w.Write([]byte("500 - Something bad happened!"))
	//		return
	//	}

	if r.Method == "POST" {
		f, err := os.Create("/tmp/4kube.txt")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		nb, err := f.Write(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("wrote %d bytes\n", nb)
		f.Sync()
		f.Close()

	} else {
		name, err := os.Hostname()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}
		fmt.Fprintf(w, name)
		fmt.Fprintf(w, "\nversion ")
		fmt.Fprintf(w, strconv.Itoa(*version))
		fmt.Fprintf(w, "\n")

		data, err := ioutil.ReadFile("/tmp/4kube.txt")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
		fmt.Fprintf(w, "\n")

	}
}

func main() {
	version = flag.Int("v", 1, "version as an int")
	errHand = flag.Bool("err", false, "err as a bool")
	flag.Parse()

	http.HandleFunc("/", handlerHostName)
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

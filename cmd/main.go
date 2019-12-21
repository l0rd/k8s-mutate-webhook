package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	m "github.com/alex-leonhardt/k8s-mutate-webhook/pkg/mutate"
	v "github.com/alex-leonhardt/k8s-mutate-webhook/pkg/validate"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello %q", html.EscapeString(r.URL.Path))
}

func handleMutate(w http.ResponseWriter, r *http.Request) {

	// read the body / request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	// mutate the request
	mutated, err := m.Mutate(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	// and write it back
	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}

func handleValidate(w http.ResponseWriter, r *http.Request) {

	// read the body / request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	// validate the request
	mutated, err := v.Validate(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	// and write it back
	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/mutate", handleMutate)
	mux.HandleFunc("/validate", handleValidate)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	log.Fatal(s.ListenAndServeTLS("./ssl/mutateme.pem", "./ssl/mutateme.key"))

}

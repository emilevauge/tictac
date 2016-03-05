package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Messages struct {
	Boom   time.Time `json:"boom,omitempty"`
	Traces []string  `json:"traces,omitempty"`
}

func main() {
	http.HandleFunc("/", tictacHandler)
	http.HandleFunc("/health", healthHandler)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "OK")
}

func tictacHandler(w http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(req.URL.String())
	fmt.Printf("Incoming request %+v\n", req)
	messages := &Messages{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(messages)
	switch {
	case err == io.EOF:
		// empty body
		queryParams := u.Query()
		tictac := queryParams.Get("tictac")
		if len(tictac) > 0 {
			duration, err := time.ParseDuration(tictac)
			if err != nil {
				fmt.Println("Error", err)
				w.WriteHeader(http.StatusBadRequest)
				req.Write(w)
				return
			}
			messages = &Messages{
				Boom: time.Now().Add(duration),
			}
		} else {
			fmt.Println("Error missing tictac!")
			w.WriteHeader(http.StatusBadRequest)
			req.Write(w)
			return
		}
	case err != nil:
		log.Println("Error decoding messages", err)
		w.WriteHeader(http.StatusBadRequest)
		req.Write(w)
		return
	}
	time.Sleep(1000 * time.Millisecond)
	// do I have to explode?
	if messages.Boom.Before(time.Now()) {
		json.NewEncoder(w).Encode(*messages)
		go func() {
			time.Sleep(3000 * time.Millisecond)
			log.Fatal("Booooooooooooooooooooom !")
		}()
	} else {
		hostname, _ := os.Hostname()
		messages.Traces = append(messages.Traces, hostname+" - "+time.Now().String())
		request := gorequest.New()
		_, body, errs := request.Post("http://"+u.Host+":8080").
			Set("Host", u.Host).
			Send(*messages).
			End()
		if len(errs) > 0 {
			fmt.Println("Error", errs)
			w.WriteHeader(http.StatusBadRequest)
			req.Write(w)
			return
		}

		fmt.Fprint(w, body)
	}
}

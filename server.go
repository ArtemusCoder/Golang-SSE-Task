package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var lastMessage Message

type Message struct {
	Word string `json:"word"`
}

func ListenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(lastMessage)
	fmt.Fprintf(w, "data: %v\n\n", buf.String())
	fmt.Fprintf(w, "retry: 1000\n")
	fmt.Printf("data: %v\n", buf.String())

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

}

func SayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var jsonString map[string]interface{}
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.Unmarshal(b, &jsonString)
		if word, found := jsonString["word"]; found {
			var id string
			var ok bool
			if id, ok = word.(string); !ok {
				http.Error(w, "WORD field is not string", http.StatusBadRequest)
				return
			} else {
				message := Message{Word: id}
				lastMessage = message
			}
		} else {
			http.Error(w, "No WORD field", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Only POST", http.StatusBadRequest)
	}
}

func main() {
	lastMessage = Message{Word: "firstWord"}
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.HandleFunc("/listen", ListenHandler)
	http.HandleFunc("/say", SayHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

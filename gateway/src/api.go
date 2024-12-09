package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	utils "github.com/adamlahbib/go-ms-poc/common"
	spec "github.com/adamlahbib/go-ms-poc/spec"
)

// channels to wait for reply messages from rabbit
var rchans = make(map[string](chan spec.CreateDocumentReply))

func initApi() {
	// router
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/documents", apiDocument).Methods("POST")

	log.Printf("INFO: Starting HTTP API")

	// start server
	err := http.ListenAndServe(":7654", r)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

func apiDocument(w http.ResponseWriter, r *http.Request) {
	// read body
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	log.Printf("INFO: Received request: %s", string(data))

	// unmarshal request and create document message
	doc := &spec.Document{}
	err := json.Unmarshal(data, doc)
	if err != nil {
		log.Printf("ERROR: Failed to unmarshal request: %s", err.Error())
		response(w, "Invalid request JSON", http.StatusBadRequest)
	}

	docMsg := &spec.CreateDocumentMessage{
		Uid:      utils.Uid(),
		Document: doc,
		ReplyTo:  "gateway",
	}

	log.Printf("INFO: Document Message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan spec.CreateDocumentReply)
	rchans[docMsg.Uid] = rchan

	// send message to rabbit
	msg := RabbitMsg{
		QueueName: "storage",
		Message:   *docMsg,
	}

	pchan <- msg
	waitReply(docMsg.Uid, rchan, w)
}

func waitReply(uid string, rchan chan spec.CreateDocumentReply, w http.ResponseWriter) {
	for {
		select {
		case docReply := <-rchan:
			// response received
			log.Printf("INFO: Received reply: %v with uid: %s", docReply, uid)

			// send response back to client
			response(w, "Created", http.StatusCreated)

			// remove channel from rchans
			delete(rchans, uid)
			return
		case <-time.After(10 * time.Second):
			// timeout
			log.Printf("ERROR: Timeout waiting for reply with uid: %s", uid)

			// send response to client
			response(w, "Timeout", http.StatusRequestTimeout)

			// remove channel from rchans
			delete(rchans, uid)
			return
		}
	}
}

func response(
	w http.ResponseWriter,
	resp string,
	status int,
) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, resp)
}

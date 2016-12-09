package subscriber

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

var requestIsInProgress bool
var mutex sync.Mutex

func init() {
	http.HandleFunc("/log-pubsub-message", logPubSubMessage)

	mutex = sync.Mutex{}
}

func logPubSubMessage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	log.Infof(ctx, "Starting the execution of logPubSubMessage")

	if requestIsInProgress {
		log.Errorf(ctx, "A request is already in progress, ealy exiting")
		http.Error(w, "A request is already in progress, ealy exiting", 500)
		return
	}

	mutex.Lock()

	requestIsInProgress = true

	log.Infof(ctx, "Just locked mutex")

	defer func() {
		log.Infof(ctx, "Unlocking mutex")
		requestIsInProgress = false
		mutex.Unlock()
	}()

	log.Infof(ctx, "Starting working on pubsub message")

	s, err := getStringValueFromPubSubPayload(ctx, r.Body)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to read message from payload: %v", err)
		http.Error(w, errorMsg, 500)
		log.Errorf(ctx, errorMsg)
		return
	}

	log.Infof(ctx, "Start sleeping for 8 seconds...")

	time.Sleep(8 * time.Second)

	log.Infof(ctx, "... Waking up")

	log.Infof(ctx, "New message from pubsub: %s", s)

	fmt.Fprint(w, "Message successfully logged")
}

func getStringValueFromPubSubPayload(ctx context.Context, r io.Reader) (string, error) {
	p := payload{}
	if err := json.NewDecoder(r).Decode(&p); err != nil {
		return "", err
	}

	log.Infof(ctx, "PubSub message id: %v", p.Message.ID)

	s := p.Message.Attributes["value"]
	if s == "" {
		return "", errors.New("Message value is empty")
	}

	return s, nil
}

type payload struct {
	Message struct {
		Attributes map[string]string
		ID         string `json:"message_id"`
	}
}

package publisher

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const (
	PROJECT_ID  = "pubsub-throttle"
	TOPIC_ID    = "messages-to-log"
	NB_MESSAGES = 100
)

func init() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/publish", publish)
}

func publish(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	messages := []*pubsub.Message{}
	for i := 0; i < NB_MESSAGES; i++ {
		s := randStringRunes(15)
		messages = append(messages, &pubsub.Message{Attributes: map[string]string{"value": s}})
	}

	if err := publishMessages(ctx, messages); err != nil {
		http.Error(w, fmt.Sprintf("Failed to publish messages to PubSub: %v", err), 500)
		log.Errorf(ctx, "Failed to publish messages to PubSub: %v", err)
		return
	}

	fmt.Fprintf(w, "%v messages successfully published", NB_MESSAGES)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func publishMessages(ctx context.Context, messages []*pubsub.Message) error {
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return err
	}

	topic := client.Topic(TOPIC_ID)
	_, err = topic.Publish(ctx, messages...)

	return err
}

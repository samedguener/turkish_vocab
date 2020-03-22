package hello

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
)

// SubscriptionRequest ...
type SubscriptionRequest struct {
	Email string `json:"email"`
}

// Subscribe ...
func Subscribe(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var msg SubscriptionRequest
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Use the application default credentials
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "staging-270414"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	_, _, err = client.Collection("subscriptions").Add(ctx, map[string]interface{}{
		"email": msg.Email,
	})
	if err != nil {
		log.Fatalf("Failed adding subscription: %v", err)
	}

	output, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

package hello

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
)

type Status struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

// SubscriptionRequest ...
type SubscriptionRequest struct {
	Email               string `json:"email"`
	IsGrammarInterested bool   `json:"grammarInterested"`
	IsCultureInterested bool   `json:"cultureInterested"`
	IsHistoryInterested bool   `json:"historyInterested"`
	IsFoodInterested    bool   `json:"foodInterested"`
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

	subscribed, err := isSubscribed(msg.Email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if subscribed {
		status := Status{Status: "Failed", Description: "Email is already subscribed!"}

		output, err := json.Marshal(status)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)
		return
	}

	if err = saveSubscription(msg); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	status := Status{Status: "Success", Description: "Email is successfully subscribed!"}

	output, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
	return

}

func saveSubscription(data SubscriptionRequest) error {
	projectID := os.Getenv("GCP_PROJECT")

	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Collection("subscriptions").Doc(data.Email).Set(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func isSubscribed(email string) (bool, error) {
	projectID := os.Getenv("GCP_PROJECT")

	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return false, err
	}
	defer client.Close()

	documentRef := client.Collection("subscription").Doc(email)
	if documentRef != nil {
		return true, nil
	}

	return false, nil
}

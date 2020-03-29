package hello

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Status ...
type Status struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

// SubscriptionRequest ...
type SubscriptionRequest struct {
	Email     *string   `json:"email"`
	Interests []*string `json:"interests"`
}

// SubscriptionConfirmationRequest ...
type SubscriptionConfirmationRequest struct {
	Email *string `json:"email"`
	Hash  *string `json:"hash"`
}

// ConfirmSubscription ...
func ConfirmSubscription(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	var msg SubscriptionConfirmationRequest
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if msg.Email == nil {
		err = fmt.Errorf("missing email")
		log.Printf(err.Error())

		w.Header().Set("content-type", "application/json")
		status := Status{Status: "failure", Description: "Malformed request! Missing key email!"}
		output, err := json.Marshal(status)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(400)
		w.Write(output)
		return
	}

	if msg.Hash == nil {
		err = fmt.Errorf("missing interests")
		log.Printf(err.Error())

		w.Header().Set("content-type", "application/json")
		status := Status{Status: "failure", Description: "Malformed request! Missing key hash!"}
		output, err := json.Marshal(status)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(400)
		w.Write(output)
		return
	}

	subscribed, err := isSubscribed(*msg.Email)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if !subscribed {
		status := Status{Status: "failure", Description: "Email is not subscribed!"}

		output, err := json.Marshal(status)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)
	}

}

// Subscribe ...
func Subscribe(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	var msg SubscriptionRequest
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if msg.Email == nil {
		err = fmt.Errorf("missing email")
		log.Printf(err.Error())

		w.Header().Set("content-type", "application/json")
		status := Status{Status: "failure", Description: "Malformed request! Missing key email!"}
		output, err := json.Marshal(status)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(400)
		w.Write(output)
		return
	}

	if msg.Interests == nil {
		err = fmt.Errorf("missing interests")
		log.Printf(err.Error())

		w.Header().Set("content-type", "application/json")
		status := Status{Status: "failure", Description: "Malformed request! Missing key interests!"}
		output, err := json.Marshal(status)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(400)
		w.Write(output)
		return
	}

	_, err = mail.ParseAddress(*msg.Email)
	if err != nil {
		log.Printf(err.Error())

		status := Status{Status: "failure", Description: "Malformed request! Malformed email!"}
		output, errTwo := json.Marshal(status)
		if errTwo != nil {
			log.Printf(err.Error())
			http.Error(w, errTwo.Error(), 500)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(404)
		w.Write(output)
		return
	}

	subscribed, err := isSubscribed(*msg.Email)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if subscribed {
		status := Status{Status: "failure", Description: "Email is already subscribed!"}

		output, err := json.Marshal(status)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)
		return
	}

	if err = saveSubscription(msg); err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	status := Status{Status: "success", Description: "Email is successfully subscribed!"}

	output, err := json.Marshal(status)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
	return

}

func saveSubscription(data SubscriptionRequest) error {
	projectID, isSet := os.LookupEnv("GCP_PROJECT")
	if !isSet {
		return fmt.Errorf("environment variable 'GCP_PROJECT' not set")
	}

	conf := &firebase.Config{ProjectID: projectID}
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Collection("subscriptions").Doc(*data.Email).Set(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func deleteSubscription(email string) error {
	projectID, isSet := os.LookupEnv("GCP_PROJECT")
	if !isSet {
		return fmt.Errorf("environment variable 'GCP_PROJECT' not set")
	}

	conf := &firebase.Config{ProjectID: projectID}
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Collection("subscriptions").Doc(email).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func isSubscribed(email string) (bool, error) {
	projectID, isSet := os.LookupEnv("GCP_PROJECT")
	if !isSet {
		return false, fmt.Errorf("environment variable 'GCP_PROJECT' not set")
	}

	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return false, nil
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return false, err
	}
	defer client.Close()

	_, err = client.Collection("subscriptions").Doc(email).Get(ctx)
	if grpc.Code(err) == codes.NotFound {
		return false, nil
	}

	return true, nil
}

package hello

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Test struct {
	request            string
	expectedResponse   string
	expectedStatusCode int
	cleanUp            bool
}

func TestSubscribe(t *testing.T) {

	tests := []Test{
		{request: `{ "email" : "malformed", "interests" : [ "food", "grammar", "culture", "vocab"] }`, expectedResponse: `{ "status" : "failure", "description" : "Malformed request! Malformed email!" }`, expectedStatusCode: 404, cleanUp: false},
		{request: `{ "interests" : [ "food", "grammar", "culture", "vocab"] }`, expectedResponse: `{ "status" : "failure", "description" : "Malformed request! Missing key email!" }`, expectedStatusCode: 400, cleanUp: false},
		{request: `{ "email" : "mail@samedguener.com" }`, expectedResponse: `{ "status" : "failure", "description" : "Malformed request! Missing key interests!" }`, expectedStatusCode: 400, cleanUp: false},
		{request: `{ "email" : "mail@samedguener.com", "interests" : [ "food", "grammar", "culture", "vocab"] }`, expectedResponse: `{ "status" : "success", "description" : "Email is successfully subscribed!" }`, expectedStatusCode: 200, cleanUp: false},
		{request: `{ "email" : "mail@samedguener.com", "interests" : [ "food", "grammar", "culture", "vocab"] }`, expectedResponse: `{ "status" : "failure", "description" : "Email is already subscribed!" }`, expectedStatusCode: 200, cleanUp: true},
	}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/subscriptions", strings.NewReader(test.request))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Subscribe)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != test.expectedStatusCode {
			t.Errorf("Status code differs. Expected %d .\n Got %d instead", test.expectedStatusCode, status)
		}

		var resp Status
		b, err := ioutil.ReadAll(rr.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(b, &resp)
		if err != nil {
			t.Fatal(err)
		}

		var expectedResp Status
		err = json.Unmarshal([]byte(test.expectedResponse), &expectedResp)
		if err != nil {
			t.Fatal(err)
		}

		if resp.Status != expectedResp.Status {
			t.Errorf("Status differs. Expected '%s' .\n Got '%s' instead", expectedResp.Status, resp.Status)
		}

		if resp.Description != expectedResp.Description {
			t.Errorf("Description differs. Expected '%s' .\n Got '%s' instead", expectedResp.Description, resp.Description)
		}

		if test.cleanUp {
			var msg SubscriptionRequest

			b, err := ioutil.ReadAll(strings.NewReader(test.request))
			if err != nil {
				t.Fatal(err)
			}

			err = json.Unmarshal(b, &msg)
			if err != nil {
				t.Fatal(err)
			}

			err = deleteSubscription(*msg.Email)
			if err != nil {
				t.Fatal(err)
			}

		}
	}

}

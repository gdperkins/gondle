package main

import (
	"errors"
	"net/http"

	"io/ioutil"

	"encoding/json"

	"fmt"

	"github.com/google/go-github/github"
)

const (
	githubEventHeader      = "X-GitHub-Event"
	githubSignatureHeader  = "X-GitHub-Signature"
	githubDeliveryIDHeader = "X-GitHub-Delivery"
)

// HookContext represents the http request from
// the webhook source
type HookContext struct {
	Signature        string
	Payload          []byte
	DeliveryID       string
	Event            string
	PullRequestEvent github.PullRequestEvent
}

// HookResponse gets returned to the client
type HookResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// HandleHook takes a request from the source and trys to process it
func HandleHook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	e := json.NewEncoder(w)

	hook, err := parseGitHook(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e.Encode(HookResponse{"Invalid Request", 500})
		return
	}

	if err := tryDeploy(hook); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e.Encode(HookResponse{"Invalid Request", 500})

	} else {
		w.WriteHeader(http.StatusOK)
		e.Encode(HookResponse{"Success", 200})
	}
}

// Parses the request and builds a HookRequest struct to be used
// for deploying
func parseGitHook(r *http.Request) (*HookContext, error) {
	hc := HookContext{}

	if hc.Signature = r.Header.Get(githubSignatureHeader); len(hc.Signature) == 0 {
		return nil, errors.New("no signature")
	}

	if hc.Event = r.Header.Get(githubEventHeader); len(hc.Event) == 0 {
		return nil, errors.New("no event")
	}

	if hc.DeliveryID = r.Header.Get(githubDeliveryIDHeader); len(hc.DeliveryID) == 0 {
		return nil, errors.New("no delivery id")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	pr := github.PullRequestEvent{}
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, errors.New("invalid payload")
	}
	hc.PullRequestEvent = pr

	if !verifySignature("") {
		return nil, errors.New("invalid signature")
	}

	return &hc, nil
}

func verifySignature(s string) bool {
	return true
}

func tryDeploy(h *HookContext) error {

	fmt.Printf("Merged: %T", h.PullRequestEvent.PullRequest.Merged)

	//pull out config items here and execute
	fmt.Println(h.PullRequestEvent)

	return nil
}

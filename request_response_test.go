package hashdb

import (
	"testing"
)

func TestRequestResponse(t *testing.T) {
	response := make(chan *GetResponse, 1)
	request := &GetRequest{"ping", response}

	go func() {
		request.response <- &GetResponse{"ping", "pong", nil}
	}()

	resp := <-response
	if resp.err != nil || resp.password != "pong" {
		t.Fatal("Error wrong response values!")
	}
}

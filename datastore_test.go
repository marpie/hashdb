// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashdb

import (
	"strings"
	"testing"
)

func TestDatastoreOpen(t *testing.T) {
	_, err := OpenDatastore(":memory:", 0)
	if err != nil {
		t.Fatalf("Error while opening datastore: %s", err)
	}
}

func TestDatastorePut(t *testing.T) {
	ds, err := OpenDatastore(":memory:", 1)
	if err != nil {
		t.Fatalf("Opening database error: %s", err)
	}

	response := ds.Put("5f87e0f786e60b554ec522ce85ddc930", "MT1992")
	if response.err != nil {
		t.Fatalf("Put failed: %s", response.err)
	}

	count, err := ds.Count()
	if err != nil {
		t.Fatalf("Count() failed: %s", err)
	}
	if count != 1 {
		t.Fatalf("Count should be 1 but got: %d", count)
	}

	password, err := ds.GetExact("5f87e0f786e60b554ec522ce85ddc930")
	if err != nil {
		t.Fatalf("Get failed: %s", err)
	}

	if password != "MT1992" {
		t.Fatalf("Password doesn't match: '%s' vs. 'MT1992'", password)
	}
}

func TestDatastorePut2(t *testing.T) {
	ds, _ := OpenDatastore(":memory:", 2)

	response := ds.Put("world!", "hello")
	if response.err != nil {
		t.Fatalf("Put failed: %s", response.err)
	}

	response = ds.Put("World!", "Hello")
	if response.err != nil {
		t.Fatalf("Put failed: %s", response.err)
	}

	count, err := ds.Count()
	if err != nil {
		t.Fatalf("Count() failed: %s", err)
	}
	if count != 2 {
		t.Fatalf("Count should be 2 but got: %d", count)
	}

	password, _ := ds.GetExact("world!")
	if password != "hello" {
		t.Fatalf("Password doesn't match: '%s' vs. 'hello'", password)
	}

	password, _ = ds.GetExact("World!")
	if password != "Hello" {
		t.Fatalf("Password doesn't match: '%s' vs. 'Hello'", password)
	}
}

func TestDatastoreGetLike(t *testing.T) {
	ds, _ := OpenDatastore(":memory:", 2)

	response := ds.Put("world!", "hello")
	if response.err != nil {
		t.Fatalf("Put failed: %s", response.err)
	}

	response = ds.Put("World", "Hello")
	if response.err != nil {
		t.Fatal("Put failed: %s", response.err)
	}

	request := &GetRequest{request: "world%", response: make(chan *GetResponse)}
	ds.GetLike(request)

	for {
		resp, ok := <-request.response
		if !ok {
			break
		}
		if resp.err != nil {
			t.Fatalf("Like failed with error: %s", resp.err)
			continue
		}

		if !strings.HasPrefix(strings.ToLower(resp.hash), "world") {
			t.Fatalf("Hash should begin with 'world' but got '%s'", resp.hash)
			continue
		}

		if strings.ToLower(resp.password) != "hello" {
			t.Fatalf("Password should be 'hello' but got: '%s'", resp.password)
			continue
		}
	}
}

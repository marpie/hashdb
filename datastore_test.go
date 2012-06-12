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
		return
	}

	err = ds.Put("5f87e0f786e60b554ec522ce85ddc930", "MT1992")
	if err != nil {
		t.Fatalf("Put failed: %s", err)
		return
	}

	password, err := ds.GetExact("5f87e0f786e60b554ec522ce85ddc930")
	if err != nil {
		t.Fatalf("Get failed: %s", err)
		return
	}

	if password != "MT1992" {
		t.Fatalf("Password doesn't match: '%s' vs. 'MT1992'", password)
		return
	}
}

func TestDatastorePut2(t *testing.T) {
	ds, _ := OpenDatastore(":memory:", 2)

	err := ds.Put("world!", "hello")
	if err != nil {
		t.Fatalf("Put failed: %s", err)
		return
	}

	err = ds.Put("World!", "Hello")
	if err != nil {
		t.Fatalf("Put failed: %s", err)
		return
	}

	password, _ := ds.GetExact("world!")
	if password != "hello" {
		t.Fatalf("Password doesn't match: '%s' vs. 'hello'", password)
		return
	}

	password, _ = ds.GetExact("World!")
	if password != "Hello" {
		t.Fatalf("Password doesn't match: '%s' vs. 'Hello'", password)
		return
	}
}

func TestDatastoreGetLike(t *testing.T) {
	ds, _ := OpenDatastore(":memory:", 2)

	err := ds.Put("world!", "hello")
	if err != nil {
		t.Fatalf("Put failed: %s", err)
		return
	}

	err = ds.Put("World", "Hello")
	if err != nil {
		t.Fatal("Put failed: %s", err)
		return
	}

	request := &GetRequest{request: "world%", response: make(chan *GetResponse)}
	ds.GetLike(request)

	for {
		resp, ok := <-request.response
		if !ok {
			break
		}
		if resp.err != nil {
			t.Fatalf("Like failed with error: %s", err)
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

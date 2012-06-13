// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashdb

import (
	"crypto/md5"
	"testing"
)

func TestOpenDatabase(t *testing.T) {
	_, err := OpenDatabase(":memory:", md5.New(), 0)
	if err != nil {
		t.Fatalf("OpenDatabase failed with error: %s", err)
	}
}

func TestHashDbPut(t *testing.T) {
	db, _ := OpenDatabase(":memory:", md5.New(), 1)

	// Create communication channels
	respChan := make(chan *PutResponse, 1)
	redirectChan := make(chan *PutResponse, 1)
	// Set up "redirect"
	go func(respChan chan *PutResponse, redirectChan chan *PutResponse) {
		response := <-respChan
		redirectChan <- response
	}(respChan, redirectChan)

	db.Put("MT1992", respChan)

	// wait for redirect
	response := <-redirectChan

	if response.err != nil {
		t.Fatalf("Put failed with error: %s", response.err)
	}

	count, err := db.Count()
	if err != nil {
		t.Fatalf("Count() failed with error: %s", err)
	}
	if count != 1 {
		t.Fatalf("Count() should be 1 but got: %d", count)
	}
}

func TestHashDbGetExact(t *testing.T) {
	db, _ := OpenDatabase(":memory:", md5.New(), 1)

	respChan := make(chan *PutResponse, 1)
	db.Put("MT1992", respChan)
	<-respChan // discard result

	password, err := db.GetExact("5f87e0f786e60b554ec522ce85ddc930")
	if err != nil {
		t.Fatalf("GetExact failed with error: %s", err)
	}

	if password != "MT1992" {
		t.Fatalf("Wrong Password - '%s' vs. 'MT1992'", password)
	}
}

func TestHashDbGetLike(t *testing.T) {
	db, _ := OpenDatabase(":memory:", md5.New(), 1)

	respChan := make(chan *PutResponse, 1)
	db.Put("MT1992", respChan)
	<-respChan // discard result
	db.Put("SAUL69", respChan)
	<-respChan

	result, err := db.GetLike("e0f7")
	if err != nil {
		t.Fatalf("GetLike failed with error: %s", err)
	}

	if len(result) != 1 {
		t.Fatalf("There should be one item in the map but got '%d'.", len(result))
	}
}

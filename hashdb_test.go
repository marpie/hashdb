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

	err := db.Put("MT1992")
	if err != nil {
		t.Fatalf("Put failed with error: %s", err)
	}
}

func TestHashDbGetExact(t *testing.T) {
	db, _ := OpenDatabase(":memory:", md5.New(), 1)

	db.Put("MT1992")

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
	db.Put("MT1992")
	db.Put("SAUL69")

	result, err := db.GetLike("e0f7")
	if err != nil {
		t.Fatalf("GetLike failed with error: %s", err)
	}

	if len(result) != 1 {
		t.Fatalf("There should be one item in the map but got '%d'.", len(result))
	}
}

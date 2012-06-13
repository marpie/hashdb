// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashdb

import (
	"testing"
)

func TestOpenMix(t *testing.T) {
	_, err := OpenMix(":memory:", 0)
	if err != nil {
		t.Fatalf("OpenMix failed with error: %s", err)
	}
}

func TestHashMixPut(t *testing.T) {
	mix, _ := OpenMix(":memory:", 1)

	responseChan := make(chan *PutResponse, 2)

	mix.Put("MT1992", responseChan)
	res1 := <-responseChan
	res2 := <-responseChan

	if res1.err != nil {
		t.Fatalf("Put() failed: %s", res1.err)
	}
	if res2.err != nil {
		t.Fatalf("Put() failed: %s", res2.err)
	}

	md5, sha1, err := mix.Count()
	if err != nil {
		t.Fatalf("Count failed with error: %s", err)
	}
	if md5 != 1 {
		t.Fatalf("Count(MD5) should be 1 but got: %d", md5)
	}
	if sha1 != 1 {
		t.Fatalf("Count(SHA1) should be 1 but got: %d", sha1)
	}
}

func TestHashMixGetMD5(t *testing.T) {
	mix, _ := OpenMix(":memory:", 1)

	responseChan := make(chan *PutResponse, 2)
	mix.Put("MT1992", responseChan)
	<-responseChan // discard
	<-responseChan // results

	passwordMap, err := mix.GetMD5("5f87e0f786e60b554ec522ce85ddc930")
	if err != nil {
		t.Fatalf("GetExact failed with error: %s", err)
	}

	password, ok := passwordMap["5f87e0f786e60b554ec522ce85ddc930"]
	if !ok {
		t.Fatalf("Hash not found!")
	}

	if password != "MT1992" {
		t.Fatalf("Wrong Password - '%s' vs. 'MT1992'", password)
	}
}

func TestHashMixGetMD52(t *testing.T) {
	mix, _ := OpenMix(":memory:", 1)

	responseChan := make(chan *PutResponse, 2)
	mix.Put("MT1992", responseChan)
	<-responseChan // discard
	<-responseChan // results

	passwordMap, err := mix.GetMD5("0b554ec")
	if err != nil {
		t.Fatalf("GetExact failed with error: %s", err)
	}

	for hash, password := range passwordMap {
		if hash != "5f87e0f786e60b554ec522ce85ddc930" {
			t.Fatalf("Hash not found!")
		}
		if password != "MT1992" {
			t.Fatalf("Wrong Password - '%s' vs. 'MT1992'", password)
		}
	}
}

func TestHashMixGetSHA1(t *testing.T) {
	mix, _ := OpenMix(":memory:", 1)

	responseChan := make(chan *PutResponse, 2)
	mix.Put("MT1992", responseChan)
	<-responseChan // discard
	<-responseChan // results

	passwordMap, err := mix.GetSHA1("f8f446125d304bf227f770d81f0a429f9fe9b025")
	if err != nil {
		t.Fatalf("GetExact failed with error: %s", err)
	}

	password, ok := passwordMap["f8f446125d304bf227f770d81f0a429f9fe9b025"]
	if !ok {
		t.Fatalf("Hash not found!")
	}

	if password != "MT1992" {
		t.Fatalf("Wrong Password - '%s' vs. 'MT1992'", password)
	}
}

func TestHashMixGetSHA12(t *testing.T) {
	mix, _ := OpenMix(":memory:", 1)

	responseChan := make(chan *PutResponse, 2)
	mix.Put("MT1992", responseChan)
	<-responseChan // discard
	<-responseChan // results

	passwordMap, err := mix.GetSHA1("d81f0a429f9")
	if err != nil {
		t.Fatalf("GetExact failed with error: %s", err)
	}

	for hash, password := range passwordMap {
		if hash != "f8f446125d304bf227f770d81f0a429f9fe9b025" {
			t.Fatalf("Hash not found!")
		}
		if password != "MT1992" {
			t.Fatalf("Wrong Password - '%s' vs. 'MT1992'", password)
		}
	}
}

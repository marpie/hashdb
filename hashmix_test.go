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

	err := mix.Put("MT1992")
	if err != nil {
		t.Fatalf("Put failed with error: %s", err)
	}
}

func TestHashMixGetMD5(t *testing.T) {
	mix, _ := OpenMix(":memory:", 1)

	mix.Put("MT1992")

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

	mix.Put("MT1992")

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

	mix.Put("MT1992")

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

	mix.Put("MT1992")

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

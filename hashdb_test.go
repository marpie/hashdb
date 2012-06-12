package hashdb

import (
	"testing"
)

func TestOpenDatabase(t *testing.T) {
	_, err := OpenDatabase(":memory:", 0)
	if err != nil {
		t.Fatalf("OpenDatabase failed with error: %s", err)
	}
}

func TestHashDbPut(t *testing.T) {
	db, _ := OpenDatabase(":memory:", 1)

	err := db.Put("5f87e0f786e60b554ec522ce85ddc930", "MT1992")
	if err != nil {
		t.Fatalf("Put failed with error: %s", err)
	}
}

func TestHashDbGetExact(t *testing.T) {
	db, _ := OpenDatabase(":memory:", 1)

	db.Put("5f87e0f786e60b554ec522ce85ddc930", "MT1992")

	password, err := db.GetExact("5f87e0f786e60b554ec522ce85ddc930")
	if err != nil {
		t.Fatalf("GetExact failed with error: %s", err)
	}

	if password != "MT1992" {
		t.Fatalf("Wrong Password - '%s' vs. 'MT1992'", password)
	}
}

func TestHashDbGetLike(t *testing.T) {
	db, _ := OpenDatabase(":memory:", 1)
	db.Put("5f87e0f786e60b554ec522ce85ddc930", "MT1992")
	db.Put("9f5e22b402096fa932298322e43fc3ca", "SAUL69")

	result, err := db.GetLike("e0f7")
	if err != nil {
		t.Fatalf("GetLike failed with error: %s", err)
	}

	if len(result) != 1 {
		t.Fatalf("There should be one item in the map but got '%d'.", len(result))
	}
}

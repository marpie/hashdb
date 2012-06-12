// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
This http-"server" understands three different requests:
  - /md5?q=AAAAAAA      -> searches the MD5 hash in the database
  - /sha1?q=BBBBBB      -> searches the SHA1 hash in the database
  - /new?p=newPassword  -> adds the password to the database

Before the server is started the first time the "db" folder has to be created:
  - `mkdir db db/md5 db/sha1`
*/
package main

import (
	"github.com/marpie/hashdb"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "[HashDBServer] ", log.LstdFlags)

	logger.Println("[*] Initializing Hash-Database")
	hashMix, err := hashdb.OpenMix("db", 10)
	if err != nil {
		logger.Println("Error initializing HashMix:", err)
		return
	}
	ghh, err := NewGetHashHandler(hashMix, logger)
	if err != nil {
		logger.Println("Error initializing GetHashHandler:", err)
		return
	}
	phh, err := NewPutHashHandler(hashMix, logger)
	if err != nil {
		logger.Println("Error initializing PutHashHandler:", err)
		return
	}

	logger.Println("[*] Waiting for clients...")
	http.Handle("/md5", ghh)
	http.Handle("/sha1", ghh)
	http.Handle("/new", phh)
	// Redirect all other requests to...
	http.Handle("/", http.RedirectHandler("https://github.com/marpie/hashdb", 301))

	http.ListenAndServe(":8080", nil)
}

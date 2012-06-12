// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/marpie/hashdb"
	"log"
	"net/http"
	"strings"
)

type PutHashHandler struct {
	hashMix *hashdb.HashMix
	logger  *log.Logger
}

func NewPutHashHandler(hashMix *hashdb.HashMix, logger *log.Logger) (phh *PutHashHandler, err error) {
	phh = new(PutHashHandler)
	phh.hashMix = hashMix
	phh.logger = logger
	return phh, nil
}

func (phh *PutHashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	if err := r.ParseForm(); err != nil {
		LogError(phh.logger, w, r, "ParseForm() error:", err)
		return
	}

	parts := strings.SplitN(r.URL.Path, "/", 3)
	if len(parts) < 2 {
		LogError(phh.logger, w, r, "Error while parsing Path: "+r.URL.Path, nil)
		return
	}
	method := strings.ToLower(parts[1])
	if method != "new" {
		return
	}

	password := r.Form.Get("q")
	if password == "" {
		return
	}

	err := phh.hashMix.Put(password)
	if err != nil {
		LogError(phh.logger, w, r, "Error while adding new hash:", err)
		return
	}

	fmt.Fprintf(w, htmlTemplate(password, "<li>Done.</li>"))
	return
}

// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"github.com/marpie/hashdb"
	"log"
	"net/http"
	"strings"
)

type GetHashHandler struct {
	hashMix *hashdb.HashMix
	logger  *log.Logger
}

func NewGetHashHandler(hashMix *hashdb.HashMix, logger *log.Logger) (ghh *GetHashHandler, err error) {
	ghh = new(GetHashHandler)
	ghh.hashMix = hashMix
	ghh.logger = logger
	return ghh, nil
}

func (ghh *GetHashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		LogError(ghh.logger, w, r, "ParseForm() error:", err)
		return
	}

	parts := strings.SplitN(r.URL.Path, "/", 3)
	if len(parts) < 2 {
		LogError(ghh.logger, w, r, "Error while parsing Path: "+r.URL.Path, nil)
		return
	}
	method := strings.ToLower(parts[1])
	md5 := method == "md5"
	sha1 := method == "sha1"

	hash := r.Form.Get("q")
	if hash == "" {
		return
	}

	var passesMap map[string]string
	var err error
	if md5 {
		passesMap, err = ghh.hashMix.GetMD5(hash)
	} else if sha1 {
		passesMap, err = ghh.hashMix.GetSHA1(hash)
	} else {
		LogError(ghh.logger, w, r, "No hash function match!", nil)
		return
	}

	if err != nil {
		LogError(ghh.logger, w, r, "Error while searching for hash:", err)
		return
	}

	passes := new(bytes.Buffer)
	for key, value := range passesMap {
		passes.WriteString(htmlFormatHash(key, value))
	}

	fmt.Fprintf(w, htmlTemplate(hash, passes.String()))
	return
}

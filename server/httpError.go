// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
)

func ReturnError(w http.ResponseWriter) {
	w.Write([]byte("Error while processing request!"))
	return
}

func LogRequest(logger *log.Logger, r *http.Request, logString string) {
	logger.Printf("[Client-IP: %s] [Request-URI: %s] %s", r.RemoteAddr, r.RequestURI, logString)
}

func LogError(logger *log.Logger, w http.ResponseWriter, r *http.Request, errString string, err error) {
	LogRequest(logger, r, fmt.Sprintf("%s: %s", errString, err))
	ReturnError(w)
}

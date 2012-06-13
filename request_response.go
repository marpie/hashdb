// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashdb

type GetResponse struct {
	hash     string
	password string
	err      error
}

type GetRequest struct {
	request  string
	response chan *GetResponse
}

type PutResponse struct {
	password string
	hash     string
	err      error
}

type PutRequest struct {
	hash     string
	password string
	response chan *PutResponse
}

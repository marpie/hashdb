// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashdb

import (
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

const HexAlphabet = "0123456789abcdef"

type LookupTable map[string]*Datastore

type HashDB struct {
	directory   string
	lookupTable LookupTable
	hashFunc    hash.Hash
}

// OpenDatabase creates or opens a new HashDB instance. directory accepts a
// valid directory or ":memory:" if you want the database only in RAM.
func OpenDatabase(directory string, hashFunc hash.Hash, maxGetHandler int) (db *HashDB, err error) {
	if directory != ":memory:" {
		// check if directory exists
		fi, err := os.Stat(directory)
		if err != nil {
			return nil, err
		}
		if !fi.IsDir() {
			return nil, ErrDirectoryNotFound
		}
	}

	db = new(HashDB)
	db.lookupTable = make(LookupTable)
	db.hashFunc = hashFunc

	// build lookup table
	for _, x := range HexAlphabet {
		for _, y := range HexAlphabet {
			var ds *Datastore
			var err error

			name := fmt.Sprintf("%c%c", x, y)
			filename := name + ".sqlite3"
			if directory == ":memory:" {
				ds, err = OpenDatastore(":memory:", maxGetHandler)
			} else {
				ds, err = OpenDatastore(filepath.Join(directory, filename), maxGetHandler)
			}
			if err != nil {
				return nil, err
			}
			db.lookupTable[name] = ds
		}
	}

	return db, nil
}

// Count returns the number of entries in the database.
func (db *HashDB) Count() (result int64, err error) {
	for _, ds := range db.lookupTable {
		count, err := ds.Count()
		if err != nil {
			return -1, err
		}
		result += count
	}
	return result, nil
}

func (db *HashDB) getDatastoreByHash(hash string) (ds *Datastore, err error) {
	if len(hash) < 2 {
		return nil, ErrHashTooShort
	}

	key := hash[0:2]
	ds, ok := db.lookupTable[key]
	if !ok {
		return nil, ErrDatastoreNotFound
	}

	return ds, nil
}

func (db *HashDB) getHash(password string) string {
	db.hashFunc.Reset()
	io.WriteString(db.hashFunc, password)
	return fmt.Sprintf("%x", db.hashFunc.Sum(nil))
}

// Put stores a new hash in the database. The hash is stored as lower-case.
func (db *HashDB) Put(password string, resultChan chan *PutResponse) {
	if password == "" {
		resultChan <- &PutResponse{password: password, hash: "", err: ErrPasswordMissing}
		return
	}

	// Fire up a new Goroutine to handle the database insert.
	go func(password string, resultChan chan *PutResponse) {
		// Calculate hash
		hash := db.getHash(password)

		ds, err := db.getDatastoreByHash(hash)
		if err != nil {
			resultChan <- &PutResponse{password: password, hash: hash, err: err}
			return
		}
		resultChan <- ds.Put(hash, password)
	}(password, resultChan)

	return
}

// GetExact searches for the hash in the database. 
// The hash parameter has to be lower-case.
func (db *HashDB) GetExact(hash string) (password string, err error) {
	ds, err := db.getDatastoreByHash(hash)
	if err != nil {
		return "", err
	}
	return ds.GetExact(hash)
}

// GetLike searches in the database for a hash that contains the supplied
// hash-part.
func (db *HashDB) GetLike(hash string) (map[string]string, error) {
	respMap := make(map[string]string)

	for _, ds := range db.lookupTable {
		request := &GetRequest{request: fmt.Sprintf("%%%s%%", hash), response: make(chan *GetResponse)}
		ds.GetLike(request)

		for {
			resp, ok := <-request.response
			if !ok {
				break
			}
			if resp.err != nil {
				return nil, resp.err
			}
			respMap[resp.hash] = resp.password
		}
	}

	return respMap, nil
}

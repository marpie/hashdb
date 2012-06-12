// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package hashdb provides a database for hex-based hashes.

This package needs kuroneko's sqlite3 library
(https://github.com/kuroneko/gosqlite3).


Sample

  // Create a new database in memory.
  db, err := OpenDatabase(":memory:", 10)
  if err != nil {
    return err
  }

  // Add new entry
  err = db.Put("5f87e0f786e60b554ec522ce85ddc930", "MT1992")
  if err != nil {
    return err
  }

  // Get entry
  password, err := db.GetExact("5f87e0f786e60b554ec522ce85ddc930")
  if err != nil {
    return err
  }

*/
package hashdb

import (
	"fmt"
	"os"
	"path/filepath"
)

const HexAlphabet = "0123456789abcdef"

type LookupTable map[string]*Datastore

type HashDB struct {
	directory   string
	lookupTable LookupTable
}

// OpenDatabase creates or opens a new HashDB instance. directory accepts a
// valid directory or ":memory:" if you want the database only in RAM.
func OpenDatabase(directory string, maxGetHandler int) (db *HashDB, err error) {

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

// Put stores a new hash in the database.
func (db *HashDB) Put(hash string, password string) PutResponse {
	ds, err := db.getDatastoreByHash(hash)
	if err != nil {
		return err
	}
	return ds.Put(hash, password)
}

// GetExact searches for the hash in the database.
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

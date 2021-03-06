// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashdb

import (
	"crypto/md5"
	"crypto/sha1"
	"path/filepath"
	"strings"
)

type HashMix struct {
	md5  *HashDB
	sha1 *HashDB
}

// OpenMix creates a new HashMix instance.
func OpenMix(directory string, maxPerDatabaseHandler int) (hm *HashMix, err error) {
	hm = new(HashMix)
	md5_path := directory
	sha1_path := directory
	if directory != ":memory:" {
		md5_path = filepath.Join(directory, "md5")
		sha1_path = filepath.Join(directory, "sha1")
	}
	hm.md5, err = OpenDatabase(md5_path, md5.New(), maxPerDatabaseHandler)
	if err != nil {
		return nil, err
	}
	hm.sha1, err = OpenDatabase(sha1_path, sha1.New(), maxPerDatabaseHandler)
	if err != nil {
		return nil, err
	}
	return hm, nil
}

// Count returns the number of entries in the databases.
func (hm *HashMix) Count() (md5_count int64, sha1_count int64, err error) {
	md5_count, err = hm.md5.Count()
	if err != nil {
		return -1, -1, err
	}
	sha1_count, err = hm.sha1.Count()
	return md5_count, sha1_count, err
}

// Put stores the new password in the MD5 and SHA1 database. 
// The responseChan gets two messages - one for each HashDB.
func (hm *HashMix) Put(password string, responseChan chan *PutResponse) {
	hm.md5.Put(password, responseChan)
	hm.sha1.Put(password, responseChan)
}

func (hm *HashMix) getHash(db *HashDB, hash string) (passes map[string]string, err error) {
	// Since all hashes are stored lower-case:
	hash = strings.ToLower(hash)

	if len(hash) != (db.hashFunc.Size() * 2) {
		if (len(hash) < 5) || (len(hash) > db.hashFunc.Size()*2) {
			return nil, ErrWrongHashFormat
		}
		// Like-search
		return db.GetLike(hash)
	}

	// Exact-search
	passes = make(map[string]string)
	password, err := db.GetExact(hash)
	if err != nil {
		return nil, err
	}
	passes[hash] = password
	return passes, nil
}

// GetMD5 tries to find the MD5-hash in the database.
func (hm *HashMix) GetMD5(hash string) (map[string]string, error) {
	return hm.getHash(hm.md5, hash)
}

// GetSHA1 tries to find the SHA1-hash in the database.
func (hm *HashMix) GetSHA1(hash string) (map[string]string, error) {
	return hm.getHash(hm.sha1, hash)
}

// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashdb

import (
	sqlite "github.com/kuroneko/gosqlite3"
)

type Datastore struct {
	db        *sqlite.Database
	putChan   chan *PutRequest
	exactChan chan *GetRequest
	likeChan  chan *GetRequest
}

// OpenDatastore creates a new Datastore and binds it to the backend database.
func OpenDatastore(name string, maxGetHandler int) (ds *Datastore, err error) {
	db, err := sqlite.Open(name)
	if err != nil {
		return nil, err
	}

	ds = new(Datastore)
	ds.db = db

	err = ds.initializeDatabase()
	if err != nil {
		return nil, err
	}

	err = ds.initializeHandler(maxGetHandler)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

// Put creates a new entry in the database.
func (ds *Datastore) Put(hash string, password string) PutResponse {
	resp_chan := make(chan PutResponse, 1)
	ds.putChan <- &PutRequest{hash: hash, password: password, response: resp_chan}
	return <-resp_chan
}

// GetExact returns the password for the supplied hash. No fuzzy search!
func (ds *Datastore) GetExact(hash string) (password string, err error) {
	resp_chan := make(chan *GetResponse, 1)
	ds.exactChan <- &GetRequest{request: hash, response: resp_chan}
	resp := <-resp_chan
	return resp.password, resp.err
}

// GetLike returns all matching passwords for the supplied hash.
//   hash LIKE ?
func (ds *Datastore) GetLike(request *GetRequest) {
	ds.likeChan <- request
	return
}

// initializeDatabase checks if the database already exists otherwise it 
// creates it. Additionally it also sets up the prepared statements.
func (ds *Datastore) initializeDatabase() (err error) {
	// check if table(s) exist
	_, err = ds.db.Execute("SELECT COUNT(*) FROM data;")
	if err != nil {
		err = ds.createDatabase()
		if err != nil {
			return err
		}
	}

	// create the request channel
	ds.exactChan = make(chan *GetRequest)
	ds.likeChan = make(chan *GetRequest)
	ds.putChan = make(chan *PutRequest)

	return nil
}

// getExactHandler is used as a Goroutine to handle the exact searches.
func getExactHandler(ds *Datastore, initChan chan error) {
	statement, err := ds.setupStatement("SELECT password FROM data WHERE hash = ?")
	if err != nil {
		initChan <- err
		return
	}
	initChan <- nil

	for {
		statement.ClearBindings()

		request, running := <-ds.exactChan
		if !running {
			break
		}

		err, _ := statement.BindAll(request.request)
		if err != nil {
			request.response <- &GetResponse{request.request, "", err}
			continue
		}

		err = statement.Step()
		if err != sqlite.ROW {
			request.response <- &GetResponse{request.request, "", err}
			continue
		}

		request.response <- &GetResponse{request.request, statement.Column(0).(string), nil}
		statement.Reset()
	}
}

// getLikeHandler is used as a Goroutine to handle the "fuzzy" searches.
func getLikeHandler(ds *Datastore, initChan chan error) {
	statement, err := ds.setupStatement("SELECT hash, password FROM data WHERE hash LIKE ?")
	if err != nil {
		initChan <- err
		return
	}
	initChan <- nil

	for {
		statement.ClearBindings()

		request, running := <-ds.likeChan
		if !running {
			break
		}

		err, _ := statement.BindAll(request.request)
		if err != nil {
			request.response <- &GetResponse{request.request, "", err}
			continue
		}

		for err = statement.Step(); err == sqlite.ROW; err = statement.Step() {
			request.response <- &GetResponse{statement.Column(0).(string), statement.Column(1).(string), nil}
		}
		statement.Reset()
		close(request.response)
	}
}

// putHandler is used as a Goroutine to add new entries to the database.
func putHandler(ds *Datastore, initChan chan error) {
	statement, err := ds.setupStatement("INSERT INTO data (hash, password) VALUES (?, ?);")
	if err != nil {
		initChan <- err
		return
	}
	initChan <- nil

	for {
		statement.Reset()

		request, running := <-ds.putChan
		if !running {
			break
		}

		err, _ := statement.BindAll(request.hash, request.password)
		if err != nil {
			request.response <- err
			continue
		}

		err = statement.Step()
		if err != nil {
			request.response <- err
			continue
		}

		request.response <- nil
	}
}

func (ds *Datastore) initializeHandler(maxGetHandler int) (err error) {
	init_chan := make(chan error, 1)

	for i := 0; i < maxGetHandler; i++ {
		go getExactHandler(ds, init_chan)
		if <-init_chan != nil {
			return ErrInitExactHandlerFailed
		}

		go getLikeHandler(ds, init_chan)
		if <-init_chan != nil {
			return ErrInitLikeHandlerFailed
		}
	}

	go putHandler(ds, init_chan)
	if <-init_chan != nil {
		return ErrInitPutHandlerFailed
	}
	return nil
}

// createDatabase creates the table and index.
func (ds *Datastore) createDatabase() (err error) {
	err = ds.db.Begin()
	if err != nil {
		return
	}

	_, err = ds.db.Execute("CREATE TABLE data (hash TEXT, password TEXT);")
	if err != nil {
		ds.db.Rollback()
		return ErrCreatingTable
	}

	_, err = ds.db.Execute("CREATE UNIQUE INDEX idx_hash ON data(hash ASC);")
	if err != nil {
		ds.db.Rollback()
		return ErrCreatingTable
	}

	err = ds.db.Commit()
	return
}

// setupStatement creates a new prepared statement.
func (ds *Datastore) setupStatement(sql string) (st *sqlite.Statement, err error) {
	st, err = ds.db.Prepare(sql, "")
	if err != nil {
		return nil, err
	}
	err = st.Reset()
	if err != nil {
		return nil, err
	}

	return st, nil
}

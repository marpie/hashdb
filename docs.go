// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package hashdb provides a database for hex-based hashes.

This package needs kuroneko's sqlite3 library
(https://github.com/kuroneko/gosqlite3).


HashMix Sample

  // Create a new hash-mix in memory.
  mix, err := OpenMix(":memory:", 10)
  if err != nil {
    return err
  }

  // Add new entry
  err = mix.Put("MT1992")
  if err != nil {
    return err
  }

  // Get entry
  passwordMap, err := mix.GetMD5("5f87e0f786e60b554ec522ce85ddc930") // or mix.GetMD5("5f87e0f")
  if err != nil {
    return err
  }


HashDB Sample

  // Create a new database in memory.
  db, err := OpenDatabase(":memory:", md5.New(), 10)
  if err != nil {
    return err
  }

  // Add new entry
  err = db.Put("MT1992")
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

hashdb
======

Package hashdb provides a database for hex-based hashes.

Sample
------
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


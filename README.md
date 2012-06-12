hashdb
======

Package hashdb provides a database for hex-based hashes.

Sample
------

HashDatabase
```go
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
```

HashMix
```go
  // Create a new hash-mix in memory.
  mix, err := OpenMix(":memory:", 1)
  if err != nil {
      return err
  }
  
  // Add new entry
  err = mix.Put("MT1992")
  if err != nil {
      return err
  }
  
  // Get entry
  passMap, err := mix.GetMD5("5f87e0f786e60b554ec522ce85ddc930") // or mix.GetMD5("5f87e0f")
  if err != nil {
      return err
  }
```

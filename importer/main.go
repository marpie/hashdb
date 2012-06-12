// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "fmt"
  "os"
  "strings"
  "syscall"
  "runtime"
  "time"
  "github.com/marpie/hashdb"
)

func Usage() {
  fmt.Printf("Usage:\n\t%s [OutputDirectory] [InputFilename]\n", os.Args[0])
  syscall.Exit(0)
}

func FileExists(filename string) bool {
  fi, err := os.Stat(filename)
  if err != nil {
    return false
  }
  if fi.Mode()&os.ModeType == 0 {
    return true
  }
  return false
}

func DirectoryExists(directory string) bool {
  fi, err := os.Stat(directory)
  if err != nil {
    return false
  }
  if fi.IsDir() {
    return true
  }
  return false
}

func DatabaseImport(mix *hashdb.HashMix, passes chan string, result chan bool) {
  for {
    password, ok := <-passes
    if !ok {
      result <- true
      return
    }

    password_lower := strings.ToLower(password)
    password_upper := strings.ToUpper(password)
    password_capitalize := strings.Title(password_lower)

    if password_lower == password || password_lower == password_upper || password_lower == password_capitalize {
      password_lower = ""
    }

    if password_upper == password || password_upper == password_capitalize {
      password_upper = ""
    }

    if password_capitalize == password {
      password_capitalize = ""
    }

    if err := mix.Put(password); err != nil {
      println(password, "->", err.Error())
    }

    if password_lower != "" {
      if err := mix.Put(password_lower); err != nil {
        println(password_lower, "->", err.Error())
      }
    }

    if password_upper != "" {
      if err := mix.Put(password_upper); err != nil {
        println(password_upper, "->", err.Error())
      }
    }

    if password_capitalize != "" {
      if err := mix.Put(password_capitalize); err != nil {
        println(password_capitalize, "->", err.Error())
      }
    }
  }
  result <- true
}

func StatusUpdate(mix *hashdb.HashMix) {
  c := time.Tick(1 * time.Minute)
  for now := range c {
    md5, sha1, err := mix.Count()
    if err != nil {
      fmt.Printf("[Status %v] Error:", now, err)
      continue
    }
    fmt.Printf("[Status %v] MD5: %d - SHA1: %d", now, md5, sha1)
  }
}

func main() {
  if len(os.Args) < 3 {
    Usage()
    return
  }

  outputDirectory := os.Args[1]
  inputFilename := os.Args[2]

  if !DirectoryExists(outputDirectory) || !FileExists(inputFilename) {
    Usage()
    return
  }

  println("[*] Initializing Import...")
  mix, err := hashdb.OpenMix("db", 1)
  if err != nil {
    println("Error initializing database:", err)
    syscall.Exit(1)
  }

  resultChan := make(chan bool, 1)
  passes := make(chan string, 1)

  go DatabaseImport(mix, passes, resultChan)

  println("[*] Import started...")
  if err := ReadFileByLine(inputFilename, passes); err != nil {
    println("Error while importing:", err)
    syscall.Exit(1)
    return
  }

  go StatusUpdate(mix)

  if err := <-resultChan; !err {
    println("Error while importing to database:", err)
    syscall.Exit(1)
    return
  }

  println("[X] Done.")
}


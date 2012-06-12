// Copyright 2012 marpie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "bufio"
  "bytes"
  "io"
  "os"
)

func ReadFileByLine(fileName string, outputChan chan string) (err error) {
  file, err := os.Open(fileName)
  if err != nil {
    return
  }

  reader := bufio.NewReader(file)
  buffer := bytes.NewBuffer(make([]byte, 0))
  for {
    part, prefix, err := reader.ReadLine()
    if err != nil {
      break
    }
    if !prefix {
      if buffer.Len() == 0 {
        outputChan <- string(part)
      } else {
        buffer.Write(part)
        outputChan <- buffer.String()
        buffer.Reset()
      }
    } else {
      // more data available
      buffer.Write(part)
    }
  }
  close(outputChan)
  if err == io.EOF {
    return nil
  }
  return
}


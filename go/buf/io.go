// Copyright 2010 Mark C. Chu-Carroll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// File: io.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: implementations of IO methods for file-based buffers.


package buf

import (
  "io/ioutil"
  "os"
)



type FileManagerImpl struct {
  buffers map[string] *Buffer
}

func fileExists(filename string) bool {
  fi, err := os.Stat(filename)
  if err == nil {
    return fi.IsRegular()
  } else {
    return false	
  }
  return true
}

func (self *GapBuffer) Write() ResultCode {
  if !self.dirty {
    return SUCCEEDED
  }
  if fileExists(self.filename + ".bak") {
    os.Remove(self.filename + ".bak")
  }
  if fileExists(self.filename) {
    os.Rename(self.filename, self.filename + ".bak")
  }
  bytes := self.Bytes()
  err := ioutil.WriteFile(self.filename, bytes, 0544)
  if err != nil {
    return IO_ERROR	
  }
  self.dirty = false
  return SUCCEEDED
}

func (self *GapBuffer) Read() (code ResultCode, err string) {
  code, err = self.ReadFrom(self.filename)
  return
}

func (self *GapBuffer) ReadFrom(filename string) (code ResultCode, err string) {
  self.Clear()
  contents, oserr := ioutil.ReadFile(filename)
  if oserr != nil {
	// TODO: need more specific errors - use os.Error code
	// to generate some more specific error description.
	err = "OS Error reading file"
	code = IO_ERROR
	return
  } else {
    self.InsertChars(contents)
  }
  code = SUCCEEDED
  err = ""
  return
}


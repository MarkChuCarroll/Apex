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

// File: symtable.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: symbol table and bindings for the apex interpreter.


package language

type SymbolTable interface {
  Get(name string) (value string)
  Set(name string, value string)
}

type Binding struct {
  Name string
  Value string
  Next *Binding
}


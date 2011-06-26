// Copyright 2009 Mark C. Chu-Carroll
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

// File: actions.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Operations that perform buffer modifications
//    using a specified selection.

package actions

import (
  "buf"
  "expressions"
)


type Action interface {
  Execute(r expressions.Range) buf.Status
}

type DeleteAction struct{}

func NewDeleteAction() *DeleteAction { return &DeleteAction{} }

func (self *DeleteAction) Execute(r expressions.Range) buf.Status {
  status := r.Validate()
  if status.GetResultCode() != buf.SUCCEEDED {
    return status
  }
  buffer := r.GetBuffer()
  buffer.MoveTo(r.GetStart().GetAbsolute())
  buffer.Cut(r.GetLength())
  return buf.NewSuccess()
}

// Copyright 2011 Mark C. Chu-Carroll
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

// File: ast.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: ASTs for the Apex programming language

package language

type NodeType int32

const (
  NODE_APPEND_STR NodeType = iota
  NODE_APPEND_EXPR
  NODE_INSERT_STR
  NODE_INSERT_EXPR
  NODE_OPEN
  NODE_COPY
  NODE_DELETE
  NODE_MOVE
  NODE_JUMP
  NODE_LOOP
  NODE_GLOBAL
  NODE_REPLACE
  NODE_REPLACE_EXPR
  NODE_EXECUTE
  NODE_WRITE
  NODE_TYPE
  NODE_REVERT
  NODE_INVOKE
  NODE_FUN
  NODE_BLOCK
  NODE_ASSIGN
  NODE_FROMEXEC
  NODE_TOEXEC
  NODE_SEARCH
  NODE_RE_STR
  NODE_RE_CHOICE
  NODE_RE_CHARSET
  NODE_RE_REPEAT
  NODE_RE_GROUP
  NODE_RE_BIND
  NODE_RE_SEQ

  NODE_ERROR = -1
)

type AstNode interface {
  GetAstNodeType() NodeType
}

type ExpressionNode interface {
  AstNode
}

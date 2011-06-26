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

type Token struct { }

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
  NODE_COND

  NODE_ERROR = -1
)

type SourceLoc interface {
  GetSourceFile() string
  GetLine() string
}

type AstNode struct {
  nodetype NodeType
  line int32
  col int32
  left []AstNode
  mid []AstNode
  right []AstNode
}

func NewAstNode(t NodeType) *AstNode {
  result := new(AstNode)
  result.nodetype = t
  result.left = nil
  result.right = nil
  return result
}

func MakeCmdNode(t NodeType, left []AstNode, right []AstNode) *AstNode {
  result := NewAstNode(t)
  result.left = left
  result.right = right
  return result
}


// For conditionals, left = cond, mid = then, right = else.

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

// File: ast_cmd.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: ASTs for the Apex programming language commands

package language

type CommandNode interface {
  AstNode
}

type AppendStrCommand struct {
  str string
}

func (self *AppendStrCommand) GetAstNodeType() NodeType {
  return NODE_APPEND_STR
}

type AppendExprCommand struct {
  expr ExpressionNode
}

func (self *AppendExprCommand) GetAstNodeType() NodeType {
  return NODE_APPEND_EXPR
}

type InsertStringCommand struct {
  str string
}

func (self *InsertStringCommand) GetAstNodeType() NodeType {
  return NODE_INSERT_STR
}

type InsertExprCommand struct {
  expr ExpressionNode
}

func (self *InsertExprCommand) GetAstNodeType() NodeType {
  return NODE_INSERT_EXPR
}

type OpenFileCommand struct {
  filename string
}

func (self *OpenFileCommand) GetAstNodeType() NodeType {
  return NODE_OPEN
}

// Needs nothing: copies current selection
type CopyCommand struct {
  variable *string
}

func (self *CopyCommand) GetAstNodeType() NodeType {
  return NODE_COPY
}

type DeleteCommand struct {
  variable *string
}

func (self *DeleteCommand) GetAstNodeType() NodeType {
  return NODE_DELETE
}

type PositionUnit int
const (
  UNIT_CHAR PositionUnit = iota
  UNIT_LINE
  UNIT_PAGE
)

type MoveCommand struct {
  unit PositionUnit
  dist ExpressionNode
  extend bool
}

func (self *MoveCommand) GetAstNodeType() NodeType {
  return NODE_MOVE	
}

type JumpCommand struct {
  Unit PositionUnit
  dist ExpressionNode
  extend bool
}

func (self *JumpCommand) GetAstNodeType() NodeType {
  return NODE_JUMP	
}

type LoopCommand struct {
  body CommandNode 
}

func (self *LoopCommand) GetAstNodeType() NodeType {
  return NODE_LOOP
}

type GlobalCommand struct {
  pattern RegexNode
  body *BlockExpression
}

func (self *GlobalCommand) GetAstNodeType() NodeType {
  return NODE_GLOBAL
}


type ReplaceCommand struct {
  StringToken *Token
}

func (self *ReplaceCommand) GetAstNodeType() NodeType {
  return NODE_REPLACE
}

type ReplaceExprCommand struct {
  expr ExpressionNode
}

func (self *ReplaceExprCommand) GetAstNodeType() NodeType {
  return NODE_REPLACE_EXPR
}

type WriteCommand struct {
  // an optional string; if no string, then this is
  // a write default file.
  expr ExpressionNode
}

func (self *WriteCommand) GetAstNodeType() NodeType {
  return NODE_WRITE
}

type TypeCommand struct {
  expr ExpressionNode
}

func (self *TypeCommand) GetAstNodeType() NodeType {
  return NODE_TYPE
}


type RevertCommand struct {
}

func (self *RevertCommand) GetAstNodeType() NodeType {
  return NODE_REVERT
}

type AssignCommand struct {
  name string
  value ExpressionNode
}

func (self *AssignCommand) GetAstNodeType() NodeType {
  return NODE_ASSIGN
}

type SearchCommand struct {
  regex RegexNode
}

func (self *SearchCommand) GetAstNodeType() NodeType {
  return NODE_SEARCH
}




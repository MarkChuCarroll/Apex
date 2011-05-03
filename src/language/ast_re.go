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

// File: ast_fun.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: ASTs for the Apex programming language regular
//   expressions.
package language

type RegexNode interface {
  AstNode
}

type RegexStringLiteral struct {
   str string
}

func NewRegexStringLiteral(s string) *RegexStringLiteral {
  return &RegexStringLiteral{ s }
}

func (self *RegexStringLiteral) GetAstNodeType() NodeType {
  return NODE_RE_STR
}

type RegexChoice struct {
   choices []RegexNode
}

func NewRegexChoice(choices []RegexNode) *RegexChoice {
  return &RegexChoice{ choices }
}

func (self *RegexChoice) GetAstNodeType() NodeType {
  return NODE_RE_CHOICE
}

type RegexSeq struct {
  elements []RegexNode
}

func NewRegexSeq() *RegexSeq {
  return &RegexSeq{ make([]RegexNode, 0, 10) }
}

func (self *RegexSeq) AddRegexToSeq(n RegexNode) {
  self.elements = append(self.elements, n)
}

func (self *RegexSeq) GetAstNodeType() NodeType {
  return NODE_RE_SEQ
}

type RegexCharset struct {
  chars string
}

func NewRegexCharset(s string) *RegexCharset {
  return &RegexCharset{ s }
}

func (self *RegexCharset) GetAstNodeType() NodeType {
  return NODE_RE_CHARSET
}

type RegexGroup struct {
  regex RegexNode
}

func NewRegexGroup(re RegexNode) *RegexGroup {
  return &RegexGroup{ re }
}

func (self *RegexGroup) GetAstNodeType() NodeType {
  return NODE_RE_GROUP
}

type RegexRepeat struct {
  regex RegexNode
  minRep int32
  strict bool // strict means that only exactly minRep will
         // be accepted.
}

func (self *RegexRepeat) GetAstNodeType() NodeType {
  return NODE_RE_REPEAT
}

type RegexBind struct {
  regex RegexNode
  name string
}

func NewRegexBind(s string, r RegexNode) *RegexBind {
  return &RegexBind{ r, s }
}

func (self *RegexBind) GetAstNodeType() NodeType {
  return NODE_RE_BIND
}

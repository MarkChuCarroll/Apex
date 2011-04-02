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
// Description: ASTs for the Apex programming language function
//     definitions and invocations.

package language

type FunDeclCommand struct {
  pre []string
  post []string
  name string
  body []CommandNode
}

func (self *FunDeclCommand) GetAstNodeType() NodeType {
  return NODE_FUN
}

// Expression
type BlockExpression struct {
  params []string
  body []CommandNode
}

func (self *BlockExpression) GetAstNodeType() NodeType {
  return NODE_BLOCK
}

type InvokeCommand struct {
  pre []ExpressionNode
  post []ExpressionNode 
  fun ExpressionNode
}

func (self *InvokeCommand) GetAstNodeType() NodeType {
  return NODE_INVOKE
}

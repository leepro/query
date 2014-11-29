//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package algebra

import (
	"fmt"

	"github.com/couchbaselabs/query/expression"
	"github.com/couchbaselabs/query/value"
)

type PositionalParameter struct {
	expression.ExpressionBase
	position int
}

func NewPositionalParameter(position int) expression.Expression {
	rv := &PositionalParameter{
		position: position,
	}

	rv.SetExpr(rv)
	return rv
}

func (this *PositionalParameter) Accept(visitor expression.Visitor) (interface{}, error) {
	return visitor.VisitPositionalParameter(this)
}

func (this *PositionalParameter) Type() value.Type { return value.JSON }

func (this *PositionalParameter) Evaluate(item value.Value, context expression.Context) (
	value.Value, error) {
	val, ok := context.(Context).PositionalArg(this.position)

	if ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("No value for positional parameter $%d.", this.position)
	}
}

func (this *PositionalParameter) Indexable() bool {
	return false
}

func (this *PositionalParameter) EquivalentTo(other expression.Expression) bool {
	switch other := other.(type) {
	case *PositionalParameter:
		return this.position == other.position
	default:
		return false
	}
}

func (this *PositionalParameter) SubsetOf(other expression.Expression) bool {
	return this.EquivalentTo(other)
}

func (this *PositionalParameter) Children() expression.Expressions {
	return nil
}

func (this *PositionalParameter) MapChildren(mapper expression.Mapper) error {
	return nil
}

func (this *PositionalParameter) Position() int {
	return this.position
}

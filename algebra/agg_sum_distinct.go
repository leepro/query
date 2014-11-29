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

type SumDistinct struct {
	DistinctAggregateBase
}

func NewSumDistinct(operand expression.Expression) Aggregate {
	rv := &SumDistinct{
		*NewDistinctAggregateBase("sum", operand),
	}

	rv.SetExpr(rv)
	return rv
}

func (this *SumDistinct) String() string {
	return this.toString(this)
}

func (this *SumDistinct) Accept(visitor expression.Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *SumDistinct) Type() value.Type { return value.NUMBER }

func (this *SumDistinct) Evaluate(item value.Value, context expression.Context) (result value.Value, e error) {
	return this.evaluate(this, item, context)
}

func (this *SumDistinct) Constructor() expression.FunctionConstructor {
	return func(operands ...expression.Expression) expression.Function {
		return NewSumDistinct(operands[0])
	}
}

func (this *SumDistinct) Default() value.Value { return value.NULL_VALUE }

func (this *SumDistinct) CumulateInitial(item, cumulative value.Value, context Context) (value.Value, error) {
	item, e := this.Operand().Evaluate(item, context)
	if e != nil {
		return nil, e
	}

	if item.Type() != value.NUMBER {
		return cumulative, nil
	}

	return setAdd(item, cumulative)
}

func (this *SumDistinct) CumulateIntermediate(part, cumulative value.Value, context Context) (value.Value, error) {
	return cumulateSets(part, cumulative)
}

func (this *SumDistinct) ComputeFinal(cumulative value.Value, context Context) (c value.Value, e error) {
	if cumulative == value.NULL_VALUE {
		return cumulative, nil
	}

	av := cumulative.(value.AnnotatedValue)
	set := av.GetAttachment("set").(*value.Set)
	if set.Len() == 0 {
		return value.NULL_VALUE, nil
	}

	sum := 0.0
	for _, v := range set.Values() {
		a := v.Actual()
		switch a := a.(type) {
		case float64:
			sum += a
		default:
			return nil, fmt.Errorf("Invalid partial SUM %v of type %T.", a, a)
		}
	}

	return value.NewValue(sum), nil
}

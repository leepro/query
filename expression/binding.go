//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package expression

import (
	"fmt"

	"github.com/couchbaselabs/query/errors"
	"github.com/couchbaselabs/query/value"
)

type Bindings []*Binding

type Binding struct {
	variable string
	expr     Expression
	descend  bool
}

func NewBinding(variable string, expr Expression) *Binding {
	return &Binding{variable, expr, false}
}

func NewDescendantBinding(variable string, expr Expression) *Binding {
	return &Binding{variable, expr, true}
}

func (this *Binding) Variable() string {
	return this.variable
}

func (this *Binding) Expression() Expression {
	return this.expr
}

func (this *Binding) Descend() bool {
	return this.descend
}

func (this *Binding) Accept(visitor Visitor) (Expression, error) {
	var e error
	this.expr, e = visitor.Visit(this.expr)
	if e != nil {
		return nil, e
	}

	return this.expr, nil
}

func (this Bindings) VisitExpressions(visitor Visitor) (err error) {
	for _, b := range this {
		expr, err := visitor.Visit(b.expr)
		if err != nil {
			return err
		}

		b.expr = expr.(Expression)
	}

	return
}

func (this Bindings) Formalize(allowed value.Value, keyspace string) (a value.Value, err error) {
	a = value.NewScopeValue(make(map[string]interface{}, len(this)), allowed)

	for _, b := range this {
		_, ok := a.Field(b.variable)
		if ok {
			return nil, errors.NewError(nil,
				fmt.Sprintf("Bind alias %s already in scope.", b.variable))
		}

		expr, err := b.expr.Formalize(a, keyspace)
		if err != nil {
			return nil, err
		}

		b.expr = expr
		a.SetField(b.variable, b.variable)
	}

	return a, nil
}

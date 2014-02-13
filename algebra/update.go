//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package algebra

type Update struct {
	bucket    *BucketRef  `json:"bucket"`
	keys      Expression  `json:"keys"`
	set       *Set        `json:"set"`
	unset     *Unset      `json:"unset"`
	where     Expression  `json:"where"`
	limit     Expression  `json:"limit"`
	returning ResultTerms `json:"returning"`
}

func NewUpdate(bucket *BucketRef, keys Expression, set *Set, unset *Unset,
	where, limit Expression, returning ResultTerms) *Update {
	return &Update{bucket, keys, set, unset, where, limit, returning}
}

func (this *Update) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUpdate(this)
}

func (this *Update) BucketRef() *BucketRef {
	return this.bucket
}

func (this *Update) Keys() Expression {
	return this.keys
}

func (this *Update) Set() *Set {
	return this.set
}

func (this *Update) Unset() *Unset {
	return this.unset
}

func (this *Update) Where() Expression {
	return this.where
}

func (this *Update) Limit() Expression {
	return this.limit
}

func (this *Update) Returning() ResultTerms {
	return this.returning
}

type Set struct {
	terms []*SetTerm
}

func NewSet(terms []*SetTerm) *Set {
	return &Set{terms}
}

func (this *Set) Terms() []*SetTerm {
	return this.terms
}

type Unset struct {
	terms []*UnsetTerm
}

func NewUnset(terms []*UnsetTerm) *Unset {
	return &Unset{terms}
}

func (this *Unset) Terms() []*UnsetTerm {
	return this.terms
}

type SetTerm struct {
	path      Path       `json:"path"`
	value     Expression `json:"value"`
	updateFor *UpdateFor `json:"path-for"`
}

func NewSetTerm(path Path, value Expression, updateFor *UpdateFor) *SetTerm {
	return &SetTerm{path, value, updateFor}
}

func (this *SetTerm) Path() Path {
	return this.path
}

func (this *SetTerm) Value() Expression {
	return this.value
}

func (this *SetTerm) UpdateFor() *UpdateFor {
	return this.updateFor
}

type UnsetTerm struct {
	path      Path       `json:"path"`
	updateFor *UpdateFor `json:"path-for"`
}

func NewUnsetTerm(path Path, updateFor *UpdateFor) *UnsetTerm {
	return &UnsetTerm{path, updateFor}
}

func (this *UnsetTerm) Path() Path {
	return this.path
}

func (this *UnsetTerm) UpdateFor() *UpdateFor {
	return this.updateFor
}

type UpdateFor struct {
	bindings []*Binding
	when     Expression
}

func NewUpdateFor(bindings []*Binding, when Expression) *UpdateFor {
	return &UpdateFor{bindings, when}
}

func (this *UpdateFor) Bindings() []*Binding {
	return this.bindings
}

func (this *UpdateFor) When() Expression {
	return this.when
}

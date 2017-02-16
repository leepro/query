//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

import (
	"encoding/json"

	"github.com/couchbase/query/expression"
	"github.com/couchbase/query/expression/parser"
	"github.com/couchbase/query/value"
)

// DistinctScan scans multiple indexes and distincts the results.
type DistinctScan struct {
	readonly
	scan   SecondaryScan
	limit  expression.Expression
	offset expression.Expression
}

func NewDistinctScan(limit, offset expression.Expression, scan SecondaryScan) *DistinctScan {
	return &DistinctScan{
		scan:   scan,
		limit:  limit,
		offset: offset,
	}
}

func (this *DistinctScan) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitDistinctScan(this)
}

func (this *DistinctScan) New() Operator {
	return &DistinctScan{}
}

func (this *DistinctScan) Covers() expression.Covers {
	return this.scan.Covers()
}

func (this *DistinctScan) FilterCovers() map[*expression.Cover]value.Value {
	return this.scan.FilterCovers()
}

func (this *DistinctScan) Covering() bool {
	return this.scan.Covering()
}

func (this *DistinctScan) Scan() SecondaryScan {
	return this.scan
}

func (this *DistinctScan) Limit() expression.Expression {
	return this.limit
}

func (this *DistinctScan) Offset() expression.Expression {
	return this.offset
}

func (this *DistinctScan) SetLimit(limit expression.Expression) {
	this.limit = limit
	this.scan.SetLimit(limit)
}

func (this *DistinctScan) SetOffset(offset expression.Expression) {
	this.offset = offset
	this.scan.SetOffset(offset)
}

func (this *DistinctScan) String() string {
	bytes, _ := this.MarshalJSON()
	return string(bytes)
}

func (this *DistinctScan) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.MarshalBase(nil))
}

func (this *DistinctScan) MarshalBase(f func(map[string]interface{})) map[string]interface{} {
	r := map[string]interface{}{"#operator": "DistinctScan"}
	r["scan"] = this.scan

	if this.limit != nil {
		r["limit"] = expression.NewStringer().Visit(this.limit)
	}

	if this.offset != nil {
		r["offset"] = expression.NewStringer().Visit(this.offset)
	}

	if f != nil {
		f(r)
	}
	return r
}

func (this *DistinctScan) UnmarshalJSON(body []byte) error {
	var _unmarshalled struct {
		_      string          `json:"#operator"`
		Scan   json.RawMessage `json:"scan"`
		Limit  string          `json:"limit"`
		Offset string          `json:"offset"`
	}

	err := json.Unmarshal(body, &_unmarshalled)
	if err != nil {
		return err
	}

	var scan_type struct {
		Operator string `json:"#operator"`
	}

	err = json.Unmarshal(_unmarshalled.Scan, &scan_type)
	if err != nil {
		return err
	}

	scan_op, err := MakeOperator(scan_type.Operator, _unmarshalled.Scan)
	if err != nil {
		return err
	}

	this.scan = scan_op.(SecondaryScan)

	if _unmarshalled.Limit != "" {
		this.limit, err = parser.Parse(_unmarshalled.Limit)
		if err != nil {
			return err
		}
	}

	if _unmarshalled.Offset != "" {
		this.offset, err = parser.Parse(_unmarshalled.Offset)
		if err != nil {
			return err
		}
	}

	return nil
}

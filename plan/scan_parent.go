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
)

// ParentScan is used for UNNEST subqueries.
type ParentScan struct {
	readonly
}

func NewParentScan() *ParentScan {
	return &ParentScan{}
}

func (this *ParentScan) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitParentScan(this)
}

func (this *ParentScan) New() Operator {
	return &ParentScan{}
}

func (this *ParentScan) MarshalJSON() ([]byte, error) {
	r := map[string]interface{}{"#operator": "ParentScan"}
	return json.Marshal(r)
}

func (this *ParentScan) UnmarshalJSON([]byte) error {
	// NOP: ParentScan has no data structure
	return nil
}
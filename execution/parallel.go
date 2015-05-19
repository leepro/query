//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package execution

import (
	"runtime"

	"github.com/couchbase/query/value"
)

type Parallel struct {
	base
	child        Operator
	childChannel StopChannel
}

func NewParallel(child Operator) *Parallel {
	rv := &Parallel{
		base:         newBase(),
		child:        child,
		childChannel: make(StopChannel, runtime.NumCPU()),
	}

	rv.output = rv
	return rv
}

func (this *Parallel) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitParallel(this)
}

func (this *Parallel) Copy() Operator {
	return &Parallel{
		base:         this.base.copy(),
		child:        this.child.Copy(),
		childChannel: make(StopChannel, runtime.NumCPU()),
	}
}

func (this *Parallel) RunOnce(context *Context, parent value.Value) {
	this.once.Do(func() {
		defer context.Recover()       // Recover from any panic
		defer close(this.itemChannel) // Broadcast that I have stopped
		defer this.notify()           // Notify that I have stopped

		n := context.MaxParallelism()

		children := make([]Operator, n)

		// Explicitly make copies, even for the first
		// child. This ensures that the children are
		// identical, as produced by Copy().
		for i := 0; i < n; i++ {
			child := this.child
			if n > 1 {
				child = this.child.Copy()
			}
			child.SetInput(this.input)
			child.SetOutput(this.output)
			child.SetParent(this)
			child.SetStop(nil)
			children[i] = child
		}

		// Run children in parallel
		for i := 0; i < n; i++ {
			go children[i].RunOnce(context, parent)
		}

		for n > 0 {
			select {
			case <-this.childChannel: // Never closed
				// Wait for all children
				n--
			case <-this.stopChannel: // Never closed
				this.notifyStop()
				notifyChildren(children...)
			}
		}
	})
}

func (this *Parallel) ChildChannel() StopChannel {
	return this.childChannel
}

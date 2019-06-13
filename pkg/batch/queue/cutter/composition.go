/*
 *    Copyright 2019 Samsung SDS
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package cutter

import (
	"github.com/nexledger/accelerator/pkg/batch/tx"
)

type Composition func(c *compositeCutter)

type compositeCutter struct {
	cutters []Cutter
}

func (c *compositeCutter) Before(j *tx.Job, next *tx.Item) Cut {
	for _, ct := range c.cutters {
		if ct.Before(j, next) {
			return true
		}
	}
	return false
}

func (c *compositeCutter) After(j *tx.Job) Cut {
	for _, ct := range c.cutters {
		if ct.After(j) {
			return true
		}
	}
	return false
}

func (c *compositeCutter) Clear() {
	for _, ct := range c.cutters {
		ct.Clear()
	}
}

func (c *compositeCutter) Add(cutter Cutter) {
	c.cutters = append(c.cutters, cutter)
}

func WithItemCountCutter(maxItem int) Composition {
	return func(c *compositeCutter) {
		c.Add(&itemCount{maxItem})
	}
}

func WithByteLenCutter(maxByteLen int) Composition {
	return func(c *compositeCutter) {
		c.Add(&byteLen{maxByteLen})
	}
}

func WithMVCCCutter(readKeyIndices []int, writeKeyIndices []int) Composition {
	return func(c *compositeCutter) {
		c.Add(&mvcc{
			readKeyIndices:  readKeyIndices,
			writeKeyIndices: writeKeyIndices,
			writtenKeys:     make(map[string]bool),
		})
	}
}

func New(compositions ...Composition) Cutter {
	ct := &compositeCutter{}
	for _, composition := range compositions {
		composition(ct)
	}
	return ct
}

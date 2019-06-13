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

type mvcc struct {
	readKeyIndices  []int
	writeKeyIndices []int
	writtenKeys     map[string]bool
}

func (c *mvcc) Before(_ *tx.Job, item *tx.Item) Cut {
	for _, idx := range c.readKeyIndices {
		rKey := string(item.Args[idx])
		if _, ok := c.writtenKeys[rKey]; ok {
			return true
		}
	}

	return false
}

func (c *mvcc) After(job *tx.Job) Cut {
	if item, ok := job.LastItem(); ok {
		for _, idx := range c.writeKeyIndices {
			wKey := string(item.Args[idx])
			c.writtenKeys[wKey] = true
		}
	}

	return false
}

func (c *mvcc) Clear() {
	c.writtenKeys = make(map[string]bool)
}

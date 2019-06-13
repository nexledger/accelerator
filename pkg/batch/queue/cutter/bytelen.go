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

type byteLen struct {
	maxByteLen int
}

func (c *byteLen) Before(job *tx.Job, item *tx.Item) Cut {
	return job.Size() > 0 && job.ByteLen()+itemByteLen(item) >= c.maxByteLen
}

func (c *byteLen) After(job *tx.Job) Cut {
	//In case of one itemByteLen exceeds the maxByteLen, it is sent to an isolated job
	return job.ByteLen() >= c.maxByteLen
}

func (c *byteLen) Clear() {}

func itemByteLen(i *tx.Item) int {
	var size int
	for _, row := range i.Args {
		size += len(row)
	}
	return size
}

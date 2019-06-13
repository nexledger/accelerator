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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nexledger/accelerator/pkg/batch/tx"
)

const (
	maxByteLen  = 1024
	maxItemSize = 5
)

func TestByteLen(t *testing.T) {
	cutter := byteLen{maxByteLen}

	job := &tx.Job{}
	item := &tx.Item{Args: [][]byte{make([]byte, 1000)}}

	assert.False(t, bool(cutter.Before(job, item)))
	job.Add(item)
	assert.False(t, bool(cutter.After(job)))

	assert.True(t, bool(cutter.Before(job, item)), "Expected to be cut because job length exceeds maxByteLen")
	cutter.Clear()
	job = &tx.Job{}
	assert.False(t, bool(cutter.After(job)))

	overflowItem := &tx.Item{Args: [][]byte{make([]byte, 2000)}}

	assert.False(t, bool(cutter.Before(job, item)))
	job.Add(overflowItem)
	assert.True(t, bool(cutter.Before(job, item)), "Expected to be cut because a Item byte size exceeds maxByteLen of job")
}

func TestItemCount(t *testing.T) {
	cutter := itemCount{maxItemSize}

	job := &tx.Job{}
	item := &tx.Item{Args: [][]byte{}}

	for i := 1; i <= maxItemSize; i++ {
		assert.False(t, bool(cutter.Before(job, item)))
		job.Add(item)
		if cutter.After(job) {
			assert.Equal(t, i, maxItemSize, "Expected to be cut when item size reaches maxItemSize")
		}
	}
}

func TestMVCC(t *testing.T) {
	cutter := mvcc{readKeyIndices: []int{0, 2}, writeKeyIndices: []int{1, 3}, writtenKeys: make(map[string]bool)}
	job := &tx.Job{}
	item := &tx.Item{Args: [][]byte{[]byte("A"), []byte("B"), []byte("A"), []byte("B")}}
	job.Add(item)

	assert.False(t, bool(cutter.Before(job, item)))
	job.Add(item)
	assert.False(t, bool(cutter.After(job)))

	conflictItem := &tx.Item{Args: [][]byte{[]byte("B"), []byte("A"), []byte("C"), []byte("A")}}
	assert.True(t, bool(cutter.Before(job, conflictItem)), "Expected to be cut because read B after write B")
	job = &tx.Job{}
	cutter.Clear()
	job.Add(conflictItem)
	assert.False(t, bool(cutter.After(job)))
}

func TestComposition(t *testing.T) {
	cutterOpts := make([]Composition, 0)
	cutterOpts = append(cutterOpts, WithByteLenCutter(maxByteLen))
	cutterOpts = append(cutterOpts, WithItemCountCutter(maxItemSize))
	cutterOpts = append(cutterOpts, WithMVCCCutter([]int{0}, []int{2}))
	cutters := New(cutterOpts...)

	job := &tx.Job{}
	item := &tx.Item{Args: [][]byte{[]byte("A"), make([]byte, 500), []byte("B"), make([]byte, 500)}}
	job.Add(item)

	assert.True(t, bool(cutters.Before(job, item)), "Expected to be cut because byteLen cutter works")
	job = &tx.Job{}
	cutters.Clear()
	job.Add(item)
	assert.False(t, bool(cutters.After(job)))

	item = &tx.Item{Args: [][]byte{[]byte("B"), make([]byte, 100), []byte("A"), make([]byte, 100)}}
	assert.True(t, bool(cutters.Before(job, item)), "Expected to be cut because mvcc cutter works")
	job = &tx.Job{}
	cutters.Clear()
	job.Add(item)
	assert.False(t, bool(cutters.After(job)))

	for i := 2; i <= maxItemSize; i++ {
		assert.False(t, bool(cutters.Before(job, item)))
		job.Add(item)
		if cutters.After(job) {
			assert.Equal(t, i, maxItemSize, "Expected to be cut because itemCount cutter works")
		}
	}
}

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

package tx

type Resolver func(items []*Item)

type Item struct {
	Args     [][]byte
	Notifier chan *Result
}

type Result struct {
	TxId            string
	ValidationCode  int32
	ChaincodeStatus int32
	Payload         []byte
	Error           error
}

type Job struct {
	args      [][][]byte
	byteLen   int
	items     []*Item
	notifiers []chan *Result
	Retry     bool
}

func (j *Job) Add(i *Item) *Job {
	j.items = append(j.items, i)

	if i.Args == nil {
		j.args = append(j.args, make([][]byte, 0, 0))
	} else {
		j.args = append(j.args, i.Args)
	}

	for _, row := range i.Args {
		j.byteLen = j.byteLen + len(row)
	}

	j.notifiers = append(j.notifiers, i.Notifier)

	return j
}

func (j *Job) Size() int {
	return len(j.items)
}

func (j *Job) ByteLen() int {
	return j.byteLen
}

func (j *Job) Items() []*Item {
	return j.items
}

func (j *Job) LastItem() (item *Item, exist bool) {
	if j.Size() == 0 {
		return nil, false
	}
	return j.items[j.Size()-1], true
}

func (j *Job) Notifiers() []chan *Result {
	return j.notifiers
}

func (j *Job) Args() [][][]byte {
	return j.args
}

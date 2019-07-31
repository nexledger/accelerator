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

package queue

import (
	"github.com/nexledger/accelerator/pkg/batch/queue/cutter"
	"github.com/nexledger/accelerator/pkg/batch/route"
	"github.com/nexledger/accelerator/pkg/batch/tx"
)

type Processor interface {
	Submit(i *tx.Item) bool
	Process()
	Empty() bool
}

type processor struct {
	sender route.Sender
	cutter cutter.Cutter
	job    *tx.Job
}

func (p *processor) Submit(i *tx.Item) bool {
	processed := false

	if p.cutter.Before(p.job, i) {
		processed = true
		p.Process()
	}

	p.job.Add(i)

	if p.cutter.After(p.job) {
		processed = true
		p.Process()
	}

	return processed
}

func (p *processor) Process() {
	go p.sender.Send(p.job)
	p.job = &tx.Job{}
	p.cutter.Clear()
}

func (p *processor) Empty() bool {
	return p.job.Size() == 0
}

func NewProcessor(sender route.Sender, compositions []cutter.Composition) Processor {
	return &processor{sender, cutter.New(compositions...), &tx.Job{}}
}

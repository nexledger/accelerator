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
	"time"

	"github.com/nexledger/accelerator/pkg/batch/queue/cutter"
	"github.com/nexledger/accelerator/pkg/batch/route"
	"github.com/nexledger/accelerator/pkg/batch/tx"
)

type Scheduler struct {
	maxWaitTime    time.Duration
	timer          *time.Timer
	scheduledItems chan *tx.Item
	processor      *processor
}

func (s *Scheduler) Schedule(i *tx.Item) {
	s.scheduledItems <- i
}

func (s *Scheduler) Start() {
	s.timer = time.NewTimer(s.maxWaitTime)
	s.timer.Stop()
	go func() {
		for {
			select {
			case i := <-s.scheduledItems:
				s.scheduled(i)
			case <-s.timer.C:
				s.timeout()
			}
		}
	}()
}

func (s *Scheduler) Close() {
	// Do nothing
}

func (s *Scheduler) scheduled(i *tx.Item) {
	if s.processor.Empty() {
		s.timer.Reset(s.maxWaitTime)
	}

	if s.processor.Submit(i) {
		if !s.processor.Empty() {
			s.timer.Reset(s.maxWaitTime)
		} else {
			s.timer.Stop()
		}
	}
}

func (s *Scheduler) timeout() {
	if !s.processor.Empty() {
		s.processor.Process()
		s.timer.Stop()
	}
}

func New(sender *route.Sender, compositions []cutter.Composition, maxWaitTime time.Duration, queueSize int) *Scheduler {
	processor := &processor{sender, cutter.New(compositions...), &tx.Job{}}
	scheduler := &Scheduler{
		maxWaitTime,
		nil,
		make(chan *tx.Item, queueSize),
		processor,
	}
	return scheduler
}

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

package route

import (
	"strings"
	"sync"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"github.com/nexledger/accelerator/pkg/batch/route/encoding"
	"github.com/nexledger/accelerator/pkg/batch/route/fab"
	"github.com/nexledger/accelerator/pkg/batch/tx"
)

type Sender struct {
	invoker   fab.Invoker
	responder *responder
	recovery  bool
}

func (s *Sender) Send(job *tx.Job) {
	fabresp, err := s.invoker(job)
	if err != nil {
		if s.recovery && strings.Contains(err.Error(), "MVCC_READ_CONFLICT") {
			s.retry(job, fabresp)
			return
		}
		s.responder.JobFailure(job, err)
		return
	}

	s.responder.JobSuccess(job, fabresp)
}

func (s *Sender) retry(job *tx.Job, resp *channel.Response) {
	var wg sync.WaitGroup
	wg.Add(len(job.Items()))
	for _, i := range job.Items() {
		go func(i *tx.Item, resp *channel.Response, wg *sync.WaitGroup) {
			defer wg.Done()

			job := &tx.Job{Retry: true}
			fabresp, err := s.invoker(job.Add(i))
			if err != nil {
				s.responder.JobFailure(job, err)
				return
			}
			s.responder.JobSuccess(job, fabresp)
		}(i, resp, &wg)
	}
	wg.Wait()
}

func New(invoker fab.Invoker, encoder encoding.Encoder, recovery bool) (*Sender, error) {
	return &Sender{
		invoker,
		&responder{encoder},
		recovery,
	}, nil
}

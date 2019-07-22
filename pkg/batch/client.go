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

package batch

import (
	"time"

	"github.com/pkg/errors"

	"github.com/nexledger/accelerator/pkg/batch/queue"
	"github.com/nexledger/accelerator/pkg/batch/queue/cutter"
	"github.com/nexledger/accelerator/pkg/batch/route"
	"github.com/nexledger/accelerator/pkg/batch/route/encoding"
	"github.com/nexledger/accelerator/pkg/batch/route/fab"
	"github.com/nexledger/accelerator/pkg/batch/tx"
	"github.com/nexledger/accelerator/pkg/core"
)

type Acceleration struct {
	Type               string
	ChannelId          string
	ChaincodeName      string
	Fcn                string
	QueueSize          int
	MaxBatchItems      int
	MaxWaitTimeSeconds int64
	MaxBatchBytes      int
	ReadKeyIndices     []int
	WriteKeyIndices    []int
	Encoding           string
	Recovery           bool
}

type Client struct {
	ctx               *core.Context
	executeSchedulers map[string]*queue.Scheduler
	querySchedulers   map[string]*queue.Scheduler
}

func (s *Client) Execute(channelId, chaincodeName, fcn string, args [][]byte) (*tx.Result, error) {
	name := nameOf(channelId, chaincodeName, fcn)
	if scheduler, ok := s.executeSchedulers[name]; !ok {
		return nil, errors.New("Execute Scheduler not found: " + name)
	} else {
		return process(scheduler, args)
	}
}

func (s *Client) Query(channelId, chaincodeName, fcn string, args [][]byte) (*tx.Result, error) {
	name := nameOf(channelId, chaincodeName, fcn)
	if scheduler, ok := s.querySchedulers[name]; !ok {
		return nil, errors.New("Query Scheduler not found: " + name)
	} else {
		return process(scheduler, args)
	}
}

func (s *Client) Register(acc *Acceleration) error {
	var schedulers map[string]*queue.Scheduler
	if acc.Type == "execute" {
		schedulers = s.executeSchedulers
	} else if acc.Type == "query" {
		schedulers = s.querySchedulers
	} else {
		return errors.New("Unsupported type: " + acc.Type)
	}

	name := nameOf(acc.ChannelId, acc.ChaincodeName, acc.Fcn)
	if _, ok := schedulers[name]; ok {
		return errors.New("Scheduler already registered: " + name)
	}

	encoder, err := encoding.New(acc.Encoding)
	if err != nil {
		return err
	}

	invoker, err := fab.New(s.ctx, acc.ChannelId, acc.ChaincodeName, acc.Fcn, acc.Type, encoder)
	if err != nil {
		return err
	}

	sender, err := route.New(invoker, encoder, acc.Recovery)
	if err != nil {
		return err
	}

	cutterOpts := make([]cutter.Composition, 0)
	if acc.MaxBatchItems > 0 {
		cutterOpts = append(cutterOpts, cutter.WithItemCountCutter(acc.MaxBatchItems))
	}
	if acc.MaxBatchBytes > 0 {
		cutterOpts = append(cutterOpts, cutter.WithByteLenCutter(acc.MaxBatchBytes))
	}
	if len(acc.ReadKeyIndices) > 0 {
		cutterOpts = append(cutterOpts, cutter.WithMVCCCutter(acc.ReadKeyIndices, acc.WriteKeyIndices))
	}

	scheduler := queue.New(
		sender,
		cutterOpts,
		time.Duration(acc.MaxWaitTimeSeconds)*time.Second,
		acc.QueueSize,
	)
	scheduler.Start()
	schedulers[name] = scheduler

	return nil
}

func nameOf(channelId, chaincodeName, fcn string) string {
	return channelId + ":" + chaincodeName + ":" + fcn
}

func process(s *queue.Scheduler, args [][]byte) (*tx.Result, error) {
	notify := make(chan *tx.Result)
	s.Schedule(&tx.Item{Args: args, Notifier: notify})
	result := <-notify
	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

func New(ctx *core.Context) *Client {
	return &Client{
		ctx,
		make(map[string]*queue.Scheduler),
		make(map[string]*queue.Scheduler),
	}
}

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

package ping

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"google.golang.org/grpc"

	pbbatch "github.com/nexledger/accelerator/protos"
)

const (
	channelId     = "accelerator"
	chaincodeName = "ping"
	numOfPings    = 50
	address       = "127.0.0.1:5050"
)

func TestBatch(t *testing.T) {
	TestPing(t)
	TestPong(t)
}

func TestPing(t *testing.T) {
	client := pbbatch.NewAcceleratorServiceClient(connect(t))
	notifiers := make([]chan string, numOfPings)
	for i := 0; i < numOfPings; i++ {
		notifier := make(chan string)
		notifiers[i] = notifier
		go func(i int, notifier chan string) {
			req := &pbbatch.TxRequest{
				ChannelId:     channelId,
				ChaincodeName: chaincodeName,
				Fcn:           "ping",
				Args:          [][]byte{[]byte(strconv.Itoa(i)), []byte("value of " + strconv.Itoa(i))},
			}
			resp, err := client.Execute(context.Background(), req)
			if err != nil {
				notifier <- "Failed to execute" + err.Error()
			} else {
				notifier <- strconv.Itoa(i) + ":" + resp.TxId
			}
		}(i, notifier)
	}

	for i := 0; i < numOfPings; i++ {
		fmt.Println(<-notifiers[i])
	}
}

func TestPong(t *testing.T) {
	client := pbbatch.NewAcceleratorServiceClient(connect(t))
	notifiers := make([]chan string, numOfPings)
	for i := 0; i < numOfPings; i++ {
		notifier := make(chan string)
		notifiers[i] = notifier
		go func(i int, notifier chan string) {
			req := &pbbatch.TxRequest{
				ChannelId:     channelId,
				ChaincodeName: chaincodeName,
				Fcn:           "pong",
				Args:          [][]byte{[]byte(strconv.Itoa(i))},
			}
			resp, err := client.Query(context.Background(), req)
			if err != nil {
				notifier <- "Failed to query" + err.Error()
			} else {
				notifier <- strconv.Itoa(i) + ":" + string(resp.Payload)
			}
		}(i, notifier)
	}

	for i := 0; i < numOfPings; i++ {
		fmt.Println(<-notifiers[i])
	}
}

func connect(t *testing.T) *grpc.ClientConn {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect server.", err)
	}
	return cc
}

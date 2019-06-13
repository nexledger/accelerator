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

package fab

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/pkg/errors"

	"github.com/nexledger/accelerator/pkg/core"
)

type Invoker func([][]byte, int, ...channel.RequestOption) (*channel.Response, error)

func New(ctx *core.Context, channelId string, ccId string, fcn string, typ string) (Invoker, error) {
	switch typ {
	case "execute":
		return func(args [][]byte, txCnt int, opts ...channel.RequestOption) (*channel.Response, error) {
			client, err := ctx.ChannelClient(channelId)
			if err != nil {
				return nil, err
			}

			resp, err := client.Execute(
				channel.Request{ChaincodeID: ccId, Fcn: fcn, Args: args},
				opts...,
			)
			return &resp, err
		}, nil
	case "query":
		return func(args [][]byte, txCnt int, opts ...channel.RequestOption) (*channel.Response, error) {
			client, err := ctx.ChannelClient(channelId)
			if err != nil {
				return nil, err
			}

			resp, err := client.Query(
				channel.Request{ChaincodeID: ccId, Fcn: fcn, Args: args},
				opts...,
			)
			return &resp, err
		}, nil
	}
	return nil, errors.New("Unsupported type: " + typ)
}

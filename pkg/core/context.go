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

package core

import (
	"sync"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

type Context struct {
	sdk *fabsdk.FabricSDK

	user string
	org  string

	channelClientsLock sync.Mutex
	channelClients     sync.Map
}

func (ctx *Context) ChannelClient(channelId string) (*channel.Client, error) {
	if client, ok := ctx.channelClients.Load(channelId); ok {
		return client.(*channel.Client), nil
	}

	ctx.channelClientsLock.Lock()
	defer ctx.channelClientsLock.Unlock()
	if client, ok := ctx.channelClients.Load(channelId); ok {
		return client.(*channel.Client), nil
	}

	provider := ctx.sdk.ChannelContext(channelId, fabsdk.WithUser(ctx.user), fabsdk.WithOrg(ctx.org))
	client, err := channel.New(provider)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create channel context")
	}
	ctx.channelClients.Store(channelId, client)
	return client, nil
}

func (ctx *Context) FabricSDK() *fabsdk.FabricSDK {
	return ctx.sdk
}

func (ctx *Context) User() string {
	return ctx.user
}

func (ctx *Context) Organization() string {
	return ctx.org
}

func (ctx *Context) Close() {
	ctx.sdk.Close()
}

func New(confFilePath string, user string, org string, opts ...fabsdk.Option) (*Context, error) {
	if sdk, err := fabsdk.New(config.FromFile(confFilePath), opts...); err != nil {
		return nil, err
	} else {
		return &Context{sdk: sdk, user: user, org: org}, nil
	}
}

func Wrap(sdk *fabsdk.FabricSDK, user string, org string) *Context {
	return &Context{sdk: sdk, user: user, org: org}
}

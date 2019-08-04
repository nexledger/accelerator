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

package mocks

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	fabctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	fcmocks "github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	mspmocks "github.com/hyperledger/fabric-sdk-go/pkg/msp/test/mockmsp"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"

	"github.com/nexledger/accelerator/pkg/fabwrap"
	"github.com/nexledger/accelerator/pkg/fabwrap/network"
)

const channelID = "testChannel"

type mockContext struct{}

func (mockContext) ResourceClient() (*resmgmt.Client, error) {
	return resmgmt.New(mockClientProvider())
}

func (mockContext) ChannelClient(channelId string) (*channel.Client, error) {
	return channel.New(mockChannelProvider(channelID))
}

func (mockContext) NetworkClient() (*network.Client, error) {
	config := &fab.NetworkConfig{
		Channels:      make(map[string]fab.ChannelEndpointConfig, 0),
		Organizations: make(map[string]fab.OrganizationConfig, 0),
		Orderers:      make(map[string]fab.OrdererConfig, 0),
		Peers:         make(map[string]fab.PeerConfig, 0),
	}

	return &network.Client{Org: "", Config: config}, nil
}

func mockChannelProvider(channelID string) fabctx.ChannelProvider {
	channelProvider := func() (fabctx.Channel, error) {
		return mocks.NewMockChannel(channelID)
	}
	return channelProvider
}

func mockClientProvider() fabctx.ClientProvider {
	ctx := mocks.NewMockContext(mspmocks.NewMockSigningIdentity("test", "Org1MSP"))

	// Create mock orderer with simple mock block
	orderer := mocks.NewMockOrderer("", nil)
	orderer.EnqueueForSendDeliver(mocks.NewSimpleMockBlock())
	orderer.EnqueueForSendDeliver(common.Status_SUCCESS)
	orderer.CloseQueue()

	setupCustomOrderer(ctx, orderer)

	clientProvider := func() (fabctx.Client, error) {
		return ctx, nil
	}

	return clientProvider
}

func setupCustomOrderer(ctx *fcmocks.MockContext, mockOrderer fab.Orderer) *fcmocks.MockContext {
	mockInfraProvider := &fcmocks.MockInfraProvider{}
	mockInfraProvider.SetCustomOrderer(mockOrderer)
	ctx.SetCustomInfraProvider(mockInfraProvider)
	return ctx
}

func NewMockContext() (fabwrap.Context, error) {
	return &mockContext{}, nil
}

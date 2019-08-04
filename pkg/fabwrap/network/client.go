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

package network

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type Client struct {
	Org    string
	Config *fab.NetworkConfig
}

func (c *Client) GetPeerUrls() []string {
	peerUrls := make([]string, 0)
	for _, peerConfig := range c.Config.Peers {
		peerUrls = append(peerUrls, peerConfig.URL)
	}
	return peerUrls
}

func (c *Client) GetOrdererUrls() []string {
	ordererUrls := make([]string, 0)
	for _, ordererConfig := range c.Config.Orderers {
		ordererUrls = append(ordererUrls, ordererConfig.URL)
	}
	return ordererUrls
}

func (c *Client) GetMspId() string {
	return c.Config.Organizations[c.Org].MSPID
}

func NewNetworkClient(org string, networkConfig *fab.NetworkConfig) *Client {
	return &Client{org, networkConfig}
}

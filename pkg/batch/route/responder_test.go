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
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/stretchr/testify/assert"

	"github.com/nexledger/accelerator/pkg/batch/mocks"
	"github.com/nexledger/accelerator/pkg/batch/tx"
)

const itemSize = 5

func TestResponder_JobSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockEncoder := mocks.NewMockEncoder(mockCtrl)

	msg := [][]byte{[]byte("1"), []byte("2"), []byte("3"), []byte("4"), []byte("5")}
	mockEncoder.EXPECT().EncodeRequest(gomock.Any()).Return(msg, nil).AnyTimes()
	mockEncoder.EXPECT().DecodeResponse(gomock.Any()).Return(msg, nil).AnyTimes()

	job := &tx.Job{}
	notify := make(chan *tx.Result, itemSize)
	for i := 0; i < itemSize; i++ {
		item := &tx.Item{Args: [][]byte{}, Notifier: notify}
		job.Add(item)
	}

	responder := responder{mockEncoder}
	responder.JobSuccess(job, &channel.Response{})

	for i := 0; i < itemSize; i++ {
		result := <-notify
		assert.NoError(t, result.Error)
	}
}

func TestResponder_JobFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockEncoder := mocks.NewMockEncoder(mockCtrl)

	msg := [][]byte{[]byte("1"), []byte("2"), []byte("3"), []byte("4"), []byte("5")}
	mockEncoder.EXPECT().EncodeRequest(gomock.Any()).Return(msg, nil).AnyTimes()
	//return empty struct
	mockEncoder.EXPECT().DecodeResponse(gomock.Any()).Return([][]byte{}, nil).AnyTimes()

	job := &tx.Job{}
	notify := make(chan *tx.Result, itemSize)
	for i := 0; i < itemSize; i++ {
		item := &tx.Item{Args: [][]byte{}, Notifier: notify}
		job.Add(item)
	}

	responder := responder{mockEncoder}
	responder.JobFailure(job, errors.New("failure"))

	for i := 0; i < itemSize; i++ {
		result := <-notify
		assert.Error(t, result.Error)
	}

	responder.JobSuccess(job, &channel.Response{})

	for i := 0; i < itemSize; i++ {
		result := <-notify
		assert.Error(t, result.Error)
	}
}

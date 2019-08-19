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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/nexledger/accelerator/pkg/batch/mocks"
	"github.com/nexledger/accelerator/pkg/batch/tx"
)

func TestSenderOnSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockResponder := mocks.NewMockResponder(mockCtrl)
	mockResponder.EXPECT().JobSuccess(gomock.Any(), gomock.Any()).Return().Times(1)

	fakeInvoker := func(job *tx.Job, opts ...channel.RequestOption) (*channel.Response, error) {
		return nil, nil
	}

	sender, err := NewSender(fakeInvoker, mockResponder, false)
	assert.NoError(t, err)
	sender.Send(&tx.Job{})
}

func TestSenderOnFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockResponder := mocks.NewMockResponder(mockCtrl)
	mockResponder.EXPECT().JobFailure(gomock.Any(), gomock.Any()).Return().Times(1)

	errInvoker := func(job *tx.Job, opts ...channel.RequestOption) (*channel.Response, error) {
		return nil, errors.New("Unknown Error")
	}

	job := &tx.Job{}
	job.Add(&tx.Item{})

	sender, err := NewSender(errInvoker, mockResponder, true)
	assert.NoError(t, err)
	sender.Send(job)
}

func TestSenderOnRetry(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockResponder := mocks.NewMockResponder(mockCtrl)
	mockResponder.EXPECT().JobSuccess(gomock.Any(), gomock.Any()).Return().Times(1)

	mvccInvoker := func(job *tx.Job, opts ...channel.RequestOption) (*channel.Response, error) {
		if job.Retry {
			return nil, nil
		}
		return nil, errors.New("MVCC_READ_CONFLICT")
	}

	job := &tx.Job{}
	job.Add(&tx.Item{})

	sender, err := NewSender(mvccInvoker, mockResponder, true)
	assert.NoError(t, err)
	sender.Send(job)
}

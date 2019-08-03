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
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/nexledger/accelerator/pkg/batch/mocks"
	"github.com/nexledger/accelerator/pkg/batch/tx"
)

func TestSchedulerAfterTimout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProcessor := mocks.NewMockProcessor(mockCtrl)
	mockProcessor.EXPECT().Process().Return().Times(1)
	mockProcessor.EXPECT().Submit(gomock.Any()).Return(false)
	mockProcessor.EXPECT().Empty().Return(true).Times(1)

	scheduler := NewScheduler(mockProcessor, time.Millisecond, 1000)
	scheduler.Start()
	scheduler.Schedule(&tx.Item{})

	mockProcessor.EXPECT().Submit(gomock.Any()).Return(true)
	mockProcessor.EXPECT().Empty().Return(false).AnyTimes()

	scheduler.Schedule(&tx.Item{})

	mockProcessor.EXPECT().Submit(gomock.Any()).Return(false)
	mockProcessor.EXPECT().Empty().Return(true).AnyTimes()

	scheduler.Schedule(&tx.Item{})

	time.Sleep(1 * time.Second)
}

func TestSchedulerBeforeTimout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProcessor := mocks.NewMockProcessor(mockCtrl)
	mockProcessor.EXPECT().Process().Return().Times(0)
	mockProcessor.EXPECT().Submit(gomock.Any()).Return(false)
	mockProcessor.EXPECT().Empty().Return(true).Times(1)

	scheduler := NewScheduler(mockProcessor, 5*time.Second, 1000)
	scheduler.Start()
	scheduler.Schedule(&tx.Item{})

	mockProcessor.EXPECT().Submit(gomock.Any()).Return(true)
	mockProcessor.EXPECT().Empty().Return(false).AnyTimes()

	scheduler.Schedule(&tx.Item{})

	mockProcessor.EXPECT().Submit(gomock.Any()).Return(false)
	mockProcessor.EXPECT().Empty().Return(true).AnyTimes()

	scheduler.Schedule(&tx.Item{})

	time.Sleep(1 * time.Second)
}

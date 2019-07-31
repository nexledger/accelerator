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

package queue

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/nexledger/accelerator/pkg/batch/mocks"
	"github.com/nexledger/accelerator/pkg/batch/queue/cutter"
	"github.com/nexledger/accelerator/pkg/batch/tx"
)

func TestProcessor(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSender := mocks.NewMockSender(mockCtrl)
	mockSender.EXPECT().Send(gomock.Any()).Return().AnyTimes()

	mockCutter := mocks.NewMockCutter(mockCtrl)
	processor := processor{mockSender, mockCutter, &tx.Job{}}

	mockCutter.EXPECT().Before(gomock.Any(), gomock.Any()).Return(cutter.Cut(false))
	mockCutter.EXPECT().After(gomock.Any()).Return(cutter.Cut(false))
	mockCutter.EXPECT().Clear().Return().AnyTimes()

	processor.Submit(&tx.Item{})
	if processor.job.Size() != 1 {
		t.Fatalf("Should have job size 1")
	}

	mockCutter.EXPECT().Before(gomock.Any(), gomock.Any()).Return(cutter.Cut(false))
	mockCutter.EXPECT().After(gomock.Any()).Return(cutter.Cut(true))

	res := processor.Submit(&tx.Item{})
	if !res {
		t.Fatalf("Should have submmited")
	}
	if !processor.Empty() {
		t.Fatalf("Should have job size 0")
	}

	mockCutter.EXPECT().Before(gomock.Any(), gomock.Any()).Return(cutter.Cut(true))
	mockCutter.EXPECT().After(gomock.Any()).Return(cutter.Cut(false))

	res = processor.Submit(&tx.Item{})
	if !res {
		t.Fatalf("Should have submmited")
	}
	if processor.Empty() {
		t.Fatalf("Should have job size 1")
	}
}

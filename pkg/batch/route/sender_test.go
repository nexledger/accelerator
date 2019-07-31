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

func TestSenderOnRetry(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockResponder := mocks.NewMockResponder(mockCtrl)
	mockResponder.EXPECT().JobFailure(gomock.Any(), gomock.Any()).Return().Times(1)

	mvccInvoker := func(job *tx.Job, opts ...channel.RequestOption) (*channel.Response, error) {
		return nil, errors.New("MVCC_READ_CONFLICT")
	}

	job := &tx.Job{}
	job.Add(&tx.Item{})

	sender, err := NewSender(mvccInvoker, mockResponder, true)
	assert.NoError(t, err)
	sender.Send(job)
}

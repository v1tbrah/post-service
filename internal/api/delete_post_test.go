package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/post-service/internal/api/mocks"
	"gitlab.com/pet-pr-social-network/post-service/internal/msgsndr"
	"gitlab.com/pet-pr-social-network/post-service/internal/storage"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
)

func TestAPI_DeletePost(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name            string
		mockStorage     func(t *testing.T) *mocks.Storage
		req             *ppbapi.DeletePostRequest
		expectedResp    *ppbapi.Empty
		wantErr         bool
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name: "OK",
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("DeletePost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1)).
					Return(int64(1), nil).
					Once()
				return testStorage
			},
			req:          &ppbapi.DeletePostRequest{Id: int64(1)},
			expectedResp: &ppbapi.Empty{},
		},
		{
			name: "not found",
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("DeletePost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1)).
					Return(int64(-1), storage.ErrPostNotFoundByID).
					Once()
				return testStorage
			},
			req:             &ppbapi.DeletePostRequest{Id: int64(1)},
			wantErr:         true,
			expectedErr:     ppbapi.ErrPostNotFoundByID,
			expectedErrCode: codes.NotFound,
		},
		{
			name: "unexpected err on storage.DeletePost",
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("DeletePost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1)).
					Return(int64(-1), errors.New("unexpected error")).
					Once()
				return testStorage
			},
			req:             &ppbapi.DeletePostRequest{Id: int64(1)},
			wantErr:         true,
			expectedErr:     errors.New("unexpected error"),
			expectedErrCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msgSender *msgsndr.Sender
			a := &API{storage: tt.mockStorage(t), msgSender: msgSender}
			_, err := a.DeletePost(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrCode, status.Code(err))
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			}

			if !tt.wantErr {
				require.NoError(t, err)
			}
		})
	}
}

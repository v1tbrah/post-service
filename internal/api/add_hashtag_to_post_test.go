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
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
)

func TestAPI_AddHashtagToPost(t *testing.T) {
	tests := []struct {
		name            string
		mockStorage     func(t *testing.T) *mocks.Storage
		req             *ppbapi.AddHashtagToPostRequest
		wantErr         bool
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name: "OK",
			req:  &ppbapi.AddHashtagToPostRequest{PostID: 1, HashtagID: 1},
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("AddHashtagToPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1), int64(1)).
					Return(nil).
					Once()
				return testStorage
			},
		},
		{
			name: "unexpected err on storage.AddHashtagToPost",
			req:  &ppbapi.AddHashtagToPostRequest{PostID: 1, HashtagID: 1},
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("AddHashtagToPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1), int64(1)).
					Return(errors.New("unexpected err")).
					Once()
				return testStorage
			},
			wantErr:         true,
			expectedErr:     errors.New("unexpected"),
			expectedErrCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{storage: tt.mockStorage(t)}
			_, err := a.AddHashtagToPost(context.Background(), tt.req)

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

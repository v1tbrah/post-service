package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/post-service/internal/api/mocks"
)

func TestAPI_CreateHashtag(t *testing.T) {
	tests := []struct {
		name            string
		mockStorage     func(t *testing.T) *mocks.Storage
		req             *ppbapi.CreateHashtagRequest
		expectedResp    *ppbapi.CreateHashtagResponse
		wantErr         bool
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name: "OK",
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("CreateHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("model.Hashtag")).
					Return(int64(1), nil).
					Once()
				return testStorage
			},
			req: &ppbapi.CreateHashtagRequest{
				Name: "TestName",
			},
			expectedResp: &ppbapi.CreateHashtagResponse{Id: int64(1)},
		},
		{
			name: "empty name",
			mockStorage: func(t *testing.T) *mocks.Storage {
				return mocks.NewStorage(t)
			},
			req: &ppbapi.CreateHashtagRequest{
				Name: "",
			},
			wantErr:         true,
			expectedErr:     ppbapi.ErrEmptyName,
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "empty name with spaces",
			mockStorage: func(t *testing.T) *mocks.Storage {
				return mocks.NewStorage(t)
			},
			req: &ppbapi.CreateHashtagRequest{
				Name: "   ",
			},
			wantErr:         true,
			expectedErr:     ppbapi.ErrEmptyName,
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "unexpected err on storage.CreateHashtag",
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("CreateHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("model.Hashtag")).
					Return(int64(0), errors.New("unexpected err")).
					Once()
				return testStorage
			},
			req: &ppbapi.CreateHashtagRequest{
				Name: "TestName",
			},
			wantErr:         true,
			expectedErr:     errors.New("unexpected"),
			expectedErrCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{storage: tt.mockStorage(t)}
			resp, err := a.CreateHashtag(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrCode, status.Code(err))
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			}

			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}

package api

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.com/pet-pr-social-network/post-service/internal/api/mocks"
	"gitlab.com/pet-pr-social-network/post-service/internal/model"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
)

func TestAPI_GetPost(t *testing.T) {
	tests := []struct {
		name            string
		mockStorage     func(t *testing.T) *mocks.Storage
		req             *ppbapi.GetPostRequest
		expectedResp    *ppbapi.GetPostResponse
		wantErr         bool
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name: "OK",
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				postFromStorage := model.Post{
					ID:          1,
					UserID:      1,
					Description: "TestDescription",
					HashtagsID:  []int64{1, 2, 3},
					CreatedAt:   time.Unix(10, 20),
				}
				testStorage.On("GetPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1)).
					Return(postFromStorage, nil).
					Once()
				return testStorage
			},
			req: &ppbapi.GetPostRequest{Id: int64(1)},
			expectedResp: &ppbapi.GetPostResponse{
				UserID:      1,
				Description: "TestDescription",
				HashtagsID:  []int64{1, 2, 3},
				CreatedAt:   timestamppb.New(time.Unix(10, 20)),
			},
		},
		{
			name: "not found",
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("GetPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1)).
					Return(model.Post{}, sql.ErrNoRows).
					Once()
				return testStorage
			},
			req:             &ppbapi.GetPostRequest{Id: int64(1)},
			wantErr:         true,
			expectedErr:     ppbapi.ErrPostNotFoundByID,
			expectedErrCode: codes.NotFound,
		},
		{
			name: "unexpected err on storage.GetPost",
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("GetPost",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1)).
					Return(model.Post{}, errors.New("unexpected err")).
					Once()
				return testStorage
			},
			req:             &ppbapi.GetPostRequest{Id: int64(1)},
			wantErr:         true,
			expectedErr:     errors.New("unexpected"),
			expectedErrCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{storage: tt.mockStorage(t)}
			resp, err := a.GetPost(context.Background(), tt.req)

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

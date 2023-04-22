package api

import (
	"context"
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
	"gitlab.com/pet-pr-social-network/post-service/pbapi"
)

func TestAPI_GetPostsByHashtag(t *testing.T) {
	tests := []struct {
		name            string
		mockStorage     func(t *testing.T) *mocks.Storage
		req             *pbapi.GetPostsByHashtagRequest
		expectedResp    *pbapi.GetPostsByHashtagResponse
		wantErr         bool
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name: "OK",
			req:  &pbapi.GetPostsByHashtagRequest{HashtagID: 1, Direction: 0, PostOffsetID: 0, Limit: 10},
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				postsFromStorage := []model.Post{
					{
						ID:          1,
						UserID:      1,
						Description: "TestDescription1",
						HashtagsID:  []int64{1, 2},
						CreatedAt:   time.Unix(10, 20),
					},
					{
						ID:          2,
						UserID:      2,
						Description: "TestDescription2",
						HashtagsID:  []int64{3, 4},
						CreatedAt:   time.Unix(30, 40),
					},
				}
				testStorage.On("GetPostsByHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1), model.Direction(0), int64(0), int64(10)).
					Return(postsFromStorage, nil).
					Once()
				return testStorage
			},
			expectedResp: &pbapi.GetPostsByHashtagResponse{
				Posts: []*pbapi.Post{
					{
						Id:          1,
						UserID:      1,
						Description: "TestDescription1",
						HashtagsID:  []int64{1, 2},
						CreatedAt:   timestamppb.New(time.Unix(10, 20)),
					},
					{
						Id:          2,
						UserID:      2,
						Description: "TestDescription2",
						HashtagsID:  []int64{3, 4},
						CreatedAt:   timestamppb.New(time.Unix(30, 40)),
					},
				},
			},
		},
		{
			name: "unexpected err on storage.GetPostsByHashtag",
			req:  &pbapi.GetPostsByHashtagRequest{HashtagID: 1, Direction: 0, PostOffsetID: 0, Limit: 10},
			mockStorage: func(t *testing.T) *mocks.Storage {
				testStorage := mocks.NewStorage(t)
				testStorage.On("GetPostsByHashtag",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1), model.Direction(0), int64(0), int64(10)).
					Return([]model.Post{}, errors.New("unexpected err")).
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
			resp, err := a.GetPostsByHashtag(context.Background(), tt.req)

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

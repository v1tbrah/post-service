// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "gitlab.com/pet-pr-social-network/post-service/internal/model"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// AddHashtagToPost provides a mock function with given fields: ctx, postID, hashtagID
func (_m *Storage) AddHashtagToPost(ctx context.Context, postID int64, hashtagID int64) error {
	ret := _m.Called(ctx, postID, hashtagID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, postID, hashtagID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateHashtag provides a mock function with given fields: ctx, hashtag
func (_m *Storage) CreateHashtag(ctx context.Context, hashtag model.Hashtag) (int64, error) {
	ret := _m.Called(ctx, hashtag)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Hashtag) (int64, error)); ok {
		return rf(ctx, hashtag)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.Hashtag) int64); ok {
		r0 = rf(ctx, hashtag)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.Hashtag) error); ok {
		r1 = rf(ctx, hashtag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePost provides a mock function with given fields: ctx, post
func (_m *Storage) CreatePost(ctx context.Context, post model.Post) (int64, error) {
	ret := _m.Called(ctx, post)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Post) (int64, error)); ok {
		return rf(ctx, post)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.Post) int64); ok {
		r0 = rf(ctx, post)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.Post) error); ok {
		r1 = rf(ctx, post)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHashtag provides a mock function with given fields: ctx, id
func (_m *Storage) GetHashtag(ctx context.Context, id int64) (model.Hashtag, error) {
	ret := _m.Called(ctx, id)

	var r0 model.Hashtag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (model.Hashtag, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) model.Hashtag); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(model.Hashtag)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPost provides a mock function with given fields: ctx, id
func (_m *Storage) GetPost(ctx context.Context, id int64) (model.Post, error) {
	ret := _m.Called(ctx, id)

	var r0 model.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (model.Post, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) model.Post); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(model.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPostsByHashtag provides a mock function with given fields: ctx, hashtagID, direction, postOffsetID, limit
func (_m *Storage) GetPostsByHashtag(ctx context.Context, hashtagID int64, direction model.Direction, postOffsetID int64, limit int64) ([]model.Post, error) {
	ret := _m.Called(ctx, hashtagID, direction, postOffsetID, limit)

	var r0 []model.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, model.Direction, int64, int64) ([]model.Post, error)); ok {
		return rf(ctx, hashtagID, direction, postOffsetID, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, model.Direction, int64, int64) []model.Post); ok {
		r0 = rf(ctx, hashtagID, direction, postOffsetID, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, model.Direction, int64, int64) error); ok {
		r1 = rf(ctx, hashtagID, direction, postOffsetID, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPostsByUserID provides a mock function with given fields: ctx, userID
func (_m *Storage) GetPostsByUserID(ctx context.Context, userID int64) ([]model.Post, error) {
	ret := _m.Called(ctx, userID)

	var r0 []model.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]model.Post, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []model.Post); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorage(t mockConstructorTestingTNewStorage) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

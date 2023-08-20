// Code generated by mockery v2.32.4. DO NOT EDIT.

package mockery

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	git "github.com/walteh/buildrc/pkg/git"
)

// MockPullRequestProvider_git is an autogenerated mock type for the PullRequestProvider type
type MockPullRequestProvider_git struct {
	mock.Mock
}

type MockPullRequestProvider_git_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPullRequestProvider_git) EXPECT() *MockPullRequestProvider_git_Expecter {
	return &MockPullRequestProvider_git_Expecter{mock: &_m.Mock}
}

// ListRecentPullRequests provides a mock function with given fields: ctx, head
func (_m *MockPullRequestProvider_git) ListRecentPullRequests(ctx context.Context, head string) ([]*git.PullRequest, error) {
	ret := _m.Called(ctx, head)

	var r0 []*git.PullRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*git.PullRequest, error)); ok {
		return rf(ctx, head)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*git.PullRequest); ok {
		r0 = rf(ctx, head)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*git.PullRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, head)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPullRequestProvider_git_ListRecentPullRequests_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListRecentPullRequests'
type MockPullRequestProvider_git_ListRecentPullRequests_Call struct {
	*mock.Call
}

// ListRecentPullRequests is a helper method to define mock.On call
//   - ctx context.Context
//   - head string
func (_e *MockPullRequestProvider_git_Expecter) ListRecentPullRequests(ctx interface{}, head interface{}) *MockPullRequestProvider_git_ListRecentPullRequests_Call {
	return &MockPullRequestProvider_git_ListRecentPullRequests_Call{Call: _e.mock.On("ListRecentPullRequests", ctx, head)}
}

func (_c *MockPullRequestProvider_git_ListRecentPullRequests_Call) Run(run func(ctx context.Context, head string)) *MockPullRequestProvider_git_ListRecentPullRequests_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockPullRequestProvider_git_ListRecentPullRequests_Call) Return(_a0 []*git.PullRequest, _a1 error) *MockPullRequestProvider_git_ListRecentPullRequests_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPullRequestProvider_git_ListRecentPullRequests_Call) RunAndReturn(run func(context.Context, string) ([]*git.PullRequest, error)) *MockPullRequestProvider_git_ListRecentPullRequests_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPullRequestProvider_git creates a new instance of MockPullRequestProvider_git. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPullRequestProvider_git(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPullRequestProvider_git {
	mock := &MockPullRequestProvider_git{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

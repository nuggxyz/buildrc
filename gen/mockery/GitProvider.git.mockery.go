// Code generated by mockery v2.33.0. DO NOT EDIT.

package mockery

import (
	context "context"

	afero "github.com/spf13/afero"

	git "github.com/walteh/buildrc/pkg/git"

	mock "github.com/stretchr/testify/mock"

	semver "github.com/Masterminds/semver/v3"
)

// MockGitProvider_git is an autogenerated mock type for the GitProvider type
type MockGitProvider_git struct {
	mock.Mock
}

type MockGitProvider_git_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGitProvider_git) EXPECT() *MockGitProvider_git_Expecter {
	return &MockGitProvider_git_Expecter{mock: &_m.Mock}
}

// Dirty provides a mock function with given fields: ctx
func (_m *MockGitProvider_git) Dirty(ctx context.Context) bool {
	ret := _m.Called(ctx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockGitProvider_git_Dirty_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Dirty'
type MockGitProvider_git_Dirty_Call struct {
	*mock.Call
}

// Dirty is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockGitProvider_git_Expecter) Dirty(ctx interface{}) *MockGitProvider_git_Dirty_Call {
	return &MockGitProvider_git_Dirty_Call{Call: _e.mock.On("Dirty", ctx)}
}

func (_c *MockGitProvider_git_Dirty_Call) Run(run func(ctx context.Context)) *MockGitProvider_git_Dirty_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockGitProvider_git_Dirty_Call) Return(_a0 bool) *MockGitProvider_git_Dirty_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGitProvider_git_Dirty_Call) RunAndReturn(run func(context.Context) bool) *MockGitProvider_git_Dirty_Call {
	_c.Call.Return(run)
	return _c
}

// Fs provides a mock function with given fields:
func (_m *MockGitProvider_git) Fs() afero.Fs {
	ret := _m.Called()

	var r0 afero.Fs
	if rf, ok := ret.Get(0).(func() afero.Fs); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(afero.Fs)
		}
	}

	return r0
}

// MockGitProvider_git_Fs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fs'
type MockGitProvider_git_Fs_Call struct {
	*mock.Call
}

// Fs is a helper method to define mock.On call
func (_e *MockGitProvider_git_Expecter) Fs() *MockGitProvider_git_Fs_Call {
	return &MockGitProvider_git_Fs_Call{Call: _e.mock.On("Fs")}
}

func (_c *MockGitProvider_git_Fs_Call) Run(run func()) *MockGitProvider_git_Fs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockGitProvider_git_Fs_Call) Return(_a0 afero.Fs) *MockGitProvider_git_Fs_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGitProvider_git_Fs_Call) RunAndReturn(run func() afero.Fs) *MockGitProvider_git_Fs_Call {
	_c.Call.Return(run)
	return _c
}

// GetContentHashFromRef provides a mock function with given fields: ctx, ref
func (_m *MockGitProvider_git) GetContentHashFromRef(ctx context.Context, ref string) (string, error) {
	ret := _m.Called(ctx, ref)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, ref)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, ref)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_GetContentHashFromRef_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetContentHashFromRef'
type MockGitProvider_git_GetContentHashFromRef_Call struct {
	*mock.Call
}

// GetContentHashFromRef is a helper method to define mock.On call
//   - ctx context.Context
//   - ref string
func (_e *MockGitProvider_git_Expecter) GetContentHashFromRef(ctx interface{}, ref interface{}) *MockGitProvider_git_GetContentHashFromRef_Call {
	return &MockGitProvider_git_GetContentHashFromRef_Call{Call: _e.mock.On("GetContentHashFromRef", ctx, ref)}
}

func (_c *MockGitProvider_git_GetContentHashFromRef_Call) Run(run func(ctx context.Context, ref string)) *MockGitProvider_git_GetContentHashFromRef_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockGitProvider_git_GetContentHashFromRef_Call) Return(_a0 string, _a1 error) *MockGitProvider_git_GetContentHashFromRef_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_GetContentHashFromRef_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockGitProvider_git_GetContentHashFromRef_Call {
	_c.Call.Return(run)
	return _c
}

// GetCurrentBranchFromRef provides a mock function with given fields: ctx, ref
func (_m *MockGitProvider_git) GetCurrentBranchFromRef(ctx context.Context, ref string) (string, error) {
	ret := _m.Called(ctx, ref)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, ref)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, ref)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_GetCurrentBranchFromRef_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrentBranchFromRef'
type MockGitProvider_git_GetCurrentBranchFromRef_Call struct {
	*mock.Call
}

// GetCurrentBranchFromRef is a helper method to define mock.On call
//   - ctx context.Context
//   - ref string
func (_e *MockGitProvider_git_Expecter) GetCurrentBranchFromRef(ctx interface{}, ref interface{}) *MockGitProvider_git_GetCurrentBranchFromRef_Call {
	return &MockGitProvider_git_GetCurrentBranchFromRef_Call{Call: _e.mock.On("GetCurrentBranchFromRef", ctx, ref)}
}

func (_c *MockGitProvider_git_GetCurrentBranchFromRef_Call) Run(run func(ctx context.Context, ref string)) *MockGitProvider_git_GetCurrentBranchFromRef_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockGitProvider_git_GetCurrentBranchFromRef_Call) Return(_a0 string, _a1 error) *MockGitProvider_git_GetCurrentBranchFromRef_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_GetCurrentBranchFromRef_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockGitProvider_git_GetCurrentBranchFromRef_Call {
	_c.Call.Return(run)
	return _c
}

// GetCurrentCommitFromRef provides a mock function with given fields: ctx, ref
func (_m *MockGitProvider_git) GetCurrentCommitFromRef(ctx context.Context, ref string) (string, error) {
	ret := _m.Called(ctx, ref)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, ref)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, ref)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_GetCurrentCommitFromRef_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrentCommitFromRef'
type MockGitProvider_git_GetCurrentCommitFromRef_Call struct {
	*mock.Call
}

// GetCurrentCommitFromRef is a helper method to define mock.On call
//   - ctx context.Context
//   - ref string
func (_e *MockGitProvider_git_Expecter) GetCurrentCommitFromRef(ctx interface{}, ref interface{}) *MockGitProvider_git_GetCurrentCommitFromRef_Call {
	return &MockGitProvider_git_GetCurrentCommitFromRef_Call{Call: _e.mock.On("GetCurrentCommitFromRef", ctx, ref)}
}

func (_c *MockGitProvider_git_GetCurrentCommitFromRef_Call) Run(run func(ctx context.Context, ref string)) *MockGitProvider_git_GetCurrentCommitFromRef_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockGitProvider_git_GetCurrentCommitFromRef_Call) Return(_a0 string, _a1 error) *MockGitProvider_git_GetCurrentCommitFromRef_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_GetCurrentCommitFromRef_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockGitProvider_git_GetCurrentCommitFromRef_Call {
	_c.Call.Return(run)
	return _c
}

// GetCurrentCommitMessageFromRef provides a mock function with given fields: ctx, ref
func (_m *MockGitProvider_git) GetCurrentCommitMessageFromRef(ctx context.Context, ref string) (string, error) {
	ret := _m.Called(ctx, ref)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, ref)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, ref)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_GetCurrentCommitMessageFromRef_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrentCommitMessageFromRef'
type MockGitProvider_git_GetCurrentCommitMessageFromRef_Call struct {
	*mock.Call
}

// GetCurrentCommitMessageFromRef is a helper method to define mock.On call
//   - ctx context.Context
//   - ref string
func (_e *MockGitProvider_git_Expecter) GetCurrentCommitMessageFromRef(ctx interface{}, ref interface{}) *MockGitProvider_git_GetCurrentCommitMessageFromRef_Call {
	return &MockGitProvider_git_GetCurrentCommitMessageFromRef_Call{Call: _e.mock.On("GetCurrentCommitMessageFromRef", ctx, ref)}
}

func (_c *MockGitProvider_git_GetCurrentCommitMessageFromRef_Call) Run(run func(ctx context.Context, ref string)) *MockGitProvider_git_GetCurrentCommitMessageFromRef_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockGitProvider_git_GetCurrentCommitMessageFromRef_Call) Return(_a0 string, _a1 error) *MockGitProvider_git_GetCurrentCommitMessageFromRef_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_GetCurrentCommitMessageFromRef_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockGitProvider_git_GetCurrentCommitMessageFromRef_Call {
	_c.Call.Return(run)
	return _c
}

// GetCurrentShortHashFromRef provides a mock function with given fields: ctx, ref
func (_m *MockGitProvider_git) GetCurrentShortHashFromRef(ctx context.Context, ref string) (string, error) {
	ret := _m.Called(ctx, ref)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, ref)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, ref)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_GetCurrentShortHashFromRef_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrentShortHashFromRef'
type MockGitProvider_git_GetCurrentShortHashFromRef_Call struct {
	*mock.Call
}

// GetCurrentShortHashFromRef is a helper method to define mock.On call
//   - ctx context.Context
//   - ref string
func (_e *MockGitProvider_git_Expecter) GetCurrentShortHashFromRef(ctx interface{}, ref interface{}) *MockGitProvider_git_GetCurrentShortHashFromRef_Call {
	return &MockGitProvider_git_GetCurrentShortHashFromRef_Call{Call: _e.mock.On("GetCurrentShortHashFromRef", ctx, ref)}
}

func (_c *MockGitProvider_git_GetCurrentShortHashFromRef_Call) Run(run func(ctx context.Context, ref string)) *MockGitProvider_git_GetCurrentShortHashFromRef_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockGitProvider_git_GetCurrentShortHashFromRef_Call) Return(_a0 string, _a1 error) *MockGitProvider_git_GetCurrentShortHashFromRef_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_GetCurrentShortHashFromRef_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockGitProvider_git_GetCurrentShortHashFromRef_Call {
	_c.Call.Return(run)
	return _c
}

// GetLatestSemverTagFromRef provides a mock function with given fields: ctx, ref
func (_m *MockGitProvider_git) GetLatestSemverTagFromRef(ctx context.Context, ref string) (*semver.Version, error) {
	ret := _m.Called(ctx, ref)

	var r0 *semver.Version
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*semver.Version, error)); ok {
		return rf(ctx, ref)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *semver.Version); ok {
		r0 = rf(ctx, ref)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*semver.Version)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_GetLatestSemverTagFromRef_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatestSemverTagFromRef'
type MockGitProvider_git_GetLatestSemverTagFromRef_Call struct {
	*mock.Call
}

// GetLatestSemverTagFromRef is a helper method to define mock.On call
//   - ctx context.Context
//   - ref string
func (_e *MockGitProvider_git_Expecter) GetLatestSemverTagFromRef(ctx interface{}, ref interface{}) *MockGitProvider_git_GetLatestSemverTagFromRef_Call {
	return &MockGitProvider_git_GetLatestSemverTagFromRef_Call{Call: _e.mock.On("GetLatestSemverTagFromRef", ctx, ref)}
}

func (_c *MockGitProvider_git_GetLatestSemverTagFromRef_Call) Run(run func(ctx context.Context, ref string)) *MockGitProvider_git_GetLatestSemverTagFromRef_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockGitProvider_git_GetLatestSemverTagFromRef_Call) Return(_a0 *semver.Version, _a1 error) *MockGitProvider_git_GetLatestSemverTagFromRef_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_GetLatestSemverTagFromRef_Call) RunAndReturn(run func(context.Context, string) (*semver.Version, error)) *MockGitProvider_git_GetLatestSemverTagFromRef_Call {
	_c.Call.Return(run)
	return _c
}

// GetLocalRepositoryMetadata provides a mock function with given fields: ctx
func (_m *MockGitProvider_git) GetLocalRepositoryMetadata(ctx context.Context) (*git.LocalRepositoryMetadata, error) {
	ret := _m.Called(ctx)

	var r0 *git.LocalRepositoryMetadata
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*git.LocalRepositoryMetadata, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *git.LocalRepositoryMetadata); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*git.LocalRepositoryMetadata)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_GetLocalRepositoryMetadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLocalRepositoryMetadata'
type MockGitProvider_git_GetLocalRepositoryMetadata_Call struct {
	*mock.Call
}

// GetLocalRepositoryMetadata is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockGitProvider_git_Expecter) GetLocalRepositoryMetadata(ctx interface{}) *MockGitProvider_git_GetLocalRepositoryMetadata_Call {
	return &MockGitProvider_git_GetLocalRepositoryMetadata_Call{Call: _e.mock.On("GetLocalRepositoryMetadata", ctx)}
}

func (_c *MockGitProvider_git_GetLocalRepositoryMetadata_Call) Run(run func(ctx context.Context)) *MockGitProvider_git_GetLocalRepositoryMetadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockGitProvider_git_GetLocalRepositoryMetadata_Call) Return(_a0 *git.LocalRepositoryMetadata, _a1 error) *MockGitProvider_git_GetLocalRepositoryMetadata_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_GetLocalRepositoryMetadata_Call) RunAndReturn(run func(context.Context) (*git.LocalRepositoryMetadata, error)) *MockGitProvider_git_GetLocalRepositoryMetadata_Call {
	_c.Call.Return(run)
	return _c
}

// GetRemoteURL provides a mock function with given fields: ctx
func (_m *MockGitProvider_git) GetRemoteURL(ctx context.Context) (string, error) {
	ret := _m.Called(ctx)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_GetRemoteURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRemoteURL'
type MockGitProvider_git_GetRemoteURL_Call struct {
	*mock.Call
}

// GetRemoteURL is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockGitProvider_git_Expecter) GetRemoteURL(ctx interface{}) *MockGitProvider_git_GetRemoteURL_Call {
	return &MockGitProvider_git_GetRemoteURL_Call{Call: _e.mock.On("GetRemoteURL", ctx)}
}

func (_c *MockGitProvider_git_GetRemoteURL_Call) Run(run func(ctx context.Context)) *MockGitProvider_git_GetRemoteURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockGitProvider_git_GetRemoteURL_Call) Return(_a0 string, _a1 error) *MockGitProvider_git_GetRemoteURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_GetRemoteURL_Call) RunAndReturn(run func(context.Context) (string, error)) *MockGitProvider_git_GetRemoteURL_Call {
	_c.Call.Return(run)
	return _c
}

// TryGetPRNumber provides a mock function with given fields: ctx
func (_m *MockGitProvider_git) TryGetPRNumber(ctx context.Context) (uint64, error) {
	ret := _m.Called(ctx)

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (uint64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) uint64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_TryGetPRNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TryGetPRNumber'
type MockGitProvider_git_TryGetPRNumber_Call struct {
	*mock.Call
}

// TryGetPRNumber is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockGitProvider_git_Expecter) TryGetPRNumber(ctx interface{}) *MockGitProvider_git_TryGetPRNumber_Call {
	return &MockGitProvider_git_TryGetPRNumber_Call{Call: _e.mock.On("TryGetPRNumber", ctx)}
}

func (_c *MockGitProvider_git_TryGetPRNumber_Call) Run(run func(ctx context.Context)) *MockGitProvider_git_TryGetPRNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockGitProvider_git_TryGetPRNumber_Call) Return(_a0 uint64, _a1 error) *MockGitProvider_git_TryGetPRNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_TryGetPRNumber_Call) RunAndReturn(run func(context.Context) (uint64, error)) *MockGitProvider_git_TryGetPRNumber_Call {
	_c.Call.Return(run)
	return _c
}

// TryGetSemverTag provides a mock function with given fields: ctx
func (_m *MockGitProvider_git) TryGetSemverTag(ctx context.Context) (*semver.Version, error) {
	ret := _m.Called(ctx)

	var r0 *semver.Version
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*semver.Version, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *semver.Version); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*semver.Version)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_git_TryGetSemverTag_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TryGetSemverTag'
type MockGitProvider_git_TryGetSemverTag_Call struct {
	*mock.Call
}

// TryGetSemverTag is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockGitProvider_git_Expecter) TryGetSemverTag(ctx interface{}) *MockGitProvider_git_TryGetSemverTag_Call {
	return &MockGitProvider_git_TryGetSemverTag_Call{Call: _e.mock.On("TryGetSemverTag", ctx)}
}

func (_c *MockGitProvider_git_TryGetSemverTag_Call) Run(run func(ctx context.Context)) *MockGitProvider_git_TryGetSemverTag_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockGitProvider_git_TryGetSemverTag_Call) Return(_a0 *semver.Version, _a1 error) *MockGitProvider_git_TryGetSemverTag_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGitProvider_git_TryGetSemverTag_Call) RunAndReturn(run func(context.Context) (*semver.Version, error)) *MockGitProvider_git_TryGetSemverTag_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGitProvider_git creates a new instance of MockGitProvider_git. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGitProvider_git(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGitProvider_git {
	mock := &MockGitProvider_git{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

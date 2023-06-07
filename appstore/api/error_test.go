package api

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_As(t *testing.T) {
	tests := []struct {
		SrcError   error
		ExpectedAs bool // Check If SrcError can be 'As' to Error.
	}{
		{SrcError: AccountNotFoundError, ExpectedAs: true},
		{SrcError: AppNotFoundError, ExpectedAs: true},
		{SrcError: fmt.Errorf("custom error"), ExpectedAs: false},
		{SrcError: fmt.Errorf("wrapping: %w", AccountNotFoundError), ExpectedAs: true},
		{SrcError: errors.Unwrap(fmt.Errorf("wrapping: %w", AccountNotFoundError)), ExpectedAs: true},
	}

	for _, test := range tests {
		var apiErr *Error
		as := errors.As(test.SrcError, &apiErr)
		assert.Equal(t, test.ExpectedAs, as)
		if test.ExpectedAs {
			assert.NotZero(t, apiErr.errorCode)
			assert.NotZero(t, apiErr.errorMessage)
		} else {
			assert.Nil(t, apiErr)
		}
	}

}

func TestError_Is(t *testing.T) {
	tests := []struct {
		ErrBytes    []byte
		TargetError error
		ExpectedIs  bool // Check if error (constructed by ErrBytes) Is TargetError Or not.
	}{
		{ErrBytes: []byte(`{"errorCode": 4040001, "errorMessage": "Account not found."}`), TargetError: AccountNotFoundError, ExpectedIs: true},
		{ErrBytes: []byte(`{"errorCode": 4040001, "errorMessage": "Account not found."}`), TargetError: AppNotFoundError, ExpectedIs: false},
		{ErrBytes: []byte(`{"errorCode": 4040001, "errorMessage": "Account not found."}`), TargetError: fmt.Errorf("custom error"), ExpectedIs: false},
		{ErrBytes: []byte(`{"errorCode": 4040001, "errorMessage": "Account not found."}`), TargetError: fmt.Errorf("wrapping: %w", AccountNotFoundError), ExpectedIs: false},
		{ErrBytes: []byte(`{"errorCode": 4040001, "errorMessage": "Account not found."}`), TargetError: errors.Unwrap(fmt.Errorf("wrapping: %w", AccountNotFoundError)), ExpectedIs: true},
	}
	for _, test := range tests {
		err, ok := newErrorFromJSON(test.ErrBytes)
		assert.True(t, ok)
		assert.Equal(t, test.ExpectedIs, errors.Is(err, test.TargetError))
	}
}

func TestError_Is2(t *testing.T) {
	tests := []struct {
		SrcError    error
		TargetError error
		ExpectedIs  bool // Check if SrcError is TargetError or not.
	}{
		{SrcError: AccountNotFoundError, TargetError: AccountNotFoundError, ExpectedIs: true},
		{SrcError: AppNotFoundError, TargetError: AccountNotFoundError, ExpectedIs: false},
		{SrcError: fmt.Errorf("custom error"), TargetError: AccountNotFoundError, ExpectedIs: false},
		{SrcError: fmt.Errorf("wrapping: %w", AccountNotFoundError), TargetError: AccountNotFoundError, ExpectedIs: true},
		{SrcError: errors.Unwrap(fmt.Errorf("wrapping: %w", AccountNotFoundError)), TargetError: AccountNotFoundError, ExpectedIs: true},
	}
	for _, test := range tests {
		assert.Equal(t, test.ExpectedIs, errors.Is(test.SrcError, test.TargetError))
	}
}

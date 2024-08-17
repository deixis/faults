package faults_test

import (
	"testing"

	"github.com/deixis/faults"
)

// TestIs ensures all `Is*` functions return true for the error they are
// supposed to match.
func TestIs(t *testing.T) {
	table := []struct {
		Error error
		Is    func(err error) bool
	}{
		{
			Error: faults.NotFound,
			Is:    faults.IsNotFound,
		},
		{
			Error: faults.PermissionDenied,
			Is:    faults.IsPermissionDenied,
		},
		{
			Error: faults.Unauthenticated,
			Is:    faults.IsUnauthenticated,
		},
		{
			Error: faults.Bad(),
			Is:    faults.IsBad,
		},
		{
			Error: faults.FailedPrecondition(),
			Is:    faults.IsFailedPrecondition,
		},
		{
			Error: faults.Aborted(),
			Is:    faults.IsAborted,
		},
		{
			Error: faults.Unavailable(0),
			Is:    faults.IsUnavailable,
		},
		{
			Error: faults.ResourceExhausted(),
			Is:    faults.IsResourceExhausted,
		},
	}

	for i, test := range table {
		if !test.Is(test.Error) {
			t.Errorf("%d - expect error Is to return true for error %s", i, test.Error)
		}
	}
}

// TestAs ensures all `As*` functions return true for the error they are
// supposed to match.
func TestAs(t *testing.T) {
	table := []struct {
		Error error
		As    func(err error) bool
	}{
		{
			Error: faults.NotFound,
			As: func(err error) bool {
				_, ok := faults.AsNotFound(err)
				return ok
			},
		},
		{
			Error: faults.PermissionDenied,
			As: func(err error) bool {
				_, ok := faults.AsPermissionDenied(err)
				return ok
			},
		},
		{
			Error: faults.Unauthenticated,
			As: func(err error) bool {
				_, ok := faults.AsUnauthenticated(err)
				return ok
			},
		},
		{
			Error: faults.Bad(),
			As: func(err error) bool {
				_, ok := faults.AsBad(err)
				return ok
			},
		},
		{
			Error: faults.FailedPrecondition(),
			As: func(err error) bool {
				_, ok := faults.AsFailedPrecondition(err)
				return ok
			},
		},
		{
			Error: faults.Aborted(),
			As: func(err error) bool {
				_, ok := faults.AsAborted(err)
				return ok
			},
		},
		{
			Error: faults.Unavailable(0),
			As: func(err error) bool {
				_, ok := faults.AsUnavailable(err)
				return ok
			},
		},
		{
			Error: faults.ResourceExhausted(),
			As: func(err error) bool {
				_, ok := faults.AsResourceExhausted(err)
				return ok
			},
		},
	}

	for i, test := range table {
		if !test.As(test.Error) {
			t.Errorf("%d - expect error As to return true for error %s", i, test.Error)
		}
	}
}

package faults_test

import (
	"errors"
	"testing"
	"time"

	"github.com/deixis/faults"
)

var (
	// Pre-created errors for benchmarking
	errBase            = errors.New("base error")
	permissionErr      = faults.WithPermissionDenied(errBase)
	unauthenticatedErr = faults.WithUnauthenticated(errBase)
	notFoundErr        = faults.WithNotFound(errBase)
	badErr             = faults.WithBad(errBase, &faults.FieldViolation{Field: "test", Description: "invalid"})
	preconditionErr    = faults.WithFailedPrecondition(errBase, &faults.PreconditionViolation{Type: "test", Subject: "test", Description: "failed"})
	abortedErr         = faults.WithAborted(errBase, &faults.ConflictViolation{Resource: "test", Description: "conflict"})
	unavailableErr     = faults.WithUnavailable(errBase, time.Second)
	exhaustedErr       = faults.WithResourceExhausted(errBase, &faults.QuotaViolation{Subject: "test", Description: "exhausted"})
	unimplementedErr   = faults.WithUnimplemented(errBase)

	// Sink variables to prevent compiler optimizations
	resultError  error
	resultBool   bool
	resultString string
)

// Benchmark With* functions
func BenchmarkWithPermissionDenied(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultError = faults.WithPermissionDenied(errBase)
	}
}

func BenchmarkWithUnauthenticated(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultError = faults.WithUnauthenticated(errBase)
	}
}

func BenchmarkWithNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultError = faults.WithNotFound(errBase)
	}
}

func BenchmarkWithBad(b *testing.B) {
	violation := &faults.FieldViolation{Field: "test", Description: "invalid"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultError = faults.WithBad(errBase, violation)
	}
}

func BenchmarkWithFailedPrecondition(b *testing.B) {
	violation := &faults.PreconditionViolation{Type: "test", Subject: "test", Description: "failed"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultError = faults.WithFailedPrecondition(errBase, violation)
	}
}

func BenchmarkWithAborted(b *testing.B) {
	violation := &faults.ConflictViolation{Resource: "test", Description: "conflict"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultError = faults.WithAborted(errBase, violation)
	}
}

func BenchmarkWithUnavailable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultError = faults.WithUnavailable(errBase, time.Second)
	}
}

func BenchmarkWithResourceExhausted(b *testing.B) {
	violation := &faults.QuotaViolation{Subject: "test", Description: "exhausted"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultError = faults.WithResourceExhausted(errBase, violation)
	}
}

func BenchmarkWithUnimplemented(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultError = faults.WithUnimplemented(errBase)
	}
}

// Benchmark constructor functions
func BenchmarkBad(b *testing.B) {
	violation := &faults.FieldViolation{Field: "test", Description: "invalid"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultError = faults.Bad(violation)
	}
}

func BenchmarkFailedPrecondition(b *testing.B) {
	violation := &faults.PreconditionViolation{Type: "test", Subject: "test", Description: "failed"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultError = faults.FailedPrecondition(violation)
	}
}

func BenchmarkAborted(b *testing.B) {
	violation := &faults.ConflictViolation{Resource: "test", Description: "conflict"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultError = faults.Aborted(violation)
	}
}

func BenchmarkUnavailable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultError = faults.Unavailable(time.Second)
	}
}

func BenchmarkResourceExhausted(b *testing.B) {
	violation := &faults.QuotaViolation{Subject: "test", Description: "exhausted"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultError = faults.ResourceExhausted(violation)
	}
}

// Benchmark Is* functions
func BenchmarkIsPermissionDenied(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsPermissionDenied(permissionErr)
	}
}

func BenchmarkIsUnauthenticated(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsUnauthenticated(unauthenticatedErr)
	}
}

func BenchmarkIsNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsNotFound(notFoundErr)
	}
}

func BenchmarkIsBad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsBad(badErr)
	}
}

func BenchmarkIsFailedPrecondition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsFailedPrecondition(preconditionErr)
	}
}

func BenchmarkIsAborted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsAborted(abortedErr)
	}
}

func BenchmarkIsUnavailable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsUnavailable(unavailableErr)
	}
}

func BenchmarkIsResourceExhausted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsResourceExhausted(exhaustedErr)
	}
}

func BenchmarkIsUnimplemented(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsUnimplemented(unimplementedErr)
	}
}

// Benchmark As* functions
func BenchmarkAsPermissionDenied(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsPermissionDenied(permissionErr)
	}
}

func BenchmarkAsUnauthenticated(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsUnauthenticated(unauthenticatedErr)
	}
}

func BenchmarkAsNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsNotFound(notFoundErr)
	}
}

func BenchmarkAsBad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsBad(badErr)
	}
}

func BenchmarkAsFailedPrecondition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsFailedPrecondition(preconditionErr)
	}
}

func BenchmarkAsAborted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsAborted(abortedErr)
	}
}

func BenchmarkAsUnavailable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsUnavailable(unavailableErr)
	}
}

func BenchmarkAsResourceExhausted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsResourceExhausted(exhaustedErr)
	}
}

func BenchmarkAsUnimplemented(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsUnimplemented(unimplementedErr)
	}
}

// Benchmark Error() methods for each failure type
func BenchmarkAvailabilityFailureError(b *testing.B) {
	err := faults.Unavailable(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = err.Error()
	}
}

func BenchmarkQuotaFailureError(b *testing.B) {
	violation := &faults.QuotaViolation{Subject: "test", Description: "exhausted"}
	err := faults.ResourceExhausted(violation)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = err.Error()
	}
}

func BenchmarkPreconditionFailureError(b *testing.B) {
	violation := &faults.PreconditionViolation{Type: "test", Subject: "test", Description: "failed"}
	err := faults.FailedPrecondition(violation)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = err.Error()
	}
}

func BenchmarkBadRequestError(b *testing.B) {
	violation := &faults.FieldViolation{Field: "test", Description: "invalid"}
	err := faults.Bad(violation)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = err.Error()
	}
}

func BenchmarkConflictFailureError(b *testing.B) {
	violation := &faults.ConflictViolation{Resource: "test", Description: "conflict"}
	err := faults.Aborted(violation)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = err.Error()
	}
}

func BenchmarkMissingFailureError(b *testing.B) {
	err := faults.NotFound
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = err.Error()
	}
}

func BenchmarkPermissionFailureError(b *testing.B) {
	err := faults.PermissionDenied
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = err.Error()
	}
}

func BenchmarkAuthenticationFailureError(b *testing.B) {
	err := faults.Unauthenticated
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = err.Error()
	}
}

func BenchmarkUnimplementedFailureError(b *testing.B) {
	err := faults.Unimplemented
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = err.Error()
	}
}

// Benchmark String() methods for violation types
func BenchmarkQuotaViolationString(b *testing.B) {
	violation := &faults.QuotaViolation{Subject: "clientip:192.168.1.1", Description: "Daily limit exceeded"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = violation.String()
	}
}

func BenchmarkPreconditionViolationString(b *testing.B) {
	violation := &faults.PreconditionViolation{Type: "TOS", Subject: "google.com/cloud", Description: "Terms of service not accepted"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = violation.String()
	}
}

func BenchmarkFieldViolationString(b *testing.B) {
	violation := &faults.FieldViolation{Field: "field_violations.field", Description: "Invalid field value"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = violation.String()
	}
}

func BenchmarkConflictViolationString(b *testing.B) {
	violation := &faults.ConflictViolation{Resource: "user:123e4567-e89b-12d3-a456-426614174000", Description: "Resource is locked"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultString = violation.String()
	}
}

// Benchmark complex scenarios with multiple violations
func BenchmarkBadRequestWithMultipleViolations(b *testing.B) {
	violations := []*faults.FieldViolation{
		{Field: "name", Description: "Name is required"},
		{Field: "email", Description: "Invalid email format"},
		{Field: "age", Description: "Age must be positive"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := faults.Bad(violations...)
		resultString = err.Error()
	}
}

func BenchmarkQuotaFailureWithMultipleViolations(b *testing.B) {
	violations := []*faults.QuotaViolation{
		{Subject: "clientip:192.168.1.1", Description: "Daily API limit exceeded"},
		{Subject: "project:my-project", Description: "Storage quota exceeded"},
		{Subject: "user:john@example.com", Description: "Request rate limit exceeded"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := faults.ResourceExhausted(violations...)
		resultString = err.Error()
	}
}

// Benchmark error wrapping scenarios
func BenchmarkNestedErrorWrapping(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := errors.New("original error")
		err = faults.WithNotFound(err)
		err = faults.WithUnauthenticated(err)
		resultString = err.Error()
	}
}

// Benchmark Is/As functions with non-matching errors
func BenchmarkIsNotFoundWithWrongError(b *testing.B) {
	wrongErr := faults.PermissionDenied
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultBool = faults.IsNotFound(wrongErr)
	}
}

func BenchmarkAsNotFoundWithWrongError(b *testing.B) {
	wrongErr := faults.PermissionDenied
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, resultBool = faults.AsNotFound(wrongErr)
	}
}

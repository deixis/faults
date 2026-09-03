package faults_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/deixis/faults"
)

// TestIs validates that each Is* function correctly identifies its error type,
// including through wrapping layers and negative cases.
func TestIs(t *testing.T) {
	cause := fmt.Errorf("root cause")

	tests := []struct {
		scenario string
		err      error
		is       func(error) bool
		want     bool
	}{
		// Singletons
		{
			scenario: "NotFound singleton is recognized",
			err:      faults.NotFound,
			is:       faults.IsNotFound,
			want:     true,
		},
		{
			scenario: "PermissionDenied singleton is recognized",
			err:      faults.PermissionDenied,
			is:       faults.IsPermissionDenied,
			want:     true,
		},
		{
			scenario: "Unauthenticated singleton is recognized",
			err:      faults.Unauthenticated,
			is:       faults.IsUnauthenticated,
			want:     true,
		},
		{
			scenario: "Unimplemented singleton is recognized",
			err:      faults.Unimplemented,
			is:       faults.IsUnimplemented,
			want:     true,
		},
		// Constructors
		{
			scenario: "Bad error without violations is recognized",
			err:      faults.Bad(),
			is:       faults.IsBad,
			want:     true,
		},
		{
			scenario: "Bad error with field violations is recognized",
			err:      faults.Bad(&faults.FieldViolation{Field: "email", Description: "invalid format"}),
			is:       faults.IsBad,
			want:     true,
		},
		{
			scenario: "FailedPrecondition error is recognized",
			err:      faults.FailedPrecondition(),
			is:       faults.IsFailedPrecondition,
			want:     true,
		},
		{
			scenario: "Aborted error is recognized",
			err:      faults.Aborted(),
			is:       faults.IsAborted,
			want:     true,
		},
		{
			scenario: "Unavailable error is recognized",
			err:      faults.Unavailable(0),
			is:       faults.IsUnavailable,
			want:     true,
		},
		{
			scenario: "ResourceExhausted error is recognized",
			err:      faults.ResourceExhausted(),
			is:       faults.IsResourceExhausted,
			want:     true,
		},
		// With* wrappers
		{
			scenario: "WithNotFound wrapping a cause is recognized as NotFound",
			err:      faults.WithNotFound(cause),
			is:       faults.IsNotFound,
			want:     true,
		},
		{
			scenario: "WithPermissionDenied wrapping a cause is recognized as PermissionDenied",
			err:      faults.WithPermissionDenied(cause),
			is:       faults.IsPermissionDenied,
			want:     true,
		},
		{
			scenario: "WithUnauthenticated wrapping a cause is recognized as Unauthenticated",
			err:      faults.WithUnauthenticated(cause),
			is:       faults.IsUnauthenticated,
			want:     true,
		},
		{
			scenario: "WithBad wrapping a cause is recognized as Bad",
			err:      faults.WithBad(cause),
			is:       faults.IsBad,
			want:     true,
		},
		{
			scenario: "WithFailedPrecondition wrapping a cause is recognized as FailedPrecondition",
			err:      faults.WithFailedPrecondition(cause),
			is:       faults.IsFailedPrecondition,
			want:     true,
		},
		{
			scenario: "WithAborted wrapping a cause is recognized as Aborted",
			err:      faults.WithAborted(cause),
			is:       faults.IsAborted,
			want:     true,
		},
		{
			scenario: "WithUnavailable wrapping a cause is recognized as Unavailable",
			err:      faults.WithUnavailable(cause, 5*time.Second),
			is:       faults.IsUnavailable,
			want:     true,
		},
		{
			scenario: "WithResourceExhausted wrapping a cause is recognized as ResourceExhausted",
			err:      faults.WithResourceExhausted(cause),
			is:       faults.IsResourceExhausted,
			want:     true,
		},
		{
			scenario: "WithUnimplemented wrapping a cause is recognized as Unimplemented",
			err:      faults.WithUnimplemented(cause),
			is:       faults.IsUnimplemented,
			want:     true,
		},
		// Errors wrapped further with fmt.Errorf are still recognized
		{
			scenario: "ResourceExhausted wrapped in fmt.Errorf is still recognized",
			err:      fmt.Errorf("op failed: %w", faults.ResourceExhausted()),
			is:       faults.IsResourceExhausted,
			want:     true,
		},
		{
			scenario: "NotFound wrapped in fmt.Errorf is still recognized",
			err:      fmt.Errorf("op failed: %w", faults.NotFound),
			is:       faults.IsNotFound,
			want:     true,
		},
		// Negative cases
		{
			scenario: "NotFound is not recognized as PermissionDenied",
			err:      faults.NotFound,
			is:       faults.IsPermissionDenied,
			want:     false,
		},
		{
			scenario: "unrelated error is not recognized as NotFound",
			err:      fmt.Errorf("something went wrong"),
			is:       faults.IsNotFound,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			got := tt.is(tt.err)
			if got != tt.want {
				t.Errorf("got %v, want %v (err: %v)", got, tt.want, tt.err)
			}
		})
	}
}

// TestAs validates that each As* function correctly identifies and extracts
// its error type, including field values and through wrapping layers.
func TestAs(t *testing.T) {
	cause := fmt.Errorf("root cause")

	t.Run("NotFound singleton is extractable", func(t *testing.T) {
		e, ok := faults.AsNotFound(faults.NotFound)
		if !ok || e == nil {
			t.Fatal("expected AsNotFound to extract a non-nil MissingFailure")
		}
	})

	t.Run("PermissionDenied singleton is extractable", func(t *testing.T) {
		e, ok := faults.AsPermissionDenied(faults.PermissionDenied)
		if !ok || e == nil {
			t.Fatal("expected AsPermissionDenied to extract a non-nil PermissionFailure")
		}
	})

	t.Run("Unauthenticated singleton is extractable", func(t *testing.T) {
		e, ok := faults.AsUnauthenticated(faults.Unauthenticated)
		if !ok || e == nil {
			t.Fatal("expected AsUnauthenticated to extract a non-nil AuthenticationFailure")
		}
	})

	t.Run("Unimplemented singleton is extractable", func(t *testing.T) {
		e, ok := faults.AsUnimplemented(faults.Unimplemented)
		if !ok || e == nil {
			t.Fatal("expected AsUnimplemented to extract a non-nil UnimplementedFailure")
		}
	})

	t.Run("Bad error is extractable with field violations intact", func(t *testing.T) {
		violation := &faults.FieldViolation{Field: "email", Description: "invalid format"}
		e, ok := faults.AsBad(faults.Bad(violation))
		if !ok {
			t.Fatal("expected AsBad to return true")
		}
		if len(e.Violations) != 1 || e.Violations[0].Field != "email" {
			t.Errorf("unexpected violations: %v", e.Violations)
		}
	})

	t.Run("FailedPrecondition is extractable with precondition violations intact", func(t *testing.T) {
		violation := &faults.PreconditionViolation{Type: "TOS", Subject: "google.com", Description: "not accepted"}
		e, ok := faults.AsFailedPrecondition(faults.FailedPrecondition(violation))
		if !ok {
			t.Fatal("expected AsFailedPrecondition to return true")
		}
		if len(e.Violations) != 1 || e.Violations[0].Type != "TOS" {
			t.Errorf("unexpected violations: %v", e.Violations)
		}
	})

	t.Run("Aborted is extractable with conflict violations intact", func(t *testing.T) {
		violation := &faults.ConflictViolation{Resource: "user:123", Description: "concurrent update"}
		e, ok := faults.AsAborted(faults.Aborted(violation))
		if !ok {
			t.Fatal("expected AsAborted to return true")
		}
		if len(e.Violations) != 1 || e.Violations[0].Resource != "user:123" {
			t.Errorf("unexpected violations: %v", e.Violations)
		}
	})

	t.Run("Unavailable is extractable with retry delay intact", func(t *testing.T) {
		e, ok := faults.AsUnavailable(faults.Unavailable(5 * time.Second))
		if !ok {
			t.Fatal("expected AsUnavailable to return true")
		}
		if e.RetryInfo.RetryDelay != 5*time.Second {
			t.Errorf("expected retry delay 5s, got %v", e.RetryInfo.RetryDelay)
		}
	})

	t.Run("ResourceExhausted is extractable with quota violations intact", func(t *testing.T) {
		violation := &faults.QuotaViolation{Subject: "project:123", Description: "daily limit exceeded"}
		e, ok := faults.AsResourceExhausted(faults.ResourceExhausted(violation))
		if !ok {
			t.Fatal("expected AsResourceExhausted to return true")
		}
		if len(e.Violations) != 1 || e.Violations[0].Subject != "project:123" {
			t.Errorf("unexpected violations: %v", e.Violations)
		}
	})

	t.Run("WithNotFound wrapping a cause is still extractable", func(t *testing.T) {
		e, ok := faults.AsNotFound(faults.WithNotFound(cause))
		if !ok || e == nil {
			t.Fatal("expected AsNotFound to extract from a WithNotFound error")
		}
	})

	t.Run("ResourceExhausted wrapped in fmt.Errorf is still extractable", func(t *testing.T) {
		_, ok := faults.AsResourceExhausted(fmt.Errorf("op failed: %w", faults.ResourceExhausted()))
		if !ok {
			t.Fatal("expected AsResourceExhausted to extract through fmt.Errorf wrapping")
		}
	})

	t.Run("NotFound is not extractable as PermissionDenied", func(t *testing.T) {
		_, ok := faults.AsPermissionDenied(faults.NotFound)
		if ok {
			t.Fatal("expected AsPermissionDenied to return false for a NotFound error")
		}
	})

	t.Run("unrelated error is not extractable as NotFound", func(t *testing.T) {
		_, ok := faults.AsNotFound(fmt.Errorf("unrelated"))
		if ok {
			t.Fatal("expected AsNotFound to return false for an unrelated error")
		}
	})
}

// TestWith validates that With* functions wrap the parent error, preserving the
// fault type and keeping the original cause reachable in the error chain.
func TestWith(t *testing.T) {
	cause := fmt.Errorf("root cause")

	t.Run("WithNotFound preserves the parent cause in the error chain", func(t *testing.T) {
		err := faults.WithNotFound(cause)
		if !faults.IsNotFound(err) {
			t.Error("expected IsNotFound to return true")
		}
		if !errors.Is(err, cause) {
			t.Error("expected the original cause to be reachable via errors.Is")
		}
	})

	t.Run("WithPermissionDenied preserves the parent cause in the error chain", func(t *testing.T) {
		err := faults.WithPermissionDenied(cause)
		if !faults.IsPermissionDenied(err) {
			t.Error("expected IsPermissionDenied to return true")
		}
		if !errors.Is(err, cause) {
			t.Error("expected the original cause to be reachable via errors.Is")
		}
	})

	t.Run("WithUnauthenticated preserves the parent cause in the error chain", func(t *testing.T) {
		err := faults.WithUnauthenticated(cause)
		if !faults.IsUnauthenticated(err) {
			t.Error("expected IsUnauthenticated to return true")
		}
		if !errors.Is(err, cause) {
			t.Error("expected the original cause to be reachable via errors.Is")
		}
	})

	t.Run("WithBad preserves the parent cause and carries field violations", func(t *testing.T) {
		violation := &faults.FieldViolation{Field: "name", Description: "required"}
		err := faults.WithBad(cause, violation)
		if !faults.IsBad(err) {
			t.Error("expected IsBad to return true")
		}
		if !errors.Is(err, cause) {
			t.Error("expected the original cause to be reachable via errors.Is")
		}
		e, _ := faults.AsBad(err)
		if len(e.Violations) != 1 || e.Violations[0].Field != "name" {
			t.Errorf("unexpected violations: %v", e.Violations)
		}
	})

	t.Run("WithFailedPrecondition preserves the parent cause in the error chain", func(t *testing.T) {
		err := faults.WithFailedPrecondition(cause)
		if !faults.IsFailedPrecondition(err) {
			t.Error("expected IsFailedPrecondition to return true")
		}
		if !errors.Is(err, cause) {
			t.Error("expected the original cause to be reachable via errors.Is")
		}
	})

	t.Run("WithAborted preserves the parent cause in the error chain", func(t *testing.T) {
		err := faults.WithAborted(cause)
		if !faults.IsAborted(err) {
			t.Error("expected IsAborted to return true")
		}
		if !errors.Is(err, cause) {
			t.Error("expected the original cause to be reachable via errors.Is")
		}
	})

	t.Run("WithUnavailable preserves the parent cause and carries the retry delay", func(t *testing.T) {
		err := faults.WithUnavailable(cause, 3*time.Second)
		if !faults.IsUnavailable(err) {
			t.Error("expected IsUnavailable to return true")
		}
		if !errors.Is(err, cause) {
			t.Error("expected the original cause to be reachable via errors.Is")
		}
		e, _ := faults.AsUnavailable(err)
		if e.RetryInfo.RetryDelay != 3*time.Second {
			t.Errorf("expected retry delay 3s, got %v", e.RetryInfo.RetryDelay)
		}
	})

	t.Run("WithResourceExhausted preserves the parent cause in the error chain", func(t *testing.T) {
		err := faults.WithResourceExhausted(cause)
		if !faults.IsResourceExhausted(err) {
			t.Error("expected IsResourceExhausted to return true")
		}
		if !errors.Is(err, cause) {
			t.Error("expected the original cause to be reachable via errors.Is")
		}
	})

	t.Run("WithUnimplemented preserves the parent cause in the error chain", func(t *testing.T) {
		err := faults.WithUnimplemented(cause)
		if !faults.IsUnimplemented(err) {
			t.Error("expected IsUnimplemented to return true")
		}
		if !errors.Is(err, cause) {
			t.Error("expected the original cause to be reachable via errors.Is")
		}
	})
}

// TestErrorMessage validates the human-readable message produced by each error type
// under various conditions (no violations, with violations, wrapping a cause).
func TestErrorMessage(t *testing.T) {
	cause := fmt.Errorf("underlying issue")

	tests := []struct {
		scenario string
		err      error
		wantMsg  string
	}{
		// Singletons always return a static message regardless of context
		{
			scenario: "NotFound produces a not-found message",
			err:      faults.NotFound,
			wantMsg:  "resource not found",
		},
		{
			scenario: "PermissionDenied produces a permission-denied message",
			err:      faults.PermissionDenied,
			wantMsg:  "permission denied",
		},
		{
			scenario: "Unauthenticated produces an authentication-failure message",
			err:      faults.Unauthenticated,
			wantMsg:  "failed to authenticate request",
		},
		{
			scenario: "Unimplemented produces an unimplemented message",
			err:      faults.Unimplemented,
			wantMsg:  "unimplemented (yet)",
		},
		// Constructors without violations produce a base message
		{
			scenario: "Bad without violations produces a bad-request message",
			err:      faults.Bad(),
			wantMsg:  "bad request",
		},
		{
			scenario: "FailedPrecondition without violations produces a precondition-failure message",
			err:      faults.FailedPrecondition(),
			wantMsg:  "precondition failure",
		},
		{
			scenario: "Aborted without violations produces a conflict message",
			err:      faults.Aborted(),
			wantMsg:  "conflict",
		},
		{
			scenario: "Unavailable without retry delay produces a service-unavailable message",
			err:      faults.Unavailable(0),
			wantMsg:  "service temporarily unavailable",
		},
		{
			scenario: "Unavailable with retry delay includes the delay in the message",
			err:      faults.Unavailable(5 * time.Second),
			wantMsg:  "service temporarily unavailable, retry in 5s",
		},
		{
			scenario: "ResourceExhausted without violations produces a quota-failure message",
			err:      faults.ResourceExhausted(),
			wantMsg:  "quota failure",
		},
		// Constructors with violations include violation descriptions
		{
			scenario: "Bad with a field violation includes the violation description",
			err:      faults.Bad(&faults.FieldViolation{Field: "email", Description: "invalid format"}),
			wantMsg:  "invalid format",
		},
		{
			scenario: "Bad with multiple field violations joins their descriptions",
			err: faults.Bad(
				&faults.FieldViolation{Field: "email", Description: "invalid format"},
				&faults.FieldViolation{Field: "name", Description: "required"},
			),
			wantMsg: "invalid format. required",
		},
		{
			scenario: "FailedPrecondition with a violation includes the violation description",
			err:      faults.FailedPrecondition(&faults.PreconditionViolation{Type: "TOS", Description: "not accepted"}),
			wantMsg:  "not accepted",
		},
		{
			scenario: "Aborted with a conflict violation includes the violation description",
			err:      faults.Aborted(&faults.ConflictViolation{Resource: "user:1", Description: "concurrent update"}),
			wantMsg:  "concurrent update",
		},
		{
			scenario: "ResourceExhausted with a quota violation includes the violation description",
			err:      faults.ResourceExhausted(&faults.QuotaViolation{Subject: "project:1", Description: "daily limit exceeded"}),
			wantMsg:  "daily limit exceeded",
		},
		// With* wrapping a cause includes the cause in the message
		{
			scenario: "WithBad wrapping a cause includes the cause in the message",
			err:      faults.WithBad(cause),
			wantMsg:  "bad request: underlying issue",
		},
		{
			scenario: "WithFailedPrecondition wrapping a cause includes the cause in the message",
			err:      faults.WithFailedPrecondition(cause),
			wantMsg:  "precondition failure: underlying issue",
		},
		{
			scenario: "WithAborted wrapping a cause includes the cause in the message",
			err:      faults.WithAborted(cause),
			wantMsg:  "conflict: underlying issue",
		},
		{
			scenario: "WithResourceExhausted wrapping a cause includes the cause in the message",
			err:      faults.WithResourceExhausted(cause),
			wantMsg:  "quota failure: underlying issue",
		},
		{
			scenario: "WithNotFound wrapping a cause includes the cause in the message",
			err:      faults.WithNotFound(cause),
			wantMsg:  "resource not found: underlying issue",
		},
		{
			scenario: "WithPermissionDenied wrapping a cause includes the cause in the message",
			err:      faults.WithPermissionDenied(cause),
			wantMsg:  "permission denied: underlying issue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantMsg {
				t.Errorf("got %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

// TestViolationString validates the human-readable representation of each
// violation type with all fields populated.
func TestViolationString(t *testing.T) {
	t.Run("FieldViolation formats field and description separated by a dash", func(t *testing.T) {
		v := &faults.FieldViolation{Field: "email", Description: "invalid format"}
		want := "email - invalid format"
		if got := v.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("QuotaViolation formats subject and description separated by a dash", func(t *testing.T) {
		v := &faults.QuotaViolation{Subject: "project:123", Description: "daily limit exceeded"}
		want := "project:123 - daily limit exceeded"
		if got := v.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("PreconditionViolation formats type, subject, and description separated by dashes", func(t *testing.T) {
		v := &faults.PreconditionViolation{Type: "TOS", Subject: "google.com", Description: "not accepted"}
		want := "TOS - google.com - not accepted"
		if got := v.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("ConflictViolation formats resource and description separated by a dash", func(t *testing.T) {
		v := &faults.ConflictViolation{Resource: "user:123", Description: "concurrent update"}
		want := "user:123 - concurrent update"
		if got := v.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

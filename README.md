# Faults

- [Faults](#faults)
  - [Benefits](#benefits)
  - [Failure types](#failure-types)
    - [Authentication](#authentication)
    - [Availability](#availability)
    - [Bad](#bad)
    - [Conflict](#conflict)
    - [Missing](#missing)
    - [Permission](#permission)
    - [Pre-condition](#pre-condition)
    - [Quota](#quota)
    - [Unimplemented](#unimplemented)
  - [Litmus test](#litmus-test)
  - [Error chain](#error-chain)
  - [Design](#design)
  - [Disclaimer](#disclaimer)


Package `faults` is an error-handling library that provides simple primitives to represent common failures within a system. Categorising errors simplifies error management and facilitates the propagation of failures across different boundaries, such as HTTP or gRPC.

This package is inspired by `google.golang.org/genproto/googleapis/rpc/errdetails`, which offers a similar set of primitives to describe standard problems encountered in systems.

The primary aim of this library is to establish a common language across different protocols, making it easier to propagate issues. Consider the typical scenario of a service that handles requests via a REST API and communicates with a SQL database. SQL databases have their own set of error codes, which often need to be manually mapped to their HTTP counterparts.

With `faults`, the essence of errors can be abstracted to more effectively communicate the nature of an error to the caller.

## Benefits
Effective error handling is a critical component of robust software systems. It provides several key benefits:

1. **Improved Reliability**: By categorising and managing errors consistently, systems can recover more gracefully from unexpected failures, leading to increased reliability.
1. **Enhanced Debugging**: Clear and consistent error reporting allows developers to diagnose and fix issues more quickly. This reduces downtime and improves the overall stability of the system.
1. **Better User Experience**: When errors are handled well, end users receive meaningful feedback instead of cryptic messages. This leads to a more user-friendly experience, as users can understand what went wrong and, in some cases, how to resolve the issue.
1. **Seamless Cross-Boundary Communication**: In distributed systems, errors often need to be communicated across different services or protocols. A standardised approach to error handling ensures that errors are propagated correctly, maintaining the integrity of the system and reducing the likelihood of miscommunication between components.
1. **Easier Maintenance and Scalability**: As systems grow, maintaining consistent error handling becomes increasingly important. Well-defined error primitives make it easier to extend and scale the system without introducing new points of failure.

By leveraging the `faults` library, developers can create more reliable, maintainable, and user-friendly systems that handle errors in a consistent and predictable manner.

## Failure types

1. Authentication
2. Availability
3. Bad
4. Conflict
5. Missing
6. Permission
7. Pre-condition
8. Quota
9. Unimplemented

### Authentication

This error indicates that the request does not have valid authentication credentials for the operation.

```go
func SensitiveOperation(ctx context.Context) error {
  user, ok := user.FromContext(ctx)
  if !ok {
    return nil, faults.Unauthenticated
  }

  // Perform operation
  return nil
}
```

### Availability

This error describes a temporary state that prevents the request from being fulfilled. The error can contain a delay that advises the caller when it is considered safe to retry.

```go
func OutboundCall() error {
  res, err := http.Get("http://flaky-endpoint")
  if err != nil {
    return err
  }
  if res.StatusCode >= 500 {
    return faults.Unavailable(1 * time.Second)
  }

  // Process response

  return nil
}
```

### Bad

This describes a violation in a client request, usually focusing on the syntactic aspect of the request. For example, a missing field or a name that is too short. It can also involve receiving an unexpected data format. This error is never safe to retry.

```go
violations := []*faults.FieldViolation{
  {
    Field: "firstname",
    Description: "Field required",
  },
  {
    Field: "locality",
    Description: "Field required",
  },
}
err := faults.Bad(violations...)
```

### Conflict

This error indicates that the request conflicts with the current state of the target resource. When this error occurs, the caller typically needs to restart a sequence of operations from the beginning.

```go
func RegisterAccount(email string) error {
  acc, ok := accounts.LoadByEmail(email)
  if ok {
    return faults.Aborted(&faults.ConflictViolation{
      Resource:    fmt.Sprintf("account:%s", email),
      Description: "This email has already been registered",
    })
  }

  // Register account

  return nil
}
```

### Missing

This error means the requested resource was not found. This is the equivalent of a `404` in HTTP.

```go
func LoadAccount(id string) (Account, error) {
  acc, ok := accounts.Load(id)
  if !ok {
    return nil, faults.NotFound
  }

  return acc, nil
}
```

### Permission

This error indicates that the caller does not have permission to execute the specified operation.
It must not be used for rejections caused by exhausting some resource. It must also not be used if the caller cannot be identified.

```go
func SensitiveResource(ctx context.Context) error {
  user, ok := user.FromContext(ctx)
  if !ok {
    return nil, faults.Unauthenticated
  }
  if !user.IsAdmin() {
    return nil, faults.PermissionDenied
  }

  // Perform operation

  return nil
}
```

### Pre-condition

This error indicates that an operation was rejected because the system is not in a state required for the operation's execution. For example, a directory to be deleted may be non-empty, or an rmdir operation may be applied to a non-directory.

```go
func LoginAccount(email, hash string) error {
  acc, ok := accounts.LoadByEmail(email)
  if !ok {
      return faults.FailedPrecondition(&faults.PreconditionViolation{
        Type:        "account",
        Subject:     fmt.Sprintf("account:%s", email),
        Description: "Account does not exist. Please register first",
      })
  }

  // Authenticate..

  return nil
}
```

### Quota

This error describes a failure in a quota check.

For example, if a daily limit is exceeded for the calling project, a service could respond with this error, including details such as the project ID and a description of the exceeded quota limit.

```go
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
  if quotaExceeded() {
    err := faults.ResourceExhausted(&faults.QuotaViolation{
      Subject:     "clientip:<ip address of client>",
      Description: "Daily Limit for read operations exceeded",
    })

    http.Error(w, err.Error(), http.StatusTooManyRequests)
    return
  }

  // Handle the request...
}
```

### Unimplemented

This indicates the operation is not implemented or not supported.

```go
func NewFeature() error {
  return faults.Unimplemented
}
```

## Litmus test

A litmus test that may help a service implementor in deciding between
a pre-condition failure, a conflict, and an unavailability error.

* Use `faults.Unavailable` if the client can retry just the failing call.
* Use `faults.Aborted` if the client should retry at a higher-level (e.g., restarting a read-modify-write sequence).
* Use `faults.FailedPrecondition` if the client should not retry until
      the system state has been explicitly fixed. E.g., if an "rmdir"
      fails because the directory is non-empty, FailedPrecondition
      should be returned since the client should not retry unless
      they have first fixed up the directory by deleting files from it.
* Use `faults.FailedPrecondition` if the client performs conditional
      REST Get/Update/Delete on a resource and the resource on the
      server does not match the condition. E.g., conflicting
      read-modify-write on the same resource.


## Error chain

Go 1.13 introduces the concept of wrapping errors to trace back to the root cause of an issue. An error can be wrapped in this way: `fmt.Errorf("wrapped error: %w", err)`.

With `faults`, errors can also be wrapped easily by calling the prefix `faults.With*`. Once wrapped, an error will be categorised, but the underlying error can still be retrieved.

Example:

```go
_, err := os.Stat(path)
if os.IsNotExist(err) {
  return faults.WithNotFound(err)
}

// Carry on...
```

## Design

This repository was initially hosted at `github.com/deixis/errors` but has since been renamed to `faults`. The original concept was to fully wrap the standard `errors` package, similar to `github.com/pkg/errors`. This approach allowed developers to simply rename their `errors` import and immediately benefit from enhanced functionality while maintaining the familiar API.

However, this approach is no longer ideal. Go has since introduced native support for error wrapping within the standard `errors` package. Moreover, using a custom `errors` package can impede linters and static analysis tools from accurately detecting the misuse of functions like `errors.Is` or `errors.As`.

Consequently, this new version has removed all standard error calls and now focuses exclusively on providing specialised primitives to help developers categorise errors more effectively.

## Disclaimer

The code snippets provided above are intended to demonstrate how to use the different primitives. The examples are purposefully oversimplified and should not be used as-is in production environments.
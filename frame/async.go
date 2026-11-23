package frame

import (
	"context"
	"time"
)

// RedirectResultReceiver carries the value and error delivered by RedirectResult.
type RedirectResultReceiver[T any] struct {
	// Ret is the operation's returned value.
	Ret T
	// Err is the operation's error, if any.
	Err error
}

// AsyncResult represents a single operation. Set its context or timeout before starting it, then choose
// one way to receive the result.
type AsyncResult[T any] interface {
	// SetContext sets the context for the operation and returns the same result handle.
	SetContext(ctx context.Context) AsyncResult[T]
	// SetTimeout sets a deadline relative to now, using the current context as its parent.
	SetTimeout(timeout time.Duration) AsyncResult[T]
	// AsyncGetResult starts the operation and calls callback with the value and error in a separate goroutine.
	AsyncGetResult(callback func(ret T, err error))
	// BlockGetResult starts the operation and waits for its value or an error.
	BlockGetResult() (ret T, err error)
	// RedirectResult sends the value and error to receiver. If block is true, it waits until the send
	// completes; otherwise it starts a goroutine.
	RedirectResult(receiver chan RedirectResultReceiver[T], block bool)
}

// AsyncAPISetHandler installs the implementation of a named API. Choose the form that matches how the
// handler completes.
type AsyncAPISetHandler[inT any, outT any] interface {
	// InstantAPI installs a handler that runs inline and should return promptly.
	InstantAPI(func(in inT) (outT, error))
	// BlockingAPI installs a handler that may block. It runs in a separate goroutine.
	BlockingAPI(func(in inT) (outT, error))
	// CallbackAPI installs a handler that supplies its result by calling setResult rather than returning it.
	CallbackAPI(func(in inT, setResult func(outT, error)))
}

package frame

import (
	"time"
)

// CanSetArg builds a JSON object one key at a time for an operation that does not return a result.
type CanSetArg interface {
	// WithArg adds a named argument and returns the builder. Call Launch after the last argument.
	WithArg(key string, arg any) CanSetArg
	// Launch submits the accumulated arguments without waiting for a reply.
	Launch()
}

// CanSetData accepts the payload for a call, publication or write. The WithJson methods submit it
// immediately; WithArg builds it incrementally.
type CanSetData interface {
	// WithJsonStrData submits jsonStrData as the JSON payload without encoding it again.
	WithJsonStrData(jsonStrData string)
	// WithJsonBytesData submits jsonBytesData as the JSON payload without encoding it again.
	WithJsonBytesData(jsonBytesData []byte)
	// WithJsonableAny encodes the supplied Go value as JSON and submits it.
	WithJsonableAny(jsonableData any)
	// WithArg starts an object payload with a named argument. Call Launch on the returned builder to submit
	// it.
	WithArg(key string, arg any) CanSetArg
}

// CanGetData exposes a JSON payload as raw data or decoded Go values.
type CanGetData interface {
	// RawJsonStr returns the original JSON as a string.
	RawJsonStr() string
	// RawJsonBytes returns the original JSON bytes.
	RawJsonBytes() []byte
	// AsMap decodes the payload as a JSON object.
	AsMap() (map[string]any, error)
	// Bind decodes the entire payload into target, which must be a pointer.
	Bind(target any) error
	// GetValue returns the value for key, or nil when the key is missing.
	GetValue(key string) any
	// TakeValue decodes the named field into a destination pointer and returns the same payload for reading
	// more fields.
	TakeValue(key string, value any) CanGetData
}

// CanSetArgThenResult builds a JSON argument object for a call that returns a reply.
type CanSetArgThenResult interface {
	// WithArg adds a named argument and returns the builder for chaining.
	WithArg(key string, arg any) CanSetArgThenResult
	// Launch returns the operation that submits the accumulated arguments and receives the reply.
	Launch() AsyncResult[CanGetData]
}

// CanSetDataThenResult prepares the payload for a call that returns a JSON reply.
type CanSetDataThenResult interface {
	// WithJsonStrData uses the supplied JSON string as the request and returns a handle for the reply.
	WithJsonStrData(jsonStrData string) AsyncResult[CanGetData]
	// WithJsonBytesData uses the supplied JSON bytes as the request and returns a handle for the reply.
	WithJsonBytesData(jsonBytesData []byte) AsyncResult[CanGetData]
	// WithJsonableAny encodes the value as JSON for the request and returns a handle for the reply.
	WithJsonableAny(jsonableData any) AsyncResult[CanGetData]
	// WithArg starts a named-argument request. Call Launch on the builder to obtain the result handle.
	WithArg(key string, arg any) CanSetArgThenResult
}

// FlexAccess connects components through node APIs, topics and shared values. InProcess methods share Go
// values only within this host process.
type FlexAccess interface {
	FundamentalNode
	KVDataNode
	RolesNode
	TimeLockNode
	// SoftCallOmitResult prepares a call to a named JSON API without a reply.
	SoftCallOmitResult(cmd string) CanSetData
	// SoftCall prepares a call to a named JSON API and returns a builder for its arguments.
	SoftCall(cmd string) CanSetDataThenResult
	// RegSoftAPI registers a named API that accepts a JSON payload and replies with JSON bytes.
	RegSoftAPI(cmd string) AsyncAPISetHandler[CanGetData, []byte]
	// SoftPublish prepares a JSON message for topic. Supply the payload on the returned builder.
	SoftPublish(topic string) CanSetData
	// SoftListen subscribes to JSON messages on topic. The handler must not block.
	SoftListen(topic string, nonBlockingMsgHandleFn func(CanGetData))
	// SoftGet looks up a shared JSON value by key.
	SoftGet(key string) (val CanGetData, found bool)
	// SoftSet returns a builder that writes a JSON value under key.
	SoftSet(key string) CanSetData
	// InProcessGet looks up a Go value shared within this process.
	InProcessGet(key string) (val any, found bool)
	// InProcessSet stores a Go value under key in this process.
	InProcessSet(key string, val any)
	// InProcessCompareAndDelete removes key only if its current value equals old. The comparison value must be
	// comparable.
	InProcessCompareAndDelete(key string, old any) (deleted bool)
	// InProcessCompareAndSwap replaces old with new only if key holds old. The comparison value must be
	// comparable.
	InProcessCompareAndSwap(key string, old any, new any) bool
	// InProcessDelete removes key from the process-local store.
	InProcessDelete(key string)
	// InProcessLoadAndDelete removes key and returns its previous value, if present.
	InProcessLoadAndDelete(key string) (value any, loaded bool)
	// InProcessLoadOrStore returns the existing value and true, or stores value and returns it with false.
	InProcessLoadOrStore(key string, value any) (actual any, loaded bool)
	// InProcessRange visits process-local entries until f returns false.
	InProcessRange(f func(key string, value any) bool)
	// InProcessSwap stores value and returns the previous value and whether it existed.
	InProcessSwap(key string, value any) (previous any, loaded bool)
	// InProcessListen subscribes to a process-local topic. Set newGoroutine to run the handler in a separate
	// goroutine.
	InProcessListen(topic string, onMsg func(any), newGoroutine bool)
	// InProcessPublish delivers a Go value to local subscribers of topic.
	InProcessPublish(topic string, msg any)
	// RegInProcessAPI registers a process-local API whose arguments and results are Go values.
	RegInProcessAPI(apiName string) AsyncAPISetHandler[any, any]
	// InProcessCallAPI returns an operation that calls a process-local API and receives its reply.
	InProcessCallAPI(apiName string, args any) AsyncResult[any]
	// InProcessCallAPIOmitResponse calls a process-local API without collecting a reply.
	InProcessCallAPIOmitResponse(apiName string, args any)
}

// RolesNode advertises and looks up node tags.
type RolesNode interface {
	// SetTags adds the tags to this node.
	SetTags(tags ...string)
	// CheckNetTag reports whether a tag is advertised on the node network.
	CheckNetTag(tag string) bool
	// CheckLocalTag reports whether this node has the given tag.
	CheckLocalTag(tag string) bool
}

// TimeLockNode coordinates named locks with a limited lifetime.
type TimeLockNode interface {
	// TryLock attempts to acquire the named lock for acquireTime and reports whether it succeeded.
	TryLock(name string, acquireTime time.Duration) bool
	// ResetLockTime resets the lifetime of a lock held by this node and reports whether the reset succeeded.
	ResetLockTime(name string, acquireTime time.Duration) bool
	// Unlock releases the named lock held by this node.
	Unlock(name string)
}

// APINode exposes and calls node APIs using byte fragments as arguments and results.
type APINode interface {
	// ExposeAPI registers a named API and returns the handler setter for it.
	ExposeAPI(apiName string) AsyncAPISetHandler[Values, Values]
	// CallOmitResponse calls api with args without collecting a reply.
	CallOmitResponse(api string, args Values)
	// CallWithResponse returns an operation that calls api with args and receives its reply.
	CallWithResponse(api string, args Values) AsyncResult[Values]
}

// TopicNetNode exchanges messages between nodes on named topics.
type TopicNetNode interface {
	// PublishMessage publishes msg to subscribers of topic.
	PublishMessage(topic string, msg Values)
	// ListenMessage subscribes listener to topic. Set newGoroutine to run it in a separate goroutine.
	ListenMessage(topic string, listener func(msg Values), newGoroutine bool)
}

// FundamentalNode combines node API calls with topic messaging.
type FundamentalNode interface {
	APINode
	TopicNetNode
}

// KVDataNode shares byte-fragment values by key across nodes.
type KVDataNode interface {
	// GetValue returns the value stored under key and whether it exists.
	GetValue(key string) (val Values, found bool)
	// SetValue stores val under key in the node network.
	SetValue(key string, val Values)
}

// FlexEnhance is another name for FlexAccess.
type FlexEnhance = FlexAccess

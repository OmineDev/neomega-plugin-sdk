package frame

// StorageAndLogAccess combines file storage with terminal input and logging.
type StorageAndLogAccess interface {
	StorageAndPathAccess
	BackendIO
}

// ComponentFor defines the lifetime of one component instance. F is the host interface passed to Inject.
type ComponentFor[F any] interface {
	// Init reads the configuration for this instance. Use storage for files and logging; game access arrives
	// later in Inject.
	Init(cfg *ComponentConfig, storage StorageAndLogAccess)
	// Inject attaches the host. Register menu entries and game callbacks here.
	Inject(frame F)
	// BeforeActivate runs after injection and before any component is activated. Return an error to stop
	// startup.
	BeforeActivate() error
	// Activate starts the component's work. The host runs it in its own goroutine.
	Activate()
	// Stop releases resources on orderly shutdown. It is not called if the process is killed.
	Stop() error
	// Signal handles a host notification, such as SIGNAL_DATA_CHECKPOINT.
	Signal(int) error
}

// SIGNAL_DATA_CHECKPOINT asks components to save pending changes without stopping.
const SIGNAL_DATA_CHECKPOINT = 0

// BasicComponentFor stores the configuration and host and supplies no-op lifecycle methods. Embed it to
// implement only the hooks a component needs.
type BasicComponentFor[F any] struct {
	// Config is the configuration supplied to Init.
	Config *ComponentConfig
	// Frame is the host interface supplied to Inject.
	Frame F
}

// Init stores cfg in Config.
func (bc *BasicComponentFor[F]) Init(cfg *ComponentConfig, storage StorageAndLogAccess) {
	bc.Config = cfg
}

// Inject stores the host in Frame.
func (bc *BasicComponentFor[F]) Inject(frame F) {
	bc.Frame = frame
}

// BeforeActivate returns nil. Override it to check dependencies before startup.
func (bc *BasicComponentFor[F]) BeforeActivate() error {
	return nil
}

// Activate does nothing by default.
func (bc *BasicComponentFor[F]) Activate() {
}

// Stop returns nil. Override it to release resources on shutdown.
func (bc *BasicComponentFor[F]) Stop() error {
	return nil
}

// Signal ignores host signals and returns nil.
func (bc *BasicComponentFor[F]) Signal(signal int) error {
	return nil
}

// Component is the lifecycle the loader expects from a Go plugin.
type Component = ComponentFor[ExtendOmega]

// BasicComponent is the embeddable base for a component that uses ExtendOmega.
type BasicComponent = BasicComponentFor[ExtendOmega]

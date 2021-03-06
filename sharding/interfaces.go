package sharding

// Service defines items that can be registered into a sharding node's serviceFuncs.
//
// life-cycle management is delegated to the sharding node. The service is allowed to
// initialize itself upon creation, but no goroutines should be spun up outside of the
// Start method.
type Service interface {
	// Start is called after all services have been constructed and the networking
	// layer was also initialized to spawn any goroutines required by the service.
	Start() error

	// Stop terminates all goroutines belonging to the service, blocking until they
	// are all terminated.
	Stop() error
}

// ServiceConstructor defines the callback passed in when registering a service
// to a sharding node.
type ServiceConstructor func() (Service, error)

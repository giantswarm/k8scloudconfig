package id

// Service provides pseudo random hash generation used for ID assignment.
type Service interface {
	// New tries to create a new object ID using the configured ID type. The
	// returned error might be caused by timeouts reached during the ID creation.
	New() (string, error)
	// WithType tries to create a new object ID using the given ID type. The
	// returned error might be caused by timeouts reached during the ID creation.
	WithType(idType int) (string, error)
}

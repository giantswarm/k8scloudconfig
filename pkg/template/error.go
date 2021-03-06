package template

import "github.com/giantswarm/microerror"

var componentNotFoundError = &microerror.Error{
	Kind: "componentNotFound",
}

// IsComponentNotFound asserts componentNotFoundError.
func IsComponentNotFound(err error) bool {
	return microerror.Cause(err) == componentNotFoundError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var retrieveRuntimeError = &microerror.Error{
	Kind: "retrieveRuntimeError",
}

// IsRetrieveRuntimeError asserts retrieveRuntimeError.
func IsRetrieveRuntimeError(err error) bool {
	return microerror.Cause(err) == retrieveRuntimeError
}

var validationError = &microerror.Error{
	Kind: "validationError",
}

// IsValidationError asserts validationError.
func IsValidationError(err error) bool {
	return microerror.Cause(err) == validationError
}

package v_4_1_1

import "github.com/giantswarm/microerror"

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

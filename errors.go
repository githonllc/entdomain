package entdomain

import "errors"

// Sentinel errors returned by generated repositories. Use errors.Is() or the
// provided Is* helpers to check error types without string matching.
var (
	// ErrNotFound indicates the requested entity does not exist.
	ErrNotFound = errors.New("entity not found")

	// ErrAlreadyExists indicates a uniqueness constraint violation.
	ErrAlreadyExists = errors.New("entity already exists")

	// ErrValidation indicates the input failed validation.
	ErrValidation = errors.New("validation failed")
)

// IsNotFound reports whether err (or any error in its chain) is ErrNotFound.
func IsNotFound(err error) bool { return errors.Is(err, ErrNotFound) }

// IsAlreadyExists reports whether err (or any error in its chain) is ErrAlreadyExists.
func IsAlreadyExists(err error) bool { return errors.Is(err, ErrAlreadyExists) }

// IsValidation reports whether err (or any error in its chain) is ErrValidation.
func IsValidation(err error) bool { return errors.Is(err, ErrValidation) }

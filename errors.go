package parallel

import "errors"

// ErrInvalidLimit is returned when a limit parameter is less than or equal to 0.
var ErrInvalidLimit = errors.New("limit must be greater than 0")

package dbutil

import (
	"strings"
)

// IsDuplicate checks whether error returned from db is as a result of violating unique constraint
func IsDuplicate(err error) bool {
	// Checks whether the error was due to similar account existing
	return strings.Contains(
		strings.ToLower(err.Error()), "duplicate entry",
	)
}

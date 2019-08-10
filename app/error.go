package app

import (
	"fmt"

	"github.com/pkg/errors"
)

func ServerError(err error, msg string) string {
	// errors.Wrap doesn't handle nil errors. To avoid nil pointer in error
	// message we create empty error here when error is nil
	if err == nil {
		err = errors.New("")
	}
	return fmt.Sprintf(`{"error":{"message":["%s"]}}`, errors.Wrap(err, msg))
}

/*	gretry package can be used for retrying for specific actions
	such as retry password.
	retry reset password.
	retry reset masterpassword
	and more
*/

package gretry

import "errors"

// MAXIMUMALLOWEDATTEMPTS : maximum allowed attempts for retrying master key reset
var MAXIMUMALLOWEDATTEMPTS = 5

var errMaxAttemptReached = errors.New("maximum allowed retry reached")

// RetryFunction typedef
type RetryFunction func(attempts int) error

// Retry resets
func Retry(retryFunc RetryFunction) error {

	for ret := 1; ret < MAXIMUMALLOWEDATTEMPTS; ret++ {
		err := retryFunc(ret)
		if err == nil {
			return nil
		}
	}

	return errMaxAttemptReached
}

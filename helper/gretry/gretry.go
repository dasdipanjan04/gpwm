/*	gretry package can be used for retrying for specific actions
	such as retry password.
	retry reset password.
	retry reset masterpassword
	and more
*/

package gretry

import "errors"

var MAXIMUMALLOWEDATTEMPTS = 5

var maxAttemptReachedError = errors.New("Maximum allowed retry reached.")

type RetryFunction func(attempts int) error

func Retry(retryFunc RetryFunction) error {

	for ret := 1; ret < MAXIMUMALLOWEDATTEMPTS; ret++ {
		err := retryFunc(ret)
		if err == nil {
			return nil
		}
	}

	return maxAttemptReachedError
}

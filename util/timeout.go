package util

import (
	"time"
	"errors"
)

// return true when completed, false if error or timeout
// note: the action itself must not hang
func RetryEveryUntil(action func()(bool, error), retry, limit time.Duration) (bool, error) {
	tick := time.Tick(retry)
	timeout := time.After(limit)
	for {
		select {
		case <-timeout:
			return false, errors.New("timeout")
		case <-tick:
			ok, err := action()
			if err != nil {
				return ok, err
			} else if ok {
				return ok, nil
			}
		}
	}
}

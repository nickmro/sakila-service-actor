// Package redis contains functions for the actor service cache.
package redis

import (
	"crypto/sha1"
	"fmt"
	"time"
)

const (
	// DefaultTTL is the default Redis cache TTL.
	DefaultTTL = time.Minute * 5
)

func hashedKey(key string) string {
	h := sha1.New()

	if _, err := h.Write([]byte(key)); err != nil {
		return ""
	}

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

package constants

import "time"

const (
	KEY_USER_ID  = "USER_%d"
	KEY_USERNAME = "USER_%s"
	REDIS_EXPIRE = 15 * time.Minute
	// REDIS_EXPIRE = 1 * time.Nanosecond
)

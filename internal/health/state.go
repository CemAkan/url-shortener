package health

import "sync/atomic"

var (
	// 0 = down, 1 = up
	dbStatus    atomic.Int32
	redisStatus atomic.Int32
	emailStatus atomic.Int32
)

// SetDBStatus sets database status atomic safe
func SetDBStatus(up bool) {
	if up {
		dbStatus.Store(1)
	} else {
		dbStatus.Store(0)
	}
}

// SetRedisStatus sets redis status atomic safe
func SetRedisStatus(up bool) {
	if up {
		redisStatus.Store(1)
	} else {
		redisStatus.Store(0)
	}
}

// SetEmailStatus sets mail service status atomic safe
func SetEmailStatus(up bool) {
	if up {
		emailStatus.Store(1)
	} else {
		emailStatus.Store(0)
	}
}

// GetDBStatus returns database status from atomic dbStatus status var
func GetDBStatus() bool {
	return dbStatus.Load() == 1
}

// GetRedisStatus returns redis status from atomic redisStatus status var
func GetRedisStatus() bool {
	return redisStatus.Load() == 1
}

// GetEmailStatus returns email service status from atomic emailStatus status var
func GetEmailStatus() bool {
	return emailStatus.Load() == 1
}

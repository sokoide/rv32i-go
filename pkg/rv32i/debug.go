//go:build debug

package rv32i

import log "github.com/sirupsen/logrus"

func trace(args ...interface{}) {
	log.Trace(args...)
}

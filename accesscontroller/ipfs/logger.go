package ipfs

import "go.uber.org/zap"

func logger() *zap.Logger {
	return zap.L().Named("orbitdb.accesscontroller.ipfs")
}

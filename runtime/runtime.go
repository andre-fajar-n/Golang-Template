package runtime

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Runtime struct {
	Db     *gorm.DB
	Cfg    config
	Logger zerolog.Logger
}

func NewRuntime() *Runtime {
	rt := new(Runtime)

	rt = rt.logger()

	rt = rt.config()

	rt = rt.db()

	return rt
}

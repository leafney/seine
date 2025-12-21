package memcx

import (
	rmemc "github.com/leafney/rose/memc"
	"github.com/leafney/seine/pkg/xlogx"
)

type MemcSvc struct {
	*rmemc.Memc
}

func NewMemcSvc(log *xlogx.XLogSvc) *MemcSvc {

	mm := rmemc.NewMemc()
	log.Infoln("[Memc] Load successful")

	return &MemcSvc{mm}
}

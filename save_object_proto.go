package storage

import (
	"context"

	"github.com/myfantasy/mfs"
)

type SaveObjectProto struct {
	Mx     mfs.PMutex `json:"-"`
	MxStor mfs.PMutex `json:"-"`

	Rv     int64 `json:"_rv"`
	StorRv int64 `json:"_stor_rv"`
}

func (sp *SaveObjectProto) StorLock(ctx context.Context) bool {
	return sp.MxStor.TryLock(ctx)
}
func (sp *SaveObjectProto) StorUnlock() {
	sp.MxStor.Unlock()
}

func (sp *SaveObjectProto) DataRLock(ctx context.Context) bool {
	return sp.Mx.RTryLock(ctx)
}
func (sp *SaveObjectProto) DataRUnlock() {
	sp.Mx.RUnlock()
}
func (sp *SaveObjectProto) DataLock(ctx context.Context) bool {
	return sp.Mx.TryLock(ctx)
}
func (sp *SaveObjectProto) DataUnlock() {
	sp.Mx.Unlock()
}

func (sp *SaveObjectProto) RvGetLast() int64 {
	return sp.Rv
}
func (sp *SaveObjectProto) RvGetStor() int64 {
	return sp.StorRv
}
func (sp *SaveObjectProto) RvSetStor(rv int64) {
	sp.StorRv = rv
}

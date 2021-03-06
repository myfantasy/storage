package storage

import (
	"context"
	"time"

	"github.com/myfantasy/mfs"
	"github.com/myfantasy/mft"
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

type SaveProto struct {
	SaveObjectProto

	RVG *mft.G `json:"_rvg,omitempty"`

	SaveToStorageValue  Storage `json:"-"`
	SaveToFileNameValue string  `json:"-"`

	SaveToContextValue func() (context.Context, context.CancelFunc) `json:"-"`

	SaveToContextDuration time.Duration `json:"-"`
}

func (sp *SaveProto) SaveToStorage() Storage {
	return sp.SaveToStorageValue
}
func (sp *SaveProto) SaveToFileName() string {
	return sp.SaveToFileNameValue
}
func (sp *SaveProto) SaveToContext() (context.Context, context.CancelFunc) {
	if sp.SaveToContextValue != nil {
		return sp.SaveToContextValue()
	}
	if sp.SaveToContextDuration > 0 {
		return context.WithTimeout(context.Background(), sp.SaveToContextDuration)
	}

	return context.Background(), func() {}
}

func (sp *SaveProto) SetNextPartRv() {
	sp.Rv = sp.RVG.RvGetPartOrGlobal()
}

package storage

import (
	"context"

	"github.com/myfantasy/mft"
)

type Storable interface {
	// StorLock - Lock object for storage
	StorLock(ctx context.Context) bool
	// StorUnlock - Unlock object for storage
	StorUnlock()

	// DataRLock - Lock object for read
	DataRLock(ctx context.Context) bool
	// DataRUnlock - Unlock object for read
	DataRUnlock()
	// DataLock - Lock object
	DataLock(ctx context.Context) bool
	// DataUnlock - Unlock object
	DataUnlock()

	// RvGetLast - get last object version
	RvGetLast() int64
	// RvGetLast - get last stored version
	RvGetStor() int64
	// RvSetStor - set last stored version
	RvSetStor(rv int64)

	// ToBytes - Object to bytes for save
	ToBytes() (data []byte, err *mft.Error)
	// FromBytes - Object from bytes for load
	FromBytes(data []byte) (err *mft.Error)
}

type DoFunc func() *mft.Error
type CheckFunc func() (ok bool, err *mft.Error)

func Save(ctx context.Context, s Storage, fileName string, v Storable) (err *mft.Error) {
	return SaveExtend(ctx, s, fileName, v, nil, nil, nil)
}
func SaveExtend(ctx context.Context, s Storage, fileName string, v Storable, doBeforeGetData DoFunc, doBeforeSave DoFunc, doAfterSave DoFunc) (err *mft.Error) {
	if s == nil {
		return GenerateError(10001100)
	}
	if !v.StorLock(ctx) {
		return GenerateError(10001101)
	}
	defer v.StorUnlock()

	if doBeforeGetData != nil {
		err = doBeforeGetData()
		if err != nil {
			return GenerateErrorE(10001104, err)
		}
	}

	if !v.DataRLock(ctx) {
		return GenerateError(10001102)
	}

	lastRv := v.RvGetLast()
	storRv := v.RvGetStor()

	if lastRv == storRv {
		v.DataRUnlock()
		return nil
	}

	body, err := v.ToBytes()
	if err != nil {
		v.DataRUnlock()
		return GenerateErrorE(10001103, err)
	}

	v.DataRUnlock()

	if doBeforeSave != nil {
		err = doBeforeSave()
		if err != nil {
			return GenerateErrorE(10001105, err)
		}
	}

	err = s.Save(ctx, fileName, body)
	if err != nil {
		return GenerateErrorE(10001107, err, fileName)
	}

	if doAfterSave != nil {
		err = doAfterSave()
		if err != nil {
			return GenerateErrorE(10001106, err)
		}
	}

	v.RvSetStor(storRv)

	return nil
}

// Load
func Load(ctx context.Context, s Storage, fileName string, v Storable) (err *mft.Error) {
	ok, err := LoadIfExists(ctx, s, fileName, v, nil)
	if err != nil {
		return err
	}
	if !ok {
		GenerateError(10001121, fileName)
	}
	return nil
}
func LoadIfExists(ctx context.Context, s Storage, fileName string, v Storable, doFill CheckFunc) (ok bool, err *mft.Error) {
	if s == nil {
		return false, GenerateError(10001120)
	}
	if !v.StorLock(ctx) {
		return false, GenerateError(10001122)
	}
	defer v.StorUnlock()

	ok, err = s.Exists(ctx, fileName)
	if err != nil {
		return false, GenerateErrorE(10001123, err, fileName)
	}
	if !ok {
		return false, nil
	}

	body, err := s.Get(ctx, fileName)
	if err != nil {
		return false, GenerateErrorE(10001124, err, fileName)
	}

	if !v.DataLock(ctx) {
		return false, GenerateError(10001126)
	}

	ok, err = doFill()
	if err != nil {
		v.DataUnlock()
		return false, GenerateErrorE(10001127, err)
	}
	if !ok {
		v.DataUnlock()
		return true, nil
	}

	err = v.FromBytes(body)
	if err != nil {
		v.DataUnlock()
		return false, GenerateErrorE(10001125, err, fileName)
	}
	v.RvSetStor(v.RvGetLast())
	v.DataUnlock()

	return true, nil
}

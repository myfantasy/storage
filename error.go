package storage

import (
	"fmt"

	"github.com/myfantasy/mft"
)

// Errors codes and description
var Errors map[int]string = map[int]string{
	10000000: "MapSorage: not found name: `%v`",
	10000001: "MapSorage: name exists: `%v`",
	10000002: "File: error: `%v`",
	10000003: "Mkdir error: path: `%v`",

	10001000: "Cluster.Create: Lock mutex fail wait",
	10001001: "Cluster.Create: storage type `%v` is not exists",
	10001002: "Cluster.Create: mount `%v` is not exists",

	10001100: "SaveExtend: Storage is nil",
	10001101: "SaveExtend: Storable StorLock fail",
	10001102: "SaveExtend: Storable DataRLock fail",
	10001103: "SaveExtend: Storable ToBytes fail",
	10001104: "SaveExtend: doBeforeGetData fail",
	10001105: "SaveExtend: doBeforeSave fail",
	10001106: "SaveExtend: doAfterSave fail",
	10001107: "SaveExtend: Storage Save to storage `%v` fail",

	10001120: "LoadIfExists: Storage is nil",
	10001121: "Load: File `%v` does not exists",
	10001122: "LoadIfExists: Storable StorLock fail",
	10001123: "LoadIfExists: Storage check file `%v` exists fail",
	10001124: "LoadIfExists: Storage get file `%v` fail",
	10001125: "LoadIfExists: Storable FromBytes from file `%v` fail",
	10001126: "LoadIfExists: Storable DataLock fail",
	10001127: "LoadIfExists: doFill Check fail",
}

// GenerateError -
func GenerateError(key int, a ...interface{}) *mft.Error {
	if text, ok := Errors[key]; ok {
		return mft.ErrorCS(key, fmt.Sprintf(text, a...))
	}
	panic(fmt.Sprintf("storage.GenerateError, error not found code:%v", key))
}

// GenerateErrorE -
func GenerateErrorE(key int, err error, a ...interface{}) *mft.Error {
	if text, ok := Errors[key]; ok {
		return mft.ErrorCSE(key, fmt.Sprintf(text, a...), err)
	}
	panic(fmt.Sprintf("storage.GenerateErrorE, error not found code:%v error:%v", key, err))
}

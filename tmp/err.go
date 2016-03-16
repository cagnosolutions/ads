package tmp

import "errors"

var (
	StoreExists    = errors.New("db: store already exists")
	StoreNotFound  = errors.New("db: store was not found")
	RecordTooLarge = errors.New("data: record exceeded maximum size of 4KB")
	NonPtrValue    = errors.New("data: received non pointer value, expected pointer")
	RecordMaximum  = errors.New("data: maximum record count has been reached")
)

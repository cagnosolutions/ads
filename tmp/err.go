package tmp

import "errors"

var (
	StoreExists   = errors.New("db: the store already exists")
	StoreNotFound = errors.New("db: the store was not found")
	NonPtrValue   = errors.New("store: received non pointer value, expected pointer")
	StoreIsFull   = errors.New("store: maximum size for this store has been reached; store is full")
	DocExists     = errors.New("store: document/record already exists")
	DocTooLarge   = errors.New("store: document/record size exceeded maximum (4KB)")
)

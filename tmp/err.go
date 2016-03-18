package tmp

import "errors"

var (
	ErrUnknown    = errors.New("db: oh man... an unknown error has occurred")
	StoreExists   = errors.New("db: the store already exists")
	StoreNotFound = errors.New("db: the store was not found")
	NonPtrValue   = errors.New("store: received non pointer value, expected pointer")
	StoreIsFull   = errors.New("store: maximum size for this store has been reached; store is full")
	DocExists     = errors.New("store: document/record already exists")
	DocNotFound   = errors.New("store: document/record was not found")
	DocTooLarge   = errors.New("store: document/record size exceeded maximum (4KB)")
)

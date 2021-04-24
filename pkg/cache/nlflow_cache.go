package cache

const KeyNotFound = CacheError("KeyNotFound")

type CacheError string

func (e CacheError) Error() string { return string(e) }

type NLFlowCache interface {
	Read(k string) (string, error)
	Write(k string, v string) error
	Close() error
}

package cache

type CacheType int

const (
	DAPR CacheType = iota
	MOCK
)

func NewNLFlowCache(cacheType CacheType) NLFlowCache {
	switch cacheType {
	case DAPR:
		return NewNLFlowCacheDapr()
	case MOCK:
		return NewNLFlowCacheMock()
	default:
		return nil
	}
}

package etlcache

import (
	"github.com/advancedlogic/nlflow-go-sdk/cache"
)

type nlflowCacheETLMock struct {
	nLFlowCache cache.NLFlowCache
}

func NewNLFlowCacheETLMock(nLFlowCache cache.NLFlowCache) NLFlowCacheETL {
	return NLFlowCacheETL{
		NLFlowCacheETL: &nlflowCacheETLMock{
			nLFlowCache: nLFlowCache,
		},
	}
}

func (h *nlflowCacheETLMock) ResolveMetaIngress(metaReq MetadataRequest) (string, error) {
	return "", nil
}

func (h *nlflowCacheETLMock) ResolveMetaResponse(metaReq MetadataRequest, strRes string) (MetadataResponse, error) {
	return MetadataResponse{}, nil
}

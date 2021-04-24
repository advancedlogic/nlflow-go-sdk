package etlcache

import (
	"github.com/advancedlogic/nlflow-go-sdk/pkg/cache"
)

func NewNLFlowCacheETL(cacheType string, metadataMappingType string) NLFlowCacheETL {
	if cacheType == "DAPR" && metadataMappingType == "JSONPATH_PAIRLIST" {
		return NewNLFlowCacheETLJsonPath(cache.NewNLFlowCacheDapr())
	} else if cacheType == "MOCK" && metadataMappingType == "MOCK" {
		return NewNLFlowCacheETLMock(cache.NewNLFlowCacheMock())
	}
	return NLFlowCacheETL{}
}

func BuildCacheKey(workflowId string, serviceId string, docId string) string {
	return workflowId + "|" + serviceId + "|" + docId
}

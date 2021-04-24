package etlcache

type NLFlowCacheETL struct {
	NLFlowCacheETL INLFlowCacheETL
}

type INLFlowCacheETL interface {
	ResolveMetaIngress(metaRequest MetadataRequest) (string, error)
	ResolveMetaResponse(metaRequest MetadataRequest, respStr string) (MetadataResponse, error)
}

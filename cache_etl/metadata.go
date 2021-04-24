package etlcache

type MetadataMappingType int

const (
	JSONPATH_PAIRLIST = iota
	MOCK
)

type MetadataMapping struct {
	MappingType string `json:"mappingType"`
	MappingDsl  string `json:"mappingDsl"` // domain specific language
	CacheType   string `json:"cacheType"`
}

type MetadataRequest struct {
	Version             string          `json:"version"`
	JobId               string          `json:"jobId"`
	Timestamp           int64           `json:"timestamp,omitempty"`
	SourceServiceIds    []string        `json:"sourceServiceIds"`
	MicrofrontendParams string          `json:"microfrontendParams"`
	DocId               string          `json:"docId"`
	WorkflowId          string          `json:"workflowId"`
	Mapping             MetadataMapping `json:"mapping"`
	DestServiceId       string          `json:"destServiceId"`
}

type MetadataResponse struct {
	Version   string `json:"version"`
	JobId     string `json:"jobId"`
	Timestamp int64  `json:"timestamp"`
	DocId     string `json:"docId"`
	ServiceId string `json:"serviceId"`
}

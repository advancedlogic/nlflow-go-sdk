package etlcache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/advancedlogic/nlflow-go-sdk/pkg/cache"
	kazaam "gopkg.in/qntfy/kazaam.v3"
)

type nlflowCacheETLJsonPath struct {
	nLFlowCache cache.NLFlowCache
}

func NewNLFlowCacheETLJsonPath(nLFlowCache cache.NLFlowCache) NLFlowCacheETL {
	return NLFlowCacheETL{
		NLFlowCacheETL: &nlflowCacheETLJsonPath{
			nLFlowCache: nLFlowCache,
		},
	}
}

func (h *nlflowCacheETLJsonPath) ResolveMetaIngress(metaReq MetadataRequest) (string, error) {

	ferrObj := make(map[string]interface{})
	serviceObj := make(map[string]interface{})

	for _, sourceServiceId := range metaReq.SourceServiceIds {
		cacheKey := BuildCacheKey(metaReq.WorkflowId, sourceServiceId, metaReq.DocId)
		valueInCache, err := h.nLFlowCache.Read(cacheKey)
		if err != nil && err != cache.KeyNotFound {
			return "", err
		}
		if len(valueInCache) > 0 {
			// found in cache
			var result map[string]interface{}
			err := json.Unmarshal([]byte(valueInCache), &result)
			if err != nil {
				return "", err
			}
			serviceObj[sourceServiceId] = result
		} else {
			serviceObj[sourceServiceId] = nil
		}
	}

	ferrObj[metaReq.DocId] = serviceObj
	ferrObjAsByteArray, err := json.MarshalIndent(ferrObj, "", "    ")
	if err != nil {
		return "", err
	}
	ferrObjAsString := string(ferrObjAsByteArray)

	//fmt.Println(ferrObjAsString)

	//
	//
	//
	//

	//fmt.Println(metaReq.Mapping.MappingDsl)

	// applico mappingDSL a sourceJsonString
	payloadMS := make(map[string]interface{})
	var m interface{}
	err = json.Unmarshal([]byte(metaReq.Mapping.MappingDsl), &m)
	if err != nil {
		return "", err
	}
	mapping := m.(map[string]interface{})
	pathMappings := mapping["pathMappings"]
	for _, mp := range pathMappings.([]interface{}) {
		a := mp.(map[string]interface{})
		src := metaReq.DocId + "." + a["source"].(string)
		dst := a["target"].(string)
		//fmt.Println(src)
		//fmt.Println(dst)

		k, _ := kazaam.NewKazaam(fmt.Sprintf(`[ { "operation": "shift", "spec": { "%s": "%s" } } ]`, dst, src))
		kazaamOut, err := k.TransformJSONStringToString(ferrObjAsString)
		if err != nil {
			return "", err
		}
		//fmt.Println(kazaamOut)

		var kazaamObj map[string]interface{}
		err = json.Unmarshal([]byte(kazaamOut), &kazaamObj)
		if err != nil {
			return "", err
		}
		mergeKeys(payloadMS, kazaamObj)
	}

	payloadMSAsByteArray, err := json.MarshalIndent(payloadMS, "", "    ")
	if err != nil {
		return "", err
	}
	payloadMSAsString := string(payloadMSAsByteArray)

	//fmt.Println(payloadMSAsString)

	return payloadMSAsString, nil
}

// Given two maps, recursively merge right into left, NEVER replacing any key that already exists in left
func mergeKeys(left, right map[string]interface{}) map[string]interface{} {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			//then we don't want to replace it - recurse
			left[key] = mergeKeys(leftVal.(map[string]interface{}), rightVal.(map[string]interface{}))
		} else {
			// key not in left so we can just shove it in
			left[key] = rightVal
		}
	}
	return left
}

func (h *nlflowCacheETLJsonPath) ResolveMetaResponse(metaReq MetadataRequest, strRes string) (MetadataResponse, error) {
	err := h.nLFlowCache.Write(BuildCacheKey(metaReq.WorkflowId, metaReq.DestServiceId, metaReq.DocId), strRes)
	if err != nil {
		return MetadataResponse{}, err
	}
	return MetadataResponse{
		DocId:     metaReq.DocId,
		JobId:     metaReq.JobId,
		ServiceId: metaReq.DestServiceId,
		Timestamp: time.Now().UnixNano(),
		Version:   metaReq.Version,
	}, nil
}

package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	etlcache "github.com/advancedlogic/nlflow-go-sdk/cache_etl"
	echo "github.com/labstack/echo/v4"
)

type nlflowHandler struct {
	handler     Handler
	metaHandler Handler

	Error error
}

func (h *nlflowHandler) GetHandler() Handler {
	return h.handler
}

func (h *nlflowHandler) WithEndpoint(endpoint string) *nlflowHandler {
	h.handler.endpoint = endpoint
	return h
}

func (h *nlflowHandler) WithTransportIn(transportIn TransportIn) *nlflowHandler {
	h.handler.transportIn = transportIn
	return h
}

func (h *nlflowHandler) WithTransportOut(transportOut TransportOut) *nlflowHandler {
	h.handler.transportOut = transportOut
	return h
}

func (h *nlflowHandler) WithLogic(logic Logic) *nlflowHandler {
	h.handler.logic = logic
	return h
}

/////////////////////////////////////////////////////////////

type metadataModelIn struct {
	ctx     echo.Context
	metaReq etlcache.MetadataRequest
}

func (metadataModelIn) Version() string {
	return "v1"
}

/////////////////////////////////////////////////////////////

type metadataModelOut struct {
	metaRes etlcache.MetadataResponse
}

func (metadataModelOut) Version() string {
	return "v1"
}

/////////////////////////////////////////////////////////////

type metadataTransportIn struct {
}

func (mdT *metadataTransportIn) Process(c echo.Context) (Model, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}
	defer c.Request().Body.Close()

	metaModelIn := metadataModelIn{
		ctx: c,
	}
	err = json.Unmarshal(body, &metaModelIn.metaReq)
	if err != nil {
		return nil, err
	}

	return metaModelIn, nil
}

/////////////////////////////////////////////////////////////

type metadataLogic struct {
	handler Handler
}

func (mdL *metadataLogic) Process(m Model) (Model, error) {
	metaIn := m.(metadataModelIn)

	//accedo alla cache per recuperare i dati della request
	cache := etlcache.NewNLFlowCacheETL(
		metaIn.metaReq.Mapping.CacheType,
		metaIn.metaReq.Mapping.MappingType)

	// costruisco la request input da inviare al servizio
	metaIngress, err := cache.NLFlowCacheETL.ResolveMetaIngress(metaIn.metaReq)
	if err != nil {
		return nil, err
	}

	// cambio il body della request mettendoci l'input atteso dalla logica
	metaIn.ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer([]byte(metaIngress)))
	//                          ioutil.NopCloser(bytes.NewReader([]byte(metaIngress)))
	//                          ioutil.NopCloser(bytes.NewBuffer([]byte(metaIngress)))

	// convert http request body to business login expected input model
	inputModel, err := mdL.handler.GetTransportIn().Process(metaIn.ctx)
	if err != nil {
		return nil, err
	}

	// process business logic
	outputModel, err := mdL.handler.GetLogic().Process(inputModel)
	if err != nil {
		return nil, err
	}

	boutput, err := json.Marshal(outputModel)
	if err != nil {
		return nil, err
	}

	metaRes, err := cache.NLFlowCacheETL.ResolveMetaResponse(metaIn.metaReq, string(boutput))
	if err != nil {
		return nil, err
	}
	metadataModelOut := metadataModelOut{
		metaRes: metaRes,
	}

	return metadataModelOut, nil
}

/////////////////////////////////////////////////////////////

type metadataTransportOut struct {
}

func (*metadataTransportOut) Process(c echo.Context, outputModel Model) error {
	metadataModelOut := outputModel.(metadataModelOut)
	// convert metadata output to json and send back
	boutput, err := json.Marshal(metadataModelOut.metaRes)
	if err != nil {
		return err
	}
	c.String(http.StatusOK, string(boutput))
	return nil
}

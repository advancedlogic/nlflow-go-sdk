package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/alecthomas/jsonschema"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	uuid "github.com/satori/go.uuid"
	//"go.uber.org/zap"
)

var logo = `
   ///         ____  _____  _____     ________  __                      
 (o)(o)       |_   \|_   _||_   _|   |_   __  |[  |                     
 (°__°)         |   \ | |    | |       | |_ \_| | |  .--.   _   _   __  
  (__)          | |\ \| |    | |   _   |  _|    | |/ .''\ \[ \ [ \ [  ] 
  (__)         _| |_\   |_  _| |__/ | _| |_     | || \__. | \ \/\ \/ /  
 “(__)(_)(_)  |_____|\____||________||_____|   [___]'.__.'   \__/\__/   
`

type Service struct {
	id        string
	timestamp int64
	//logger    *zap.Logger

	host            string
	port            int
	grpcEnabled     bool
	debuggerEnabled bool

	manifest           Manifest
	middlewares        []Middleware
	handlers           []Handler
	nlflowHandler      nlflowHandler
	staticFilePath     string
	staticFileEndpoint string
	rootEndpoint       string
	info               Info

	e *echo.Echo

	Error error
}

func NewService() *Service {
	//	defaultLogger, err := zap.NewDevelopment()
	//	defer defaultLogger.Sync()

	return &Service{
		id:        uuid.NewV4().String(),
		timestamp: time.Now().UnixNano(),
		host:      "127.0.0.1",
		port:      18080,
		//logger:      defaultLogger,
		handlers:    make([]Handler, 0),
		middlewares: make([]Middleware, 0),
		grpcEnabled: false,
		//Error:       err,
	}
}

func (service *Service) WithHost(host string) *Service {
	service.host = host
	return service
}

func (service *Service) WithPort(port int) *Service {
	service.port = port
	return service
}

func (service *Service) WithRootEndpoint(rootEndpoint string) *Service {
	service.rootEndpoint = rootEndpoint
	return service
}

func (service *Service) WithManifest(manifest Manifest) *Service {
	service.manifest = manifest
	return service
}

//func (service *Service) WithLogger(logger *zap.Logger) *Service {
//	service.logger = logger
//	return service
//}

func (service *Service) WithGRPC() *Service {
	service.grpcEnabled = true
	return service
}

func (service *Service) WithDebugger() *Service {
	service.debuggerEnabled = true
	return service
}

func (service *Service) WithStaticFiles(endpoint, path string) *Service {
	service.staticFileEndpoint = endpoint
	service.staticFilePath = path
	return service
}

func (service *Service) WithConfiguration() *Service {
	// TODO implementation
	return service
}

func (service *Service) WithEnvVariables() *Service {
	// TODO implementation
	return service
}

func (service *Service) WithBanner(path string) *Service {
	if path != "" {
		blogo, err := ioutil.ReadFile("./banner.txt")
		if err == nil {
			logo = string(blogo)
		}
	}
	println(logo)
	return service
}

func (service *Service) AddMiddleware(middleware Middleware) *Service {
	service.middlewares = append(service.middlewares, middleware)
	return service
}

func (service *Service) AddHandler(handler Handler) *Service {
	service.handlers = append(service.handlers, handler)
	return service
}

func (service *Service) SetNLFlowHandler(handler Handler) *Service {
	service.nlflowHandler = nlflowHandler{
		handler: handler,
		metaHandler: Handler{
			transportIn:  &metadataTransportIn{},
			transportOut: &metadataTransportOut{},
			logic: &metadataLogic{
				handler: handler,
			},
		},

		Error: handler.Error,
	}
	return service
}

func (service *Service) init() {
	e := echo.New()
	e.HideBanner = true

	if service.debuggerEnabled {
		e.Debug = true
	}

	if service.staticFileEndpoint != "" && service.staticFilePath != "" {
		e.Static(service.staticFileEndpoint, service.staticFilePath)
	}

	//if service.logger != nil { // TODO fixme
	e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	//}

	//process middlewares
	for _, middleware := range service.middlewares {
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				err := middleware.Process(c)
				if err != nil {
					return err
				}
				return next(c)
			}
		})
	}

	functionHandler := func(handler Handler) func(c echo.Context) error {
		return func(c echo.Context) error {

			// convert http request body to business login expected input model
			inputModel, err := handler.GetTransportIn().Process(c)
			if err != nil {
				return err
			}

			// process business logic
			outputModel, err := handler.GetLogic().Process(inputModel)
			if err != nil {
				return err
			}

			// convert business login output model to http response
			err = handler.GetTransportOut().Process(c, outputModel)

			return err
		}
	}

	functionNLFlowHandler := func(nlflowHandler nlflowHandler) func(c echo.Context) error {
		return func(c echo.Context) error {

			if c.QueryParam("type") == "meta" {
				// convert http request body to metadata request
				inputModel, err := nlflowHandler.metaHandler.transportIn.Process(c)
				if err != nil {
					return err
				}
				// process metadata logic
				outputModel, err := nlflowHandler.metaHandler.logic.Process(inputModel)
				if err != nil {
					return err
				}
				// convert metadata logic output to json and send back
				err = nlflowHandler.metaHandler.transportOut.Process(c, outputModel)
			} else {
				// convert http request body to business login expected input model
				inputModel, err := nlflowHandler.handler.GetTransportIn().Process(c)
				if err != nil {
					return err
				}

				// process business logic
				outputModel, err := nlflowHandler.handler.GetLogic().Process(inputModel)
				if err != nil {
					return err
				}

				// convert business login output model to http response
				err = nlflowHandler.handler.GetTransportOut().Process(c, outputModel)
			}

			return nil
		}
	}

	//process handlers
	for _, handler := range service.handlers {
		switch handler.GetMethod() {
		case http.MethodGet:
			e.GET(service.rootEndpoint+handler.GetEndpoint(), functionHandler(handler))
		case http.MethodPost:
			e.POST(service.rootEndpoint+handler.GetEndpoint(), functionHandler(handler))
		case http.MethodPut:
			e.PUT(service.rootEndpoint+handler.GetEndpoint(), functionHandler(handler))
		case http.MethodDelete:
			e.DELETE(service.rootEndpoint+handler.GetEndpoint(), functionHandler(handler))
		case http.MethodHead:
			e.HEAD(service.rootEndpoint+handler.GetEndpoint(), functionHandler(handler))
		case http.MethodOptions:
			e.OPTIONS(service.rootEndpoint+handler.GetEndpoint(), functionHandler(handler))
		default:
			e.GET(service.rootEndpoint+handler.GetEndpoint(), functionHandler(handler))
		}
	}

	switch service.nlflowHandler.handler.GetMethod() {
	case http.MethodGet:
		e.GET(service.rootEndpoint+service.nlflowHandler.handler.GetEndpoint(), functionNLFlowHandler(service.nlflowHandler))
	case http.MethodPost:
		e.POST(service.rootEndpoint+service.nlflowHandler.handler.GetEndpoint(), functionNLFlowHandler(service.nlflowHandler))
	case http.MethodPut:
		e.PUT(service.rootEndpoint+service.nlflowHandler.handler.GetEndpoint(), functionNLFlowHandler(service.nlflowHandler))
	case http.MethodDelete:
		e.DELETE(service.rootEndpoint+service.nlflowHandler.handler.GetEndpoint(), functionNLFlowHandler(service.nlflowHandler))
	case http.MethodHead:
		e.HEAD(service.rootEndpoint+service.nlflowHandler.handler.GetEndpoint(), functionNLFlowHandler(service.nlflowHandler))
	case http.MethodOptions:
		e.OPTIONS(service.rootEndpoint+service.nlflowHandler.handler.GetEndpoint(), functionNLFlowHandler(service.nlflowHandler))
	default:
		e.GET(service.rootEndpoint+service.nlflowHandler.handler.GetEndpoint(), functionNLFlowHandler(service.nlflowHandler))
	}

	service.e = e
}

func (service *Service) Start() {
	service.init()
	service.e.Logger.Fatal(service.e.Start(fmt.Sprintf("%s:%d", service.host, service.port)))
}

func (service *Service) schema(i interface{}) (string, error) {
	t := reflect.TypeOf(i)
	schema := jsonschema.ReflectFromType(t)
	data, err := schema.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (service *Service) GenerateManifest(inputModel, outputModel Model) (Manifest, error) {
	inputSchemaAsString, err := service.schema(inputModel)
	if err != nil {
		return Manifest{}, err
	}
	outputSchemaAsString, err := service.schema(outputModel)
	if err != nil {
		return Manifest{}, err
	}
	var input interface{}
	err = json.Unmarshal([]byte(inputSchemaAsString), &input)
	if err != nil {
		return Manifest{}, err
	}
	var output interface{}
	err = json.Unmarshal([]byte(outputSchemaAsString), &output)
	if err != nil {
		return Manifest{}, err
	}
	return Manifest{
		ID:           service.info.ComponentId,
		Name:         service.info.ComponentName,
		Description:  service.info.Desc,
		Type:         service.info.ManifestType,
		Icon:         service.info.Icon,
		InputSchema:  input,
		OutputSchema: output,
	}, nil
}

package cache

type nlflowCacheMock struct {
}

func NewNLFlowCacheMock() NLFlowCache {
	return &nlflowCacheMock{}
}

func (c *nlflowCacheMock) Read(k string) (string, error) {
	return k, nil
}

func (*nlflowCacheMock) Write(k string, v string) error {
	return nil
}

func (*nlflowCacheMock) Close() error {
	return nil
}

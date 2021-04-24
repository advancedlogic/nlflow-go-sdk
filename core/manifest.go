package core

import "encoding/json"

type ManifestType string

const (
	CUSTOM    ManifestType = "CUSTOM"
	CONVERTER ManifestType = "CONVERTER"
)

type Manifest struct {
	// General Info
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Name      string `json:"name"`

	// UI
	Type ManifestType `json:"type"`
	//Label       string       `json:"label"`
	Icon        string `json:"icon"`
	Description string `json:"description"`

	// Networking
	//Endpoint string `json:"endpoint"`

	// Schema
	InputSchema  interface{} `json:"input-schema"`
	OutputSchema interface{} `json:"output-schema"`

	Error string
}

func (m Manifest) String() string {
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return ""
	}
	return string(data)
}

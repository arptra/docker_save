package parsemanifest

import (
	"encoding/json"
)

type BaseStruct struct {
	MediaType string
	Size      uint64
	Digest    string
}

type LayerStruct struct {
	Layers []BaseStruct
}

type ConfigStruct struct {
	Config BaseStruct
}

type Params struct {
	SchemaVersion int
	MediaType     string
}

type Manifest struct {
	FirstParams Params
	Config      *ConfigStruct
	Layers      *LayerStruct
}

func Parse(data string) Manifest {
	var manifest Params
	var config ConfigStruct
	var layers LayerStruct

	json.Unmarshal([]byte(data), &manifest)
	json.Unmarshal([]byte(data), &config)
	json.Unmarshal([]byte(data), &layers)

	return Manifest{
		manifest,
		&config,
		&layers,
	}
}

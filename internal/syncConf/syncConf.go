package syncconf

import (
	"github.com/sivaosorg/govm/apix"
	"github.com/sivaosorg/govm/configx"
)

// Conf global
// Get & Set key-value
// Based on conf, to create new cluster / instance
var Conf *configx.KeysConfig
var Params *KeyParams

type KeyParams struct {
	Curl []apix.ApiRequestConfig `json:"curl" yaml:"curl"`
}

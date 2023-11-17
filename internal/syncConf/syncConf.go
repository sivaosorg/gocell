package syncconf

import (
	"fmt"

	"github.com/sivaosorg/govm/apix"
	"github.com/sivaosorg/govm/coltx"
	"github.com/sivaosorg/govm/configx"
	"github.com/sivaosorg/govm/utils"
)

// Global Configs yaml
var Conf *configx.KeysConfig
var Params *keyParams
var Jobs *jobParams

type keyParams struct {
	Curl []apix.ApiRequestConfig `json:"curl" yaml:"curl"`
}

type jobParams struct {
}

type sync struct {
}

func NewSync() *sync {
	return &sync{}
}

func AvailableParams() bool {
	return Params != nil
}

func AvailableConf() bool {
	return Conf != nil
}

func AvailableJobs() bool {
	return Jobs != nil
}

func (s *sync) GetClusters(args []string) (*configx.KeysConfig, error) {
	keys, err := configx.ReadConfig[configx.KeysConfig](args[0])
	return keys, err
}

func (s *sync) GetParams(args []string) (*keyParams, bool, error) {
	index := 1
	if utils.IsEmpty(args[index]) {
		return nil, false, nil
	}
	if !coltx.IndexExists(args, index) {
		return nil, true, fmt.Errorf("Out of range args params: %v", index)
	}
	params, err := configx.ReadConfig[keyParams](args[index])
	return params, true, err
}

func (s *sync) GetJobs(args []string) (*jobParams, bool, error) {
	return nil, true, nil
}

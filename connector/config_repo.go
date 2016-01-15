package connector

import (
	"encoding/json"
	"io"
	"os"

	"github.com/coreos/dex/repo"
)

func newConnectorConfigsFromReader(r io.Reader) ([]ConnectorConfig, error) {
	var ms []map[string]interface{}
	if err := json.NewDecoder(r).Decode(&ms); err != nil {
		return nil, err
	}
	cfgs := make([]ConnectorConfig, len(ms))
	for i, m := range ms {
		cfg, err := newConnectorConfigFromMap(m)
		if err != nil {
			return nil, err
		}
		cfgs[i] = cfg
	}
	return cfgs, nil
}

func NewConnectorConfigRepoFromFile(loc string) (ConnectorConfigRepo, error) {
	cf, err := os.Open(loc)
	if err != nil {
		return nil, err
	}
	defer cf.Close()

	cfgs, err := newConnectorConfigsFromReader(cf)
	if err != nil {
		return nil, err
	}

	return &memConnectorConfigRepo{configs: cfgs}, nil
}

type memConnectorConfigRepo struct {
	configs []ConnectorConfig
}

func NewConnectorConfigRepoFromConfigs(cfgs []ConnectorConfig) ConnectorConfigRepo {
	return &memConnectorConfigRepo{configs: cfgs}
}

func (r *memConnectorConfigRepo) All() ([]ConnectorConfig, error) {
	return r.configs, nil
}

func (r *memConnectorConfigRepo) GetConnectorByID(_ repo.Transaction, id string) (ConnectorConfig, error) {
	for _, cfg := range r.configs {
		if cfg.ConnectorID() == id {
			return cfg, nil
		}
	}
	return nil, ErrorNotFound
}

func (r *memConnectorConfigRepo) Set(cfgs []ConnectorConfig) error {
	r.configs = cfgs
	return nil
}

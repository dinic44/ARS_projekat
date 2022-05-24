package configstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
)

type ConfigStore struct {
	cli *api.Client
}

func New() (*ConfigStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigStore{
		cli: client,
	}, nil
}

//Create Single
func (cs *ConfigStore) CreateSingleConfig(singleConfig *SingleConfig) (*SingleConfig, error) {
	kv := cs.cli.KV()

	sid, rid := generateSingleConfigKey(singleConfig.Version)
	singleConfig.Id = rid

	data, err := json.Marshal(singleConfig)
	if err != nil {
		return nil, err
	}

	c := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}

	return singleConfig, nil
}

//Put New {id}
func (cs *ConfigStore) PutNewSingleConfigVersion(singleConfig *SingleConfig) (*SingleConfig, error) {
	kv := cs.cli.KV()

	data, err := json.Marshal(singleConfig)
	if err != nil {
		return nil, err
	}

	_, err = cs.FindSingleConfig(singleConfig.Id, singleConfig.Version)

	if err == nil {
		return nil, errors.New("error! ")
	}

	c := &api.KVPair{Key: constructSingleConfigKey(singleConfig.Id, singleConfig.Version), Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}
	return singleConfig, nil

}

//Find One Single/{id}
func (cs *ConfigStore) FindSingleConfigVersion(id string) ([]*SingleConfig, error) {
	kv := cs.cli.KV()

	key := constructSingleConfigIdKey(id)
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, err
	}

	var singleConfigs []*SingleConfig

	for _, pair := range data {
		singleConfig := &SingleConfig{}
		err := json.Unmarshal(pair.Value, singleConfig)
		if err != nil {
			return nil, err
		}

		singleConfigs = append(singleConfigs, singleConfig)
	}

	return singleConfigs, nil
}

//Find One Single/{id}/{version}
func (cs *ConfigStore) FindSingleConfig(id string, ver string) (*SingleConfig, error) {
	kv := cs.cli.KV()
	key := constructSingleConfigKey(id, ver)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		return nil, errors.New("Cannot Find")
	}

	singleConfig := &SingleConfig{}
	err = json.Unmarshal(data.Value, singleConfig)
	if err != nil {
		return nil, err
	}

	return singleConfig, nil
}

//Find All Single
/*func (cs *ConfigStore) GetAllSingleConfig() ([]*SingleConfig, error) {
	kv := cs.cli.KV()
	data, _, err := kv.List(singleConfigAll, nil)
	if err != nil {
		return nil, err
	}

	singleConfigs := []*SingleConfig{}
	for _, pair := range data {
		singleConfig := &SingleConfig{}
		err = json.Unmarshal(pair.Value, singleConfig)
		if err != nil {
			return nil, err
		}
		singleConfigs = append(singleConfigs, singleConfig)
	}

	return singleConfigs, nil
}*/

//Delete Single
func (cs *ConfigStore) DeleteSingleConfig(id, version string) (map[string]string, error) {
	kv := cs.cli.KV()
	_, err := kv.Delete(constructSingleConfigKey(id, version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

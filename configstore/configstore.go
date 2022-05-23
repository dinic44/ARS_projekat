package configstore

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
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

	sid, rid := generateSingleConfigKey(singleConfig.Id)
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
func (cs *ConfigStore) FindSingleConfig(id string, version string) (*SingleConfig, error) {
	kv := cs.cli.KV()
	key := constructSingleConfigKey(id, version)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		return nil, err
	}

	singleConfig := &SingleConfig{}
	err = json.Unmarshal(data.Value, singleConfig)
	if err != nil {
		return nil, err
	}

	return singleConfig, nil
}

//Find All Single
func (cs *ConfigStore) GetAllSingleConfig() ([]*SingleConfig, error) {
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
}

//Create Group
func (cs *ConfigStore) CreateGroupConfig(groupConfig *GroupConfig) (*GroupConfig, error) {
	kv := cs.cli.KV()

	sid, rid := generateGroupConfigKey(groupConfig.Version)
	groupConfig.Id = rid

	log.Default().Println(sid, kv)

	data, err := json.Marshal(groupConfig)
	if err != nil {
		return nil, err
	}

	g := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(g, nil)
	if err != nil {
		return nil, err
	}

	/*for i, config := range group.Configs {
		cid := constructGroupLabel(rid, group.Version, i, config)
		log.Default().Println(cid)
	}*/

	return groupConfig, nil
}

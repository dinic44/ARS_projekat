package configstore

import (
	"encoding/json"
	"errors"
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

//Delete Single
func (cs *ConfigStore) DeleteSingleConfig(id, version string) (map[string]string, error) {
	kv := cs.cli.KV()
	_, err := kv.Delete(constructSingleConfigKey(id, version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

//Create Group
func (cs *ConfigStore) CreateGroupConfig(groupConfig *GroupConfig) (*GroupConfig, error) {
	kv := cs.cli.KV()

	sid, rid := generateGroupConfigKey(groupConfig.Version)
	groupConfig.Id = rid

	data, err := json.Marshal(groupConfig)
	if err != nil {
		return nil, err
	}

	gr := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(gr, nil)
	if err != nil {
		return nil, err
	}

	err = cs.CreateLabels(groupConfig.GroupConfig, groupConfig.Id, groupConfig.Version)
	if err != nil {
		return nil, err
	}

	return groupConfig, nil
}

//Put New {id}
func (cs *ConfigStore) PutNewGroupConfigVersion(groupConfig *GroupConfig) (*GroupConfig, error) {
	kv := cs.cli.KV()

	data, err := json.Marshal(groupConfig)
	if err != nil {
		return nil, err
	}

	_, err = cs.GetGroupConfig(groupConfig.Id, groupConfig.Version)

	if err == nil {
		return nil, errors.New("already exists ")
	}

	c := &api.KVPair{Key: constructGroupConfigKey(groupConfig.Id, groupConfig.Version), Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}

	err = cs.CreateLabels(groupConfig.GroupConfig, groupConfig.Id, groupConfig.Version)
	if err != nil {
		return nil, err
	}

	return groupConfig, nil
}

//Find {id}/{ver}
func (cs *ConfigStore) GetGroupConfig(id string, ver string) (*GroupConfig, error) {
	kv := cs.cli.KV()
	key := constructGroupConfigKey(id, ver)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		return nil, errors.New("cannot find")
	}

	groupConfig := &GroupConfig{}
	err = json.Unmarshal(data.Value, groupConfig)
	if err != nil {
		return nil, err
	}

	return groupConfig, nil
}

func (cs *ConfigStore) FindSingleInGroup(id, ver, KVPair string) ([]map[string]string, error) {
	kv := cs.cli.KV()
	singleKey := fmt.Sprintf(singleInGroup, id, ver, KVPair) + "/"
	keys, _, err := kv.List(singleKey, nil)
	if err != nil {
		return nil, err
	}

	configs := make([]map[string]string, len(keys))
	for i, k := range keys {
		var singleConfig map[string]string
		json.Unmarshal(k.Value, &singleConfig)
		log.Default().Printf("%q", singleConfig)
		configs[i] = singleConfig
	}

	return configs, nil
}

func (cs *ConfigStore) DeleteGroupConfig(id, ver string) error {
	kv := cs.cli.KV()

	_, err := kv.DeleteTree(constructGroupConfigKey(id, ver), nil)

	return err
}

func (cs *ConfigStore) SaveRequestId() string {

	kv := cs.cli.KV()

	reqId := generateRequestId()

	i := &api.KVPair{Key: reqId, Value: nil}

	_, err := kv.Put(i, nil)

	if err != nil {
		return "error"
	}

	return reqId
}

func (cs *ConfigStore) FindRequestId(requestId string) bool {

	kv := cs.cli.KV()

	key, _, err := kv.Get(requestId, nil)

	fmt.Println(key)

	if err != nil || key == nil {
		return false
	}

	return true
}

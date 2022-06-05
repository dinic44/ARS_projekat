package configstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"log"
	"sort"
)

type SingleConfig struct {
	Id      string            `json:"id"`
	Version string            `json:"version"`
	Entries map[string]string `json:"entries"`
}

type GroupConfig struct {
	Id          string              `json:"id"`
	GroupConfig []map[string]string `json:"configs"`
	Version     string              `json:"version"`
}

const (
	singleConfigId = "singleConfig/%s"
	singleConfig   = "singleConfig/%s/%s"
	/*	singleConfigAll = "singleConfigs"

		groupConfigAll     = "groupConfigs"
		groupConfigId      = "groupConfig/%s"*/
	groupConfigVersion = "groupConfig/%s/%s"
	singleInGroup      = "groupConfig/%s/%s/%s"
	groupConfigLabel   = "groupConfig/%s/%s/%s/%s"
)

func generateSingleConfigKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(singleConfig, id, ver), id
}

func constructSingleConfigKey(id string, ver string) string {
	return fmt.Sprintf(singleConfig, id, ver)
}

func constructSingleConfigIdKey(id string) string {
	return fmt.Sprintf(singleConfigId, id)
}

func generateGroupConfigKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(groupConfigVersion, id, ver), id
}

func constructGroupConfigKey(id string, ver string) string {
	return fmt.Sprintf(groupConfigVersion, id, ver)
}

func (cs *ConfigStore) CreateLabels(configs []map[string]string, id, ver string) error {
	kv := cs.cli.KV()
	if keys, _, err := kv.Get(constructGroupConfigKey(id, ver), nil); err != nil || keys == nil {
		return errors.New("error")
	}

	for _, config := range configs {
		cid := constructGroupLabel(id, ver, uuid.New().String(), config)
		cdata, err := json.Marshal(config)

		log.Default().Printf("adding new config", config, cdata)
		if err != nil {
			return err
		}

		c := &api.KVPair{Key: cid, Value: cdata}
		_, err = kv.Put(c, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func constructGroupLabel(id, ver, index string, config map[string]string) string {
	keys := make([]string, 0, len(config))
	for k := range config {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var kvpairs string
	for k := range keys {

		kvpairs = kvpairs + fmt.Sprintf("%s=%s", keys[k], config[keys[k]]+"&")
	}
	kvpairs = kvpairs[:len(kvpairs)-1]
	return fmt.Sprintf(groupConfigLabel, id, ver, kvpairs, index)
}

func generateRequestId() string {

	rid := uuid.New().String()

	return rid
}

package configstore

import (
	"fmt"
	"github.com/google/uuid"
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
	singleConfigId  = "singleConfig/%s"
	singleConfig    = "singleConfig/%s/%s"
	singleConfigAll = "singleConfigs"

	groupConfigId      = "groupConfig/%s"
	groupConfigVersion = "groupConfig/%s/%s"
	groupConfigAll     = "groupConfigs"
	singleInGroup      = "groupConfig/%s/%s/%s"
	groupConfigLabel   = "groupConfig/%s/%s/%s/%s" //??
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

func constructGroupConfigIdKey(id string) string {
	return fmt.Sprintf(groupConfigId, id)
}

/*func (cs *ConfigStore) CreateLabels(configs []map[string]string, id, ver string) error {
	kv := cs.cli.KV()
	if keys, _, err := kv.Get(constructGroupConfigKey(id, ver), nil); err != nil || keys == nil {
		return errors.New("Group doesn't exists")
	}

	for _, config := range configs {
		cid := constructGroupLabel(id, ver, uuid.New().String(), config)
		cdata, err := json.Marshal(config)

		log.Default().Printf("adding new config: %q. under key %q", config, cdata)
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
}*/

/*func constructGroupLabel(id, ver, index string, config map[string]string) string {
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
	return fmt.Sprintf(groupWithLabel, id, ver, kvpairs, index)
}*/

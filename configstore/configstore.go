package configstore

import (
	"encoding/json"
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

func (cs *ConfigStore) Get(id string) (*Config, error) {
	kv := cs.cli.KV()

	pair, _, err := kv.Get(constructKey(id), nil)
	if err != nil {
		return nil, err
	}

	post := &Config{}
	err = json.Unmarshal(pair.Value, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (cs *ConfigStore) GetAll() ([]*Config, error) {
	kv := cs.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	posts := []*Config{}
	for _, pair := range data {
		post := &Config{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (cs *ConfigStore) Delete(id string) (map[string]string, error) {
	kv := cs.cli.KV()
	_, err := kv.Delete(constructKey(id), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (cs *ConfigStore) Post(post *Config) (*Config, error) {
	kv := cs.cli.KV()

	sid, rid := generateKey()
	post.Id = rid

	data, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return post, nil
}

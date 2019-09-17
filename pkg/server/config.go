/*
 *    Copyright 2019 Samsung SDS
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package server

import (
	"io/ioutil"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/nexledger/accelerator/pkg/batch"
	"github.com/nexledger/accelerator/pkg/fabwrap"
)

type Tls struct {
	Certpath string `yaml:"certPath"`
	KeyPath  string `yaml:"keyPath"`
}

type Config struct {
	Sdk          string                   `yaml:"sdk"`
	Host         string                   `yaml:"host"`
	Port         int                      `yaml:"port"`
	UserName     string                   `yaml:"userName"`
	Organization string                   `yaml:"organization"`
	Tls          Tls                      `yaml:"tls"`
	Batch        []map[string]interface{} `yaml:"batch"`
}

func (c *Config) BatchClient() (*batch.Client, error) {
	var accs []*batch.Acceleration
	for _, configs := range c.Batch {
		for _, k := range []string{"type", "channelId", "chaincodeName", "fcn"} {
			if _, ok := configs[k]; !ok {
				return nil, errors.New(k + " doesn't exist in config")
			}
		}

		reg := &batch.Acceleration{
			QueueSize:          1000,
			MaxBatchItems:      10,
			MaxBatchBytes:      0,
			MaxWaitTimeSeconds: 10,
			ReadKeyIndices:     nil,
			WriteKeyIndices:    nil,
			Encoding:           "gob",
			Recovery:           false,
		}

		err := mapstructure.Decode(configs, &reg)
		if err != nil {
			return nil, err
		}

		accs = append(accs, reg)
	}

	ctx, err := fabwrap.New(c.Sdk, c.UserName, c.Organization)
	if err != nil {
		return nil, err
	}

	client := batch.New(ctx)
	for _, acc := range accs {
		if err := client.Register(acc); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func loadConfig(configPath string) (*Config, error) {
	confFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to read targeted config file")
	}

	conf := &Config{}
	err = yaml.Unmarshal(confFile, conf)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to unmarshal config file")
	}

	conf.Sdk, err = filepath.Abs(filepath.Join(filepath.Dir(configPath), conf.Sdk))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to convert sdk config file path")
	}

	return conf, nil
}

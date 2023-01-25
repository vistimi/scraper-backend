package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type ConfigDynamodb struct {
	Databases map[string]ConfigDynamodbTable `mapstructure:"dynamodb"`
}

type ConfigDynamodbTable struct {
	Name       *string `mapstructure:"name"`
	PrimaryKey *string `mapstructure:"primaryKey"`
	SortKey    *string `mapstructure:"sortKey"`
}

func ReadConfigFile(path string) (*ConfigDynamodb, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c ConfigDynamodb
	var raw interface{}

	if err := yaml.Unmarshal(f, &raw); err != nil {
		return nil, err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &c})
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(raw); err != nil {
		return nil, err
	}

	if len(c.Databases) == 0 {
		return nil, fmt.Errorf("no database found")
	}

	for tableReference, tableConfig := range c.Databases {
		if tableConfig.Name == nil {
			return nil, fmt.Errorf("no name found for table %s", tableReference)
		}
		if tableConfig.PrimaryKey == nil {
			return nil, fmt.Errorf("no primaryKey found for table %s", tableReference)
		}
		if tableConfig.SortKey == nil {
			return nil, fmt.Errorf("no sortKey found for table %s", tableReference)
		}
	}

	return &c, nil
}

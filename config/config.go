package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port      *int                           `mapstructure:"port"`
	Databases map[string]ConfigDynamodbTable `mapstructure:"dynamodb"`
	Buckets   map[string]ConfigS3Bucket      `mapstructure:"buckets"`
}

type ConfigDynamodbTable struct {
	Name           *string `mapstructure:"name"`
	PrimaryKeyName *string `mapstructure:"primaryKeyName"`
	PrimaryKeyType *string `mapstructure:"primaryKeyType"`
	SortKeyName    *string `mapstructure:"sortKeyName"`
	SortKeyType    *string `mapstructure:"sortKeyType"`
}

type ConfigS3Bucket struct {
	Name *string `mapstructure:"name"`
}

func ReadConfigFile(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
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

	if len(c.Buckets) == 0 {
		return nil, fmt.Errorf("no buckets found")
	}

	for tableReference, tableConfig := range c.Databases {
		if tableConfig.Name == nil || tableConfig.PrimaryKeyName == nil || tableConfig.PrimaryKeyType == nil || tableConfig.SortKeyName == nil || tableConfig.SortKeyType == nil {
			return nil, fmt.Errorf("element missing for table %s: %+#v", tableReference, tableConfig)
		}
	}

	for tableReference, bucketConfig := range c.Buckets {
		if bucketConfig.Name == nil {
			return nil, fmt.Errorf("element missing for bucket %s: %+#v", tableReference, bucketConfig)
		}
	}

	return &c, nil
}

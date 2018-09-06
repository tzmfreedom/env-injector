package envinjector

import (
	"gopkg.in/yaml.v2"

	"io/ioutil"
	"os"
)

func getDefaultDecorator() envKeyDecorator {
	if path := os.Getenv("ENV_INJECTOR_MAPPER"); path != "" {
		m := &keyMapper{}
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			logger.Fatalf("failed to read mappter file.\n %v", err)
		}

		err = yaml.Unmarshal(buf, &m.Mapping)
		if err != nil {
			logger.Fatalf("failed to unmarshal yaml string.\n %v", err)
		}
		return m
	}
	return noop
}

type keyMapper struct {
	Mapping map[string]string
}

func (m *keyMapper) decorate(key string) string {
	v, ok := m.Mapping[key]
	if ok {
		return v
	}

	return key
}

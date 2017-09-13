package pipeline

import (
	"bytes"
	"text/template"
)

type Context map[string]interface{}

func (c Context) ResolveConfig(cfg map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	for k, v := range cfg {
		switch casted := v.(type) {
		case map[string]interface{}:
			casted, err := c.ResolveConfig(casted)

			if err != nil {
				return nil, err
			}

			v = casted
		case string:
			casted, err := c.ResolveTemplate(k, casted)

			if err != nil {
				return nil, err
			}

			v = casted
		}

		result[k] = v
	}

	return result, nil
}

func (c Context) ResolveTemplate(name string, content string) (string, error) {
	buf := bytes.NewBuffer(nil)

	t, err := template.New(name).Parse(content)

	if err != nil {
		return "", err
	}

	err = t.ExecuteTemplate(buf, name, c)

	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

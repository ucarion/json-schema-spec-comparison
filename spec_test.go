package comparison

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/json-schema-spec/json-schema-go"
	"github.com/stretchr/testify/assert"
)

func TestSpec(t *testing.T) {
	remotes := []interface{}{}
	err := filepath.Walk("remotes", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var schema interface{}
		err = json.Unmarshal(data, &schema)
		if err != nil {
			return err
		}

		remotes = append(remotes, schema)
		return nil
	})

	assert.NoError(t, err)

	err = filepath.Walk("tests/draft7", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			data, err := ioutil.ReadFile(path)
			assert.NoError(t, err)

			var testCases []struct {
				Description string      `json:"description"`
				Schema      interface{} `json:"schema"`
				Tests       []struct {
					Description string      `json:"description"`
					Data        interface{} `json:"data"`
					Valid       bool        `json:"valid"`
				} `json:"tests"`
			}

			err = json.Unmarshal(data, &testCases)
			assert.NoError(t, err)

			for _, tt := range testCases {
				t.Run(tt.Description, func(t *testing.T) {
					schemas := []interface{}{tt.Schema}
					schemas = append(schemas, remotes...)

					validator, err := jsonschema.NewValidator(schemas)
					assert.NoError(t, err)
					if err != nil {
						return
					}

					schemaURI := url.URL{}
					if schemaObject, ok := tt.Schema.(map[string]interface{}); ok {
						if id, ok := schemaObject["$id"]; ok {
							uri, err := url.Parse(id.(string))
							assert.NoError(t, err)

							schemaURI = *uri
						}
					}

					for _, test := range tt.Tests {
						t.Run(test.Description, func(t *testing.T) {
							result, err := validator.ValidateURI(schemaURI, test.Data)
							assert.NoError(t, err)
							assert.Equal(t, test.Valid, result.IsValid())
						})
					}
				})
			}
		})

		return nil
	})

	assert.NoError(t, err)
}

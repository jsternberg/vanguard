package vanguard

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// note: use spaces in the literal string instead of tabs
const SampleYamlPlan = `---
tasks:
  - name: modify test.yml
    file:
      path: test.yml
      mode: 0644
`

func TestYamlLoader(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBufferString(SampleYamlPlan)
	loader := &YamlLoader{}
	if plan, err := loader.Load(buffer); assert.NoError(err) {
		if assert.Len(plan.Tasks, 1) {
			task := plan.Tasks[0]
			if assert.Equal("file", task.module.Name()) {
				config := task.config.(map[interface{}]interface{})
				assert.Equal("test.yml", config["path"])
				assert.Equal(0644, config["mode"])
			}
		}
	}
}

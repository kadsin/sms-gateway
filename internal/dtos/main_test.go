package dtos_test

import (
	"testing"

	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func Test_Marshal(t *testing.T) {
	data, _ := dtos.Marshal(map[string]any{"name": "mohsen"})

	assert.Equal(t, `{"name":"mohsen"}`, string(data))
}

func Test_Unmarshal(t *testing.T) {
	type Data struct {
		Name string `json:"name"`
	}

	d, _ := dtos.Unmarshal[Data]([]byte(`{"name":"sina"}`))

	assert.Equal(t, "sina", d.Name)
}

package pack

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMeta_BuildTime(t *testing.T) {
	a := assert.New(t)

	expectedTime := time.Now()

	meta := Meta{buildTime: expectedTime}

	a.Equal(expectedTime.UnixNano(), meta.BuildTime().UnixNano())
}

func TestMeta_UpdateBuildTime(t *testing.T) {
	a := assert.New(t)

	expectedTime := time.Now()

	meta := Meta{}
	meta.UpdateBuildTime(expectedTime)

	a.Equal(expectedTime, meta.buildTime)
}

package common_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUnix(t *testing.T) {

	tm1 := NewUnixNow()
	unixTimestamp := tm1.ToUnix()
	tm2 := NewUnix(unixTimestamp)

	assert.Equal(t, tm1, tm2, "TestToUnix")
}

func TestToFormat(t *testing.T) {

	_, err := NewFormat("2022/03/18 17:04:05")
	assert.NotNil(t, err)

	tm, err := NewFormat("2022-03-18 17:04:05")
	assert.Nil(t, err)

	t.Log("ToFormat -> ", tm.ToFormat())
}

func TestToFormatForTpl(t *testing.T) {

	_, err := NewFormat("2022/03/18 17:04:05")
	assert.NotNil(t, err)

	tm, err := NewFormat("2022-03-18 17:04:05")
	assert.Nil(t, err)

	t.Log("ToFormatForTpl -> ", tm.ToFormatForTpl("2006/01/02 15:04:05"))
}

func TestToIso8601(t *testing.T) {

	_, err := NewISO8601("2022-03-18 17:04:05")
	assert.NotNil(t, err)

	tm, err := NewISO8601("2006-01-02T15:04:05.999Z")
	assert.Nil(t, err)

	t.Log("ToIso8601 -> ", tm.ToIso8601())
}

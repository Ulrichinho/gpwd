package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLength(t *testing.T) {
	assert := assert.New(t)

	// default value of length : 12
	assert.GreaterOrEqual(length, 8)
	assert.LessOrEqual(length, 512)
}

func TestQuantity(t *testing.T) {
	assert := assert.New(t)

	// default value of quantity : 1
	assert.GreaterOrEqual(quantity, 1)
	assert.Less(quantity, 30)
}

func TestExport(t *testing.T) {
	// default value of export : false
	assert.Equal(t, export, false)
}

func TestStatistic(t *testing.T) {
	// default value of statistic : false
	assert.Equal(t, statistic, false)
}

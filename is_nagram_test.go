package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAnagram(t *testing.T) {
	freq := getFreq("ümberpööramine")
	some := "ümberpömiörane"

	assert.NotNil(t, isAnagram(freq, some))
}

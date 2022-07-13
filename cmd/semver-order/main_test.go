package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStdinJSON(t *testing.T) {

	t.Run("empty input", func(t *testing.T) {
		r := strings.NewReader("")
		record, err := parseStdinJSON(r)

		assert.Nil(t, err)
		assert.Len(t, record, 0)
	})

	t.Run("valid input", func(t *testing.T) {
		r := strings.NewReader(`["A","B"]`)
		record, err := parseStdinJSON(r)

		assert.Nil(t, err)
		assert.Len(t, record, 2)
	})

	t.Run("unexpected input", func(t *testing.T) {
		r := strings.NewReader(`["A",1]`)
		_, err := parseStdinJSON(r)

		assert.NotNil(t, err)
	})

}

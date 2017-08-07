package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		filter string
		lvs    map[string]string
		result bool
	}{
		// arithmetic comparison
		{
			filter: "resptime > 6",
			lvs:    map[string]string{"resptime": "10"},
			result: true,
		},
		{
			filter: "resptime >= 6",
			lvs:    map[string]string{"resptime": "6"},
			result: true,
		},
		{
			filter: "resptime == 60",
			lvs:    map[string]string{"resptime": "60"},
			result: true,
		},
		{
			filter: "resptime < 6",
			lvs:    map[string]string{"resptime": "7"},
			result: false,
		},
		{
			filter: "resptime <= 6",
			lvs:    map[string]string{"resptime": "7"},
			result: false,
		},
		// string comparison
		{
			filter: "uri == /top",
			lvs:    map[string]string{"uri": "/top"},
			result: true,
		},
		{
			filter: "uri ==* /TOP",
			lvs:    map[string]string{"uri": "/top"},
			result: true,
		},
		{
			filter: "uri != /top",
			lvs:    map[string]string{"uri": "/bottom"},
			result: true,
		},
		{
			filter: "uri !=* /top",
			lvs:    map[string]string{"uri": "/TOP"},
			result: false,
		},
		// regular expression
		{
			filter: "uri =~ ^/",
			lvs:    map[string]string{"uri": "/top"},
			result: true,
		},
		{
			filter: "uri !~ ^/",
			lvs:    map[string]string{"uri": "/top"},
			result: false,
		},
		{
			filter: "uri =~* ^/",
			lvs:    map[string]string{"uri": "/TOP"},
			result: true,
		},
		{
			filter: "uri !~* /top",
			lvs:    map[string]string{"uri": "/TOP"},
			result: false,
		},
	}

	for _, test := range tests {
		filters := getFuncFilters([]string{test.filter})
		for k, filter := range filters {
			assert.Equal(test.result, filter(test.lvs[k]))
		}
	}
}

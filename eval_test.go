package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		expression string
		vars       Vars
		result     string
	}{
		{ // no space
			expression: "1+1",
			vars:       nil,
			result:     "2",
		},
		{
			expression: "1 + 1",
			vars:       nil,
			result:     "2",
		},
		{
			expression: "resptime - upstream_resptime",
			vars:       Vars{"resptime": "5", "upstream_resptime": "3"},
			result:     "2",
		},
		{
			expression: "(resptime - upstream_resptime) * 1000",
			vars:       Vars{"resptime": "5", "upstream_resptime": "3"},
			result:     "2000",
		},
	}

	for _, test := range tests {
		expr, err := parseExpr(test.expression)
		assert.Nil(err)

		ctx := &ExprContext{
			vars: test.vars,
		}
		v, err := evalExpr(expr, ctx)
		assert.Nil(err)
		assert.Equal(test.result, v.String())
	}
}

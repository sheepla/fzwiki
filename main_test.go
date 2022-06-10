package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := []struct {
		desc string
		args []string
		want exitCode
	}{
		{
			desc: "normal: returns exitCodeOK when arguments have a short version option",
			args: []string{"-V"},
			want: exitCodeOK,
		},
		{
			desc: "normal: returns exitCodeOK when arguments have a long version option",
			args: []string{"--version"},
			want: exitCodeOK,
		},
		{
			desc: "normal: returns exitCodeOK when arguments have a help option",
			args: []string{"-h"},
			want: exitCodeOK,
		},
		{
			desc: "abnormal: returns exitCodeErr when arguments are empty",
			args: []string{},
			want: exitCodeErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)
			got, _ := run(tt.args)
			assert.Equal(tt.want, got)
		})
	}
}

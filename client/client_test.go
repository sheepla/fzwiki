package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSearchURL(t *testing.T) {
	tests := []struct {
		desc  string
		query string
		lang  string
		want  string
	}{
		{
			desc:  "normal: a language code is english (en) when lang is empty",
			query: "",
			lang:  "",
			want:  "https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=",
		},
		{
			desc:  "normal: a language code is japanese (ja) when lang is ja",
			query: "",
			lang:  "ja",
			want:  "https://ja.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=",
		},
		{
			desc:  "normal: query is encoded with parcent encoding",
			query: "あ",
			lang:  "ja",
			want:  "https://ja.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=%E3%81%82",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)
			got := CreateSearchURL(tt.query, tt.lang)
			assert.Equal(tt.want, got)
		})
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		desc    string
		url     string
		wantErr bool
	}{
		{
			desc:    "normal: search succeeds",
			url:     "https://ja.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=%E3%81%82",
			wantErr: false,
		},
		{
			desc:    "abnormal: raises error when url is empty",
			url:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got, err := Execute(tt.url)
			if tt.wantErr {
				assert.Error(err)
				return
			}

			assert.NoError(err)
			assert.NotNil(got)
		})
	}
}

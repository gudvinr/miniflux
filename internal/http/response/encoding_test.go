// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package response_test

import (
	"crypto/rand"
	"testing"

	"miniflux.app/v2/internal/http/response"
)

func TestAcceptEncoding(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		acceptEncoding string
		want           string
	}{
		{
			name:           "Empty input",
			acceptEncoding: "",
			want:           "identity",
		},
		{
			name:           "q=0 and identity",
			acceptEncoding: "identity;q=0",
			want:           "",
		},
		{
			name:           "q=0 and *",
			acceptEncoding: "*;q=0",
			want:           "",
		},
		{
			name:           "gzip",
			acceptEncoding: "gzip",
			want:           "gzip",
		},
		{
			name:           "gzip and br",
			acceptEncoding: "gzip,br",
			want:           "gzip",
		},
		{
			name:           "br and gzip",
			acceptEncoding: "br,gzip,deflate",
			want:           "br",
		},
		{
			name:           "unsupported encoding",
			acceptEncoding: "unknown",
			want:           "identity",
		},
		{
			name:           "empty encoding",
			acceptEncoding: ",",
			want:           "identity",
		},
		{
			name:           "multiple encodings and q=0",
			acceptEncoding: "gzip;q=0,br;q=0",
			want:           "identity",
		},
		{
			// We want br here but weights are not supported.
			name:           "multiple encodings and q values",
			acceptEncoding: "gzip;q=0.5,br;q=0.8",
			want:           "gzip",
		},
		{
			name:           "multiple encodings and wildcard",
			acceptEncoding: "*;q=0,gzip,br",
			want:           "gzip",
		},
		{
			name:           "multiple encodings and wildcard and q=0",
			acceptEncoding: "*;q=0,gzip,br;q=0",
			want:           "gzip",
		},
		{
			// We want br here but weights are not supported.
			name:           "multiple encodings and wildcard and q values",
			acceptEncoding: "*;q=0.5,gzip;q=0.8,br",
			want:           "gzip",
		},
		{
			name:           "multiple encodings and wildcard and q values and q=0",
			acceptEncoding: "*;q=0.5,gzip;q=0.8,br;q=0",
			want:           "gzip",
		},
		{
			name:           "invalid q value",
			acceptEncoding: "gzip;q=abc,deflate",
			want:           "deflate",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := response.AcceptEncoding(test.acceptEncoding)
			if got != test.want {
				t.Errorf("AcceptEncoding() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestAcceptEncodingCache(t *testing.T) {
	encoding := "identity;q=0,gzip," + rand.Text() // rand for avoid clashes with other test cases
	expected := "gzip"

	got := response.AcceptEncoding(encoding)
	if got != expected {
		t.Errorf("AcceptEncoding() = %q, want %q", got, expected)
	}

	got = response.AcceptEncoding(encoding)
	if got != expected {
		t.Errorf("AcceptEncoding() = %q, want %q", got, expected)
	}
}

// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package response

import (
	"slices"
	"strconv"
	"strings"
	"sync"
)

// encodingCache stores mapping between raw Accept-Encoding value
// and chosen encoding from this value.
var encodingCache sync.Map

// AcceptEncoding parses input string according to [HTTP Semantics]
// and returns first encoding that can be understood by us.
//
// Currently this function ignores set weights other than q=0.
// Encodings with q=0 will not be considered.
//
// If string is empty or no encoding was accepted function returns "identity".
//
// For "identity;q=0" and "*;q=0" function returns an empty string. In that case,
// if no other encoding was accepted, 406 Not Acceptable should be returned.
//
// [HTTP Semantics]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Encoding.
func AcceptEncoding(acceptEncoding string) string {
	acceptEncoding = strings.TrimSpace(acceptEncoding)
	if v, ok := encodingCache.Load(acceptEncoding); ok {
		return v.(string)
	}

	accepted := "identity"

	for enc := range strings.SplitSeq(acceptEncoding, ",") {
		enc = strings.TrimSpace(enc)

		if qi := strings.IndexByte(enc, ';'); qi > -1 {
			qstr := strings.TrimPrefix(enc[qi:], ";q=")

			q, err := strconv.ParseFloat(qstr, 64)
			if err != nil {
				continue // Ignore weird float values.
			}

			enc = enc[:qi]

			if q == 0 && slices.Contains([]string{"identity", "*"}, enc) {
				accepted = "" // Explicitly disabled, so can't be used as fallback.
				continue
			}

			if q == 0 {
				continue // Skipping unwanted.
			}
		}

		// List should be in sync with [Builder.Write].
		if !slices.Contains([]string{"br", "gzip", "deflate"}, enc) {
			continue // Skipping unsupported.
		}

		accepted = strings.ToLower(enc)
		break
	}

	// Store selection as it won't change for given header value.
	if v, ok := encodingCache.LoadOrStore(acceptEncoding, accepted); ok {
		return v.(string)
	}

	return accepted
}

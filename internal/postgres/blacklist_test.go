// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package postgres

import (
	"context"
	"testing"
)

func TestIsBlacklisted(t *testing.T) {
	defer ResetTestDB(testDB, t)
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	if _, err := testDB.exec(ctx, "INSERT INTO blacklist_prefixes (prefix) VALUES ('bad')"); err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		path string
		want bool
	}{
		{"fine", false},
		{"ba", false},
		{"bad", true},
		{"badness", true},
		{"bad.com/foo", true},
	} {
		got, err := testDB.IsBlacklisted(ctx, test.path)
		if err != nil {
			t.Fatal(err)
		}
		if got != test.want {
			t.Errorf("%q: got %t, want %t", test.path, got, test.want)
		}
	}
}
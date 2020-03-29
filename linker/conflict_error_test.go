package linker_test

import (
	"testing"

	"github.com/gbrlsnchs/pilgo/linker"
)

func TestConflictError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		testCases := []struct {
			cft  *linker.ConflictError
			want string
		}{
			{new(linker.ConflictError), "linker: there are 0 conflicts"},
			{&linker.ConflictError{Errs: make([]error, 1)}, "linker: there is 1 conflict"},
			{&linker.ConflictError{Errs: make([]error, 2)}, "linker: there are 2 conflicts"},
			{&linker.ConflictError{Errs: make([]error, 99)}, "linker: there are 99 conflicts"},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				errm := tc.cft.Error()
				if want, got := tc.want, errm; got != want {
					t.Fatalf("want %q, got %q", want, got)
				}
			})
		}
	})
}

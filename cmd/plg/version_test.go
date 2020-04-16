package main

import (
	"errors"
	"testing"

	"github.com/gbrlsnchs/cli/clitest"
	"github.com/google/go-cmp/cmp"
)

func TestVersion(t *testing.T) {
	testCases := []struct {
		desc    string
		version string
		cmd     versionCmd
		want    string
		err     error
	}{
		{
			desc:    "unknown version",
			version: "",
			cmd:     versionCmd{},
			want:    "version unknown version\n",
			err:     nil,
		},
		{
			desc:    "show version",
			version: "v0.0.0",
			cmd:     versionCmd{},
			want:    "version v0.0.0\n",
			err:     nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var (
				appcfg = appConfig{
					version: tc.version,
				}
				exec = tc.cmd.register(appcfg.copy)
				prg  = clitest.NewProgram("version")
				err  = exec(prg)
			)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, prg.CombinedOutput(); got != want {
				t.Fatalf("\"version\" command combined output mismatch (-want +got):\n%s",
					cmp.Diff(want, got))
			}
			if want, got := "", prg.ErrOutput(); got != want {
				t.Fatalf("\"version\" command stderr output mismatch (-want +got):\n%s",
					cmp.Diff(want, got))
			}
		})
	}
}

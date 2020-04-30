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
			want:    "version_cmd unknown version\n",
			err:     nil,
		},
		{
			desc:    "show raw version",
			version: "refs/tags/0.0.0-test",
			cmd:     versionCmd{},
			want:    "version_cmd 0.0.0-test\n",
			err:     nil,
		},
		{
			desc:    "show version",
			version: "0.0.0-test",
			cmd:     versionCmd{},
			want:    "version_cmd 0.0.0-test\n",
			err:     nil,
		},
		{
			desc:    "show version with prefix",
			version: "0.0.0-test",
			cmd:     versionCmd{},
			want:    "version_cmd 0.0.0-test\n",
			err:     nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var (
				appcfg = appConfig{
					name:    "version_cmd",
					version: tc.version,
				}
				exec = tc.cmd.register(appcfg.copy)
				prg  = clitest.NewProgram("cmd")
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

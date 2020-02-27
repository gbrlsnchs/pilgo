package pilgrim_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim"
)

func TestConfig(t *testing.T) {
	t.Run("Init", testConfigInit)
}

func testConfigInit(t *testing.T) {
	testCases := []struct {
		c       pilgrim.Config
		targets []string
		want    pilgrim.Config
	}{
		{
			c: pilgrim.Config{
				BaseDir: "test",
			},
			targets: []string{"foo", "bar"},
			want: pilgrim.Config{
				BaseDir: "test",
				Targets: []string{"foo", "bar"},
			},
		},
		{
			c: pilgrim.Config{
				BaseDir: "test",
			},
			targets: []string{"foo", "bar", "pilgrim.yml"},
			want: pilgrim.Config{
				BaseDir: "test",
				Targets: []string{"foo", "bar"},
			},
		},
		{
			c: pilgrim.Config{
				BaseDir: "test",
			},
			targets: []string{"foo", "bar", ".git"},
			want: pilgrim.Config{
				BaseDir: "test",
				Targets: []string{"foo", "bar"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			c := tc.c.Init(tc.targets)
			if want, got := tc.want, c; !cmp.Equal(got, want) {
				t.Errorf("Config.Init mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

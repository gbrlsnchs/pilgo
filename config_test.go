package pilgrim_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim"
)

func TestConfig(t *testing.T) {
	t.Run("Init", testConfigInit)
	t.Run("Set", testConfigSet)
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

func testConfigSet(t *testing.T) {
	testCases := []struct {
		c    pilgrim.Config
		name string
		o    pilgrim.Config
		want pilgrim.Config
	}{
		{
			c: pilgrim.Config{
				Targets: []string{"foo"},
				Options: map[string]pilgrim.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
					},
				},
			},
			name: "foo",
			o: pilgrim.Config{
				BaseDir: "test",
				Targets: []string{
					"bar",
				},
			},
			want: pilgrim.Config{
				Targets: []string{"foo"},
				Options: map[string]pilgrim.Config{
					"foo": {
						BaseDir: "test",
						Targets: []string{
							"bar",
						},
					},
				},
			},
		},
		{
			c: pilgrim.Config{
				Targets: []string{"foo"},
			},
			name: "foo",
			o: pilgrim.Config{
				BaseDir: "test",
				Targets: []string{
					"bar",
				},
			},
			want: pilgrim.Config{
				Targets: []string{"foo"},
				Options: map[string]pilgrim.Config{
					"foo": {
						BaseDir: "test",
						Targets: []string{
							"bar",
						},
					},
				},
			},
		},
		{
			c: pilgrim.Config{
				Targets: []string{"foo"},
				Options: map[string]pilgrim.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
					},
				},
			},
			name: "foo",
			o: pilgrim.Config{
				BaseDir: "test",
			},
			want: pilgrim.Config{
				Targets: []string{"foo"},
				Options: map[string]pilgrim.Config{
					"foo": {
						BaseDir: "test",
					},
				},
			},
		},
		{
			c: pilgrim.Config{
				Targets: []string{"foo"},
				Options: map[string]pilgrim.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
					},
				},
			},
			name: filepath.Join("foo", "bar"),
			o: pilgrim.Config{
				Targets: []string{"baz"},
			},
			want: pilgrim.Config{
				Targets: []string{"foo"},
				Options: map[string]pilgrim.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
						Options: map[string]pilgrim.Config{
							"bar": {
								Targets: []string{
									"baz",
								},
							},
						},
					},
				},
			},
		},
		{
			c:    pilgrim.Config{},
			name: "",
			o: pilgrim.Config{
				Targets: []string{"foo"},
			},
			want: pilgrim.Config{
				Targets: []string{"foo"},
				Options: nil,
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			tc.c.Set(tc.name, tc.o)
			if want, got := tc.want, tc.c; !cmp.Equal(got, want) {
				t.Errorf("Config.Set mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

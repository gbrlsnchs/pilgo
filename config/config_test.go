package config_test

import (
	"path/filepath"
	"testing"

	"github.com/gbrlsnchs/pilgo/config"
	"github.com/google/go-cmp/cmp"
)

func TestConfig(t *testing.T) {
	t.Run("Init", testConfigInit)
	t.Run("Set", testConfigSet)
}

func testConfigInit(t *testing.T) {
	testCases := []struct {
		c        config.Config
		targets  []string
		includes map[string]struct{}
		excludes map[string]struct{}
		want     config.Config
	}{
		{
			c: config.Config{
				BaseDir: "test",
			},
			targets: []string{"foo", "bar"},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"foo", "bar"},
			},
		},
		{
			c: config.Config{
				BaseDir: "test",
			},
			targets: []string{"foo", "bar", "pilgo.yml"},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"foo", "bar"},
			},
		},
		{
			c: config.Config{
				BaseDir: "test",
			},
			targets: []string{"foo", "bar", ".git"},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"foo", "bar"},
			},
		},
		{
			c: config.Config{
				BaseDir: "test",
			},
			targets: []string{"foo", "bar", ".git"},
			excludes: map[string]struct{}{
				"bar": struct{}{},
			},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"foo"},
			},
		},
		{
			c: config.Config{
				BaseDir: "test",
			},
			targets: []string{"foo", "bar", ".git"},
			includes: map[string]struct{}{
				"bar": struct{}{},
			},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"bar"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			c := tc.c.Init(tc.targets, tc.includes, tc.excludes)
			if want, got := tc.want, c; !cmp.Equal(got, want) {
				t.Errorf("Config.Init mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testConfigSet(t *testing.T) {
	testCases := []struct {
		c    config.Config
		name string
		o    config.Config
		want config.Config
	}{
		{
			c: config.Config{
				Targets: []string{"foo"},
				Options: map[string]config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
					},
				},
			},
			name: "foo",
			o: config.Config{
				BaseDir: "test",
				Targets: []string{
					"bar",
				},
			},
			want: config.Config{
				Targets: []string{"foo"},
				Options: map[string]config.Config{
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
			c: config.Config{
				Targets: []string{"foo"},
			},
			name: "foo",
			o: config.Config{
				BaseDir: "test",
				Targets: []string{
					"bar",
				},
			},
			want: config.Config{
				Targets: []string{"foo"},
				Options: map[string]config.Config{
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
			c: config.Config{
				Targets: []string{"foo"},
				Options: map[string]config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
					},
				},
			},
			name: "foo",
			o: config.Config{
				BaseDir: "test",
			},
			want: config.Config{
				Targets: []string{"foo"},
				Options: map[string]config.Config{
					"foo": {
						BaseDir: "test",
					},
				},
			},
		},
		{
			c: config.Config{
				Targets: []string{"foo"},
				Options: map[string]config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
					},
				},
			},
			name: filepath.Join("foo", "bar"),
			o: config.Config{
				Targets: []string{"baz"},
			},
			want: config.Config{
				Targets: []string{"foo"},
				Options: map[string]config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
						Options: map[string]config.Config{
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
			c:    config.Config{},
			name: "",
			o: config.Config{
				Targets: []string{"foo"},
			},
			want: config.Config{
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

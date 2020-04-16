package config_test

import (
	"path/filepath"
	"testing"

	"github.com/gbrlsnchs/pilgo/config"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var ignoreUnexported = cmpopts.IgnoreUnexported(config.Config{})

func TestNew(t *testing.T) {
	testCases := []struct {
		targets []string
		opts    []func(*config.Config)
		want    config.Config
	}{
		{
			targets: []string{"foo", "bar"},
			opts: []func(*config.Config){
				config.MergeWith(config.Config{
					BaseDir: "test",
				}),
			},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"foo", "bar"},
			},
		},
		{
			targets: []string{"foo", "bar", ".git"},
			opts: []func(*config.Config){
				config.MergeWith(config.Config{
					BaseDir: "test",
				}),
			},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"foo", "bar"},
			},
		},
		{
			targets: []string{"foo", "bar", ".git"},
			opts: []func(*config.Config){
				config.MergeWith(config.Config{
					BaseDir: "test",
				}),
				config.Exclude(map[string]struct{}{
					"bar": struct{}{},
				}),
			},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"foo"},
			},
		},
		{
			targets: []string{"foo", "bar", ".git"},
			opts: []func(*config.Config){
				config.MergeWith(config.Config{
					BaseDir: "test",
				}),
				config.Include(map[string]struct{}{
					"bar": struct{}{},
				}),
			},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"bar"},
			},
		},
		{
			targets: []string{"foo", "bar", ".git"},
			opts: []func(*config.Config){
				config.MergeWith(config.Config{
					BaseDir: "test",
				}),
				config.IncludeHidden,
			},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"foo", "bar", ".git"},
			},
		},
		{
			targets: []string{"foo", "", "baz"},
			opts: []func(*config.Config){
				config.MergeWith(config.Config{
					BaseDir: "test",
				}),
			},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{"foo", "baz"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			c := config.New(tc.targets, tc.opts...)
			if want, got := tc.want, c; !cmp.Equal(got, want, ignoreUnexported) {
				t.Errorf("Config.Init mismatch (-want +got):\n%s", cmp.Diff(want, got, ignoreUnexported))
			}
		})
	}
}

func TestConfig(t *testing.T) {
	t.Run("Set", testConfigSet)
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
			if want, got := tc.want, tc.c; !cmp.Equal(got, want, ignoreUnexported) {
				t.Errorf("Config.Set mismatch (-want +got):\n%s", cmp.Diff(want, got, ignoreUnexported))
			}
		})
	}
}

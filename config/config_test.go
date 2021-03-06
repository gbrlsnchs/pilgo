package config_test

import (
	"path/filepath"
	"testing"

	"github.com/gbrlsnchs/pilgo/config"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var ignoreUnexported = cmpopts.IgnoreUnexported(config.Config{})

func TestConfig(t *testing.T) {
	t.Run("Set", testConfigSet)
}

func testConfigSet(t *testing.T) {
	testCases := []struct {
		c    config.Config
		name string
		o    config.Config
		m    config.SetMode
		want config.Config
	}{
		{
			c: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
						Flatten: true,
					},
				},
			},
			name: "foo",
			o: config.Config{
				BaseDir: "test",
				Flatten: false,
			},
			want: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						BaseDir: "test",
						Targets: []string{
							"bar",
						},
						Flatten: false,
					},
				},
			},
		},
		{
			c: config.Config{
				Targets: []string{"foo"},
				Flatten: true,
			},
			name: "",
			o: config.Config{
				BaseDir: "test",
				Targets: []string{
					"bar",
				},
				Flatten: false,
			},
			want: config.Config{
				BaseDir: "test",
				Targets: []string{
					"foo",
				},
				Flatten: false,
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
				Options: map[string]*config.Config{
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
				Options: map[string]*config.Config{
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
				Options: map[string]*config.Config{
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
				Options: map[string]*config.Config{
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
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
						Options: map[string]*config.Config{
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
				Options: nil,
			},
		},
		{
			c: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						BaseDir: "test",
					},
				},
			},
			name: "foo",
			o:    config.Config{},
			want: config.Config{
				Targets: []string{"foo"},
				Options: nil,
			},
		},
		{
			c: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
					},
				},
			},
			name: "foo",
			o: config.Config{
				Flatten: true,
			},
			want: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						Flatten: true,
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
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
					},
				},
			},
			name: "foo",
			o: config.Config{
				Tags: []string{"test"},
			},
			want: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
						Tags: []string{
							"test",
						},
					},
				},
			},
		},
		{
			c: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
					},
				},
			},
			name: filepath.Join("foo", "bar"),
			o:    config.Config{},
			want: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
						Options: nil,
					},
				},
			},
		},
		{
			c:    config.Config{},
			name: filepath.Join("foo", "bar"),
			o: config.Config{
				BaseDir: "test",
				Targets: []string{
					"baz",
				},
			},
			want: config.Config{Options: nil},
		},
		{
			c: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
						Flatten: true,
					},
				},
			},
			name: "foo",
			o: config.Config{
				BaseDir: "test",
				Flatten: false,
				Targets: []string{"test"},
			},
			m: config.ModeScan,
			want: config.Config{
				Targets: []string{"foo"},
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"test",
						},
						Flatten: true,
					},
				},
			},
		},
		{
			c: config.Config{
				Targets: []string{"foo"},
				Flatten: true,
			},
			name: "",
			o: config.Config{
				BaseDir: "test",
				Targets: []string{
					"bar",
				},
				Flatten: false,
			},
			m: config.ModeScan,
			want: config.Config{
				Targets: []string{
					"bar",
				},
				Flatten: true,
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			tc.c.Set(tc.name, &tc.o, tc.m)
			if want, got := tc.want, tc.c; !cmp.Equal(got, want, ignoreUnexported) {
				t.Errorf("(*Config).Set mismatch (-want +got):\n%s", cmp.Diff(want, got, ignoreUnexported))
			}
		})
	}
}

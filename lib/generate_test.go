package lib

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	tests := map[string]struct {
		schedule       func() Schedule
		rebaseStrategy string
		configChecker  func(config Configuration) bool
	}{
		"A default generate": {
			schedule: func() Schedule {
				return Schedule{
					Interval: "daily",
					Time:     "",
					TimeZone: "",
				}
			},
			rebaseStrategy: "",
			configChecker:  func(config Configuration) bool { return true },
		},
		"A default generate, with rebase strategy disabled": {
			schedule: func() Schedule {
				return Schedule{
					Interval: "daily",
					Time:     "",
					TimeZone: "",
				}
			},
			rebaseStrategy: "disabled",
			configChecker: func(config Configuration) bool {
				for _, c := range config.Updates {
					if c.RebaseStrategy != "disabled" {
						return false
					}
				}
				return true
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			afs := &afero.Afero{Fs: afero.NewMemMapFs()}

			_, err := afs.Create("go.mod")
			assert.NoError(err)
			_, err = afs.Create("requirements.txt")
			assert.NoError(err)

			config, err := Generate(afs, ".", test.rebaseStrategy, test.schedule())
			assert.NoError(err)
			assert.Equal(2, config.Version)
			assert.Len(config.Updates, 2)
			assert.NotNil(config)
			assert.True(test.configChecker(config))
		})
	}
}

func TestDirectoryParserGitHub(t *testing.T) {
	afs := &afero.Afero{Fs: afero.NewMemMapFs()}

	_ = afs.Mkdir(".github", 0755)
	_ = afs.Mkdir(".github/workflows", 0755)
	_, _ = afs.Create(".github/workflows/ci.yml")
	results := directoryParser(afs, `(.*)\.yml`, ".")

	assert.Len(t, results, 1)
	assert.Equal(t, "/.github/workflows", results[0])
}

func TestDirectoryParserGoMod(t *testing.T) {
	afs := &afero.Afero{Fs: afero.NewMemMapFs()}

	_, _ = afs.Create("go.mod")
	results := directoryParser(afs, `go\.mod`, ".")

	assert.Len(t, results, 1)
	assert.Equal(t, "/", results[0])

	if results[0] != "/" {
		t.Error()
	}
}

func TestRemoveDuplicates(t *testing.T) {
	sliceWithDupes := []string{"1", "1", "2", "2", "3", "3"}
	correctSlice := []string{"1", "2", "3"}
	uniqueSlice := removeDuplicates(sliceWithDupes)

	assert.ElementsMatch(t, correctSlice, uniqueSlice)
}

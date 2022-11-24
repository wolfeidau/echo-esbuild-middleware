package assets

import (
	"testing"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	assert := require.New(t)

	// register the asset bundler which will build then serve any asset files
	err := BuildWithConfig(BuildConfig{
		EntryPoints: []api.EntryPoint{
			{
				InputPath:  "testassets/src/index.ts",
				OutputPath: "bundle",
			},
		},
		InlineSourcemap: true,
		Define: map[string]string{
			"process.env.NODE_ENV": `"production"`,
		},
		OnBuild: func(result api.BuildResult, timeTaken time.Duration) {
			assert.Len(result.Errors, 0)
		},
		Outdir: "public/js",
	})
	assert.NoError(err)
	assert.FileExists("public/js/bundle.js")
}

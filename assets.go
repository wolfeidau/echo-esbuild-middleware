package assets

import (
	"fmt"
	"path"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/labstack/echo/v4"
)

// FuncBuildResult callback used to return result when bundling
type FuncBuildResult func(result api.BuildResult, timeTaken time.Duration)

// BundlerConfig asset bundler configuration which provides the bare minimum to keep things simple
type BundlerConfig struct {
	EntryPoints     []string
	Outfile         string
	InlineSourcemap bool
	Define          map[string]string
	// This will be invoked for a build and can be used to check errors/warnings
	OnBuild FuncBuildResult
}

func (bc BundlerConfig) sourcemap() api.SourceMap {
	return api.SourceMapInline
}

// BundlerWithConfig provide bundle files which are built on startup
func BundlerWithConfig(cfg BundlerConfig) echo.MiddlewareFunc {

	if cfg.Outfile == "" {
		cfg.Outfile = "bundle.js"
	}

	start := time.Now()

	result := api.Build(api.BuildOptions{
		Banner:            `/* generated by github.com/wolfeidau/echo-esbuild-middleware */`,
		Bundle:            true,
		EntryPoints:       cfg.EntryPoints,
		Outfile:           cfg.Outfile,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Sourcemap:         cfg.sourcemap(),
		Define:            cfg.Define,
	})

	if cfg.OnBuild != nil {
		cfg.OnBuild(result, time.Since(start))
	}

	filesMap := outputFilesToMap(result.OutputFiles)

	fmt.Println(filesMap)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			contents, ok := filesMap[c.Path()]

			// if this isn't the file we are looking for invoke next handler and return
			if !ok {
				return next(c)
			}

			return c.Blob(200, "application/javascript", contents)
		}
	}
}

func outputFilesToMap(outputFiles []api.OutputFile) map[string][]byte {

	m := make(map[string][]byte)
	for _, f := range outputFiles {

		assetPath := fmt.Sprintf("/%s", path.Base(f.Path))

		m[assetPath] = f.Contents
	}

	return m
}

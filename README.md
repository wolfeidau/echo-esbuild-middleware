# echo-esbuild-middleware

This module provides an [echo](https://echo.labstack.com/) middleware which automatically bundles assets using [esbuild](https://github.com/evanw/esbuild).

[![GitHub Actions status](https://github.com/wolfeidau/echo-esbuild-middleware/workflows/Go/badge.svg?branch=master)](https://github.com/wolfeidau/echo-esbuild-middleware/actions?query=workflow%3AGo) 
[![Go Report Card](https://goreportcard.com/badge/github.com/wolfeidau/echo-esbuild-middleware)](https://goreportcard.com/report/github.com/wolfeidau/echo-esbuild-middleware) 
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/wolfeidau/echo-esbuild-middleware)

# Why?

I am currently using to enable me to add typescript based stimulus controllers to a website while keeping the bundling process simple and integrated into the development lifecycle of the service so I can "watch" and rebuild everything.

Given the single NPM module provided by esbuild my `node_modules` folder is also lean, which means less security issues.

# Usage

The following example will provide `bundle.js` via a script tag like `<script src="/bundle.js"></script>`.

```go
	e := echo.New()

    // register the asset bundler which will build then serve any asset files
    e.Use(assets.BundlerWithConfig(assets.BundlerConfig{
		EntryPoints:     []string{"testassets/src/index.ts"},
		Outfile:         "bundle.js",
		InlineSourcemap: true,
		Define: map[string]string{
			"process.env.NODE_ENV": `"production"`,
		},
		OnBuild: func(result api.BuildResult, timeTaken time.Duration) {
			if len(result.Errors) > 0 {
				log.Fatal().Fields(map[string]interface{}{
					"errors": result.Errors,
				}).Msg("failed to build assets")
			}
		},
        OnRequest: func(path string, contentLength, code int, timeTaken time.Duration) {
            log.Info().Str("path", path).Int("code", code).Duration("timeTaken", timeTaken).Msg("asset served")
        },
	}))
```


# License

This code was authored by [Mark Wolfe](https://www.wolfe.id.au) and licensed under the [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0).

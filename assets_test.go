package assets

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	assert := require.New(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bundle.js", nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/bundle.js")

	// register the asset bundler which will build then serve any asset files
	h := BundlerWithConfig(BundlerConfig{
		EntryPoints:     []string{"testassets/src/index.ts"},
		InlineSourcemap: true,
		Define: map[string]string{
			"process.env.NODE_ENV": `"production"`,
		},
	})(func(c echo.Context) error {
		return c.NoContent(404)
	})
	err := h(c)

	assert.NoError(err)
	assert.Equal(200, rec.Code)
}

func TestMiddlewareOnError(t *testing.T) {
	assert := require.New(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bundle.js", nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/bundle.js")

	// register the asset bundler which will build then serve any asset files
	h := BundlerWithConfig(BundlerConfig{
		EntryPoints:     []string{"testassets/src/indexs.ts"},
		InlineSourcemap: true,
		Define: map[string]string{
			"process.env.NODE_ENV": `"production"`,
		},
		OnBuild: func(result api.BuildResult, timeTaken time.Duration) {
			assert.Len(result.Errors, 1)
		},
	})(func(c echo.Context) error {
		return c.NoContent(404)
	})
	err := h(c)

	assert.NoError(err)
	assert.Equal(404, rec.Code)
}

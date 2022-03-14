package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/backendplugin"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/web/webtest"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServer_MetricsBasicAuth(t *testing.T) {
	ts := &HTTPServer{
		Cfg: setting.NewCfg(),
	}

	t.Run("enabled", func(t *testing.T) {
		ts.Cfg.MetricsEndpointBasicAuthUsername = "foo"
		ts.Cfg.MetricsEndpointBasicAuthPassword = "bar"

		assert.True(t, ts.metricsEndpointBasicAuthEnabled())
	})

	t.Run("disabled", func(t *testing.T) {
		ts.Cfg.MetricsEndpointBasicAuthUsername = ""
		ts.Cfg.MetricsEndpointBasicAuthPassword = ""

		assert.False(t, ts.metricsEndpointBasicAuthEnabled())
	})
}

func TestHTTPServer_PluginMetricsEndpoint(t *testing.T) {
	t.Run("Endpoint is enabled", func(t *testing.T) {
		hs := &HTTPServer{
			Cfg: &setting.Cfg{
				MetricsEndpointEnabled: true,
			},
			pluginClient: &fakePluginClientMetrics{
				store: map[string][]byte{
					"test-plugin": []byte("http_errors=2"),
				},
			},
		}

		s := webtest.NewServer(t, routing.NewRouteRegister())
		s.Mux.Use(hs.pluginMetricsEndpoint)

		t.Run("Endpoint matches and plugin is registered", func(t *testing.T) {
			req := s.NewGetRequest("/metrics/plugins/test-plugin")
			resp, err := s.Send(req)
			require.NoError(t, err)
			require.NotNil(t, resp)

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, "http_errors=2", string(body))
			require.NoError(t, resp.Body.Close())
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
		})

		t.Run("Endpoint matches and plugin is not registered", func(t *testing.T) {
			req := s.NewGetRequest("/metrics/plugins/plugin-not-registered")
			resp, err := s.Send(req)
			require.NoError(t, err)
			require.NotNil(t, resp)

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Empty(t, string(body))
			require.NoError(t, resp.Body.Close())
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Endpoint does not match", func(t *testing.T) {
			req := s.NewGetRequest("/foo")
			resp, err := s.Send(req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NoError(t, resp.Body.Close())
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("Endpoint is disabled", func(t *testing.T) {
		hs := &HTTPServer{
			Cfg: &setting.Cfg{
				MetricsEndpointEnabled: false,
			},
			pluginClient: &fakePluginClientMetrics{
				store: map[string][]byte{
					"test-plugin": []byte("http_errors=2"),
				},
			},
		}

		s := webtest.NewServer(t, routing.NewRouteRegister())
		s.Mux.Use(hs.pluginMetricsEndpoint)

		t.Run("When plugin is registered, should return 404", func(t *testing.T) {
			req := s.NewGetRequest("/metrics/plugins/test-plugin")
			resp, err := s.Send(req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NoError(t, resp.Body.Close())
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})
}

type fakePluginClientMetrics struct {
	plugins.Client

	store map[string][]byte
}

func (c *fakePluginClientMetrics) CollectMetrics(ctx context.Context, req *backend.CollectMetricsRequest) (*backend.CollectMetricsResult, error) {
	metrics, exists := c.store[req.PluginContext.PluginID]

	if !exists {
		return nil, backendplugin.ErrPluginNotRegistered
	}

	return &backend.CollectMetricsResult{
		PrometheusMetrics: metrics,
	}, nil
}

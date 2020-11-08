package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	http_routes "github.com/chacal/thread-mgmt-server/pkg/mgmt_routes/http"
	"github.com/chacal/thread-mgmt-server/pkg/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Devices(t *testing.T) {
	dbFile := test.Tempfile()
	reg, err := device_registry.Open(dbFile)
	require.NoError(t, err)
	defer reg.Close()

	router := gin.Default()
	http_routes.RegisterRoutes(router, reg)

	assertOKJson(t, `{}`, recordGet(router, "/v1/devices"))

	err = reg.Update("12345", device_registry.Device{"D100", 5000})
	require.NoError(t, err)

	assertOKJson(t, `{"12345": {"name": "D100", "pollTime": 5000}}`, recordGet(router, "/v1/devices"))
}

func assertOKJson(t *testing.T, expected string, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, expected, w.Body.String())
}

func recordGet(router *gin.Engine, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	router.ServeHTTP(w, req)
	return w
}

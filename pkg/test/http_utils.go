package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func AssertOK(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, w.Code)
}

func AssertOKJson(t *testing.T, expected string, w *httptest.ResponseRecorder) {
	AssertOK(t, w)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, expected, w.Body.String())
}

func RecordPost(router *gin.Engine, path string, payload string) *httptest.ResponseRecorder {
	return recordReq(router, path, "POST", payload)
}

func RecordDelete(router *gin.Engine, path string) *httptest.ResponseRecorder {
	return recordReq(router, path, "DELETE", "")
}

func RecordGet(router *gin.Engine, path string) *httptest.ResponseRecorder {
	return recordReq(router, path, "GET", "")
}

func recordReq(router *gin.Engine, path string, method string, payload string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, strings.NewReader(payload))
	router.ServeHTTP(w, req)
	return w
}

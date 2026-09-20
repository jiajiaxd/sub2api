package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestContentModerationConfigWhitelistedUsersAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newTestSettingRepo()
	handler := NewContentModerationHandler(service.NewContentModerationService(repo, nil, nil, nil, nil, nil, nil, nil))
	router := gin.New()
	router.PUT("/config", handler.UpdateConfig)
	router.GET("/config", handler.GetConfig)

	request := httptest.NewRequest(http.MethodPut, "/config", bytes.NewBufferString(`{"whitelisted_user_ids":[42,7,42]}`))
	request.Header.Set("Content-Type", "application/json")
	result := httptest.NewRecorder()
	router.ServeHTTP(result, request)
	require.Equal(t, http.StatusOK, result.Code)
	var updated struct {
		Data service.ContentModerationConfigView `json:"data"`
	}
	require.NoError(t, json.Unmarshal(result.Body.Bytes(), &updated))
	require.Equal(t, []int64{7, 42}, updated.Data.WhitelistedUserIDs)

	result = httptest.NewRecorder()
	router.ServeHTTP(result, httptest.NewRequest(http.MethodGet, "/config", nil))
	require.Equal(t, http.StatusOK, result.Code)
	var loaded struct {
		Data service.ContentModerationConfigView `json:"data"`
	}
	require.NoError(t, json.Unmarshal(result.Body.Bytes(), &loaded))
	require.Equal(t, updated.Data.WhitelistedUserIDs, loaded.Data.WhitelistedUserIDs)
}

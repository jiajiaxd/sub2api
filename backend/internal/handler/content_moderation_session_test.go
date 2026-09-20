//go:build unit

package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestContentModerationSessionHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, header := range []string{
		"x-opencode-session", "session-id", "session_id", "conversation_id",
		"X-Session-Affinity", "X-Session-Id", "X-Conversation-ID", "X-Claude-Code-Session-Id",
	} {
		t.Run(header, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
			c.Request.Header.Set(header, " conversation-1 ")
			input := buildContentModerationInput(c, nil, middleware2.AuthSubject{UserID: 42}, "anthropic", "model", nil)
			require.Equal(t, "conversation-1", input.AuditSessionID)
			require.EqualValues(t, 42, input.UserID)

			c.Request.Header.Set("x-opencode-session", "explicit-session")
			input = buildContentModerationInput(c, nil, middleware2.AuthSubject{UserID: 42}, "anthropic", "model", nil)
			require.Equal(t, "explicit-session", input.AuditSessionID)
		})
	}
}

func TestContentModerationMissingSessionPreservesTenantFallback(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	c.Request.Header.Set("X-Request-Id", "changes-every-turn")
	c.Request.Header.Set("x-audit-tenant", "forged-user")
	c.Request.Header.Set("session_id", "invalid\tvalue")
	input := buildContentModerationInput(c, nil, middleware2.AuthSubject{UserID: 42}, "anthropic", "model", nil)
	require.Empty(t, input.AuditSessionID)
	require.EqualValues(t, 42, input.UserID)
}

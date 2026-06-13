package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Completes before timeout", func(t *testing.T) {
		r := gin.New()
		r.Use(Timeout(10 * time.Millisecond))
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Times out", func(t *testing.T) {
		r := gin.New()
		r.Use(Timeout(1 * time.Millisecond))
		r.GET("/test", func(c *gin.Context) {
			time.Sleep(10 * time.Millisecond)

			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		strBody := w.Body.String()
		assert.NotEqual(t, http.StatusOK, w.Code)
		assert.NotContains(t, strBody, "ok")
		assert.Equal(t, http.StatusRequestTimeout, w.Code)
		assert.Contains(t, strBody, "request timeout")
	})
}

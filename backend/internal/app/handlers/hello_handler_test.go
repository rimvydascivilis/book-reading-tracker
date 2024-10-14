package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/hello", HelloWorld)

	req, err := http.NewRequest(http.MethodGet, "/hello", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	expectedResponse := `{"message":"Hello, World!"}`
	assert.JSONEq(t, expectedResponse, resp.Body.String())
}

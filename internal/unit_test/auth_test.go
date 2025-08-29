package unit_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/aryadisastra/authentication/internal/dto"
	"github.com/aryadisastra/authentication/internal/router"
)

type envelope struct {
	Result  bool            `json:"result"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "postgres://postgres:123@localhost:5432/auth_db?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("DELETE FROM tb_m_user;")
	return db
}

func TestRegisterLoginMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	r := router.New(db, "eeee1234qq11.", 60)

	regBody, _ := json.Marshal(dto.RegisterRequest{
		Username: "test", Email: "test@test.com", Name: "Test", Password: "password123", RoleCode: "user",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	loginBody, _ := json.Marshal(dto.LoginRequest{Identifier: "test", Password: "password123"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var env envelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &env))
	require.True(t, env.Result)

	var lr dto.LoginResponse
	require.NoError(t, json.Unmarshal(env.Data, &lr))
	require.NotEmpty(t, lr.Token)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+lr.Token)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

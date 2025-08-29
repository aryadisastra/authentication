package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/aryadisastra/authentication/internal/dto"
	"github.com/aryadisastra/authentication/internal/httpx"
	"github.com/aryadisastra/authentication/internal/models"
)

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			httpx.Fail(c, http.StatusBadRequest, err.Error())
			return
		}

		var role models.Role
		if err := db.Where("code = ?", req.RoleCode).First(&role).Error; err != nil {
			httpx.Fail(c, http.StatusBadRequest, "invalid role_code")
			return
		}

		hashed, err := HashPassword(req.Password)
		if err != nil {
			httpx.Fail(c, http.StatusInternalServerError, "failed to hash password")
			return
		}

		u := models.User{
			Username: req.Username,
			Email:    req.Email,
			Name:     req.Name,
			Password: hashed,
			RoleID:   role.ID,
		}

		if err := db.Create(&u).Error; err != nil {
			httpx.Fail(c, http.StatusBadRequest, "username/email already exists")
			return
		}

		httpx.Created(c, dto.ProfileResponse{
			ID: u.ID, Username: u.Username, Email: u.Email,
			RoleID: role.ID, RoleCode: role.Code, RoleName: role.Name,
		})
	}
}

func Login(db *gorm.DB, jwtSecret string, ttlMinutes int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			httpx.Fail(c, http.StatusUnauthorized, "invalid credentials")
			return
		}

		var u models.User
		if err := db.Preload("Role").
			Where("username = ? OR email = ?", req.Identifier, req.Identifier).
			First(&u).Error; err != nil {
			httpx.Fail(c, http.StatusUnauthorized, "invalid credentials")
			return
		}

		ok, err := CheckPassword(req.Password, u.Password)
		if err != nil || !ok {
			httpx.Fail(c, http.StatusUnauthorized, "invalid credentials")
			return
		}

		claims := jwt.MapClaims{
			"sub": u.ID, "role_id": u.RoleID, "role_code": u.Role.Code, "role_name": u.Role.Name,
			"exp": time.Now().Add(time.Duration(ttlMinutes) * time.Minute).Unix(),
			"iat": time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, _ := token.SignedString([]byte(jwtSecret))
		httpx.OK(c, dto.LoginResponse{Token: signed})
	}
}

func Me(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, _ := c.Get("user_id")
		var u models.User
		if err := db.Preload("Role").First(&u, "id = ?", uid).Error; err != nil {
			httpx.Fail(c, http.StatusUnauthorized, "user not found")
			return
		}
		httpx.OK(c, dto.ProfileResponse{
			ID:       u.ID,
			Username: u.Username,
			Name:     u.Name,
			Email:    u.Email,
			RoleID:   u.RoleID,
			RoleCode: u.Role.Code,
			RoleName: u.Role.Name,
		})
	}
}

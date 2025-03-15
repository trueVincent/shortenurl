package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	InitializeDatabase()
	code := m.Run()
	os.Exit(code)
}

func TestRegister_Success(t *testing.T) {
	username := "testuser"
	password := "testpass"

	user, err := Register(username, password)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
	assert.NotEmpty(t, user.Password)
}

func TestRegister_UsernameExists(t *testing.T) {
	_, err := Register("testuser", "testpass")
	assert.EqualError(t, err, "username already exists")
}

func TestLogin_Success(t *testing.T) {
	username := "testuser"
	password := "testpass"

	user, err := Login(username, password)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
}

func TestLogin_InvalidPassword(t *testing.T) {
	_, err := Login("testuser", "wrongpass")
	assert.Error(t, err)
}

func TestLogin_UserNotFound(t *testing.T) {
	_, err := Login("unknown", "testpass")
	assert.Error(t, err)
}

func TestCreateUrlMapping_Success(t *testing.T) {
	userID := uint(1)
	originURL := "https://example.com"

	urlMapping, err := CreateUrlMapping(userID, originURL)
	assert.NoError(t, err)
	assert.NotEmpty(t, urlMapping.ID)
	assert.Equal(t, originURL, urlMapping.OriginURL)
}

func TestRedirect_Success(t *testing.T) {
	var urlMappingID string
	DB.Model(&URLMapping{}).Select("ID").Take(&urlMappingID)

	urlMapping, err := Redirect(urlMappingID)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", urlMapping.OriginURL)
}

func TestRedirect_NotFound(t *testing.T) {
	_, err := Redirect("unknown")
	assert.Error(t, err)
}
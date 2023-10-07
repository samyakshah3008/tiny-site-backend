package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"tiny-site-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	configFilePath := "../../app.env"

	if err := godotenv.Load(configFilePath); err != nil {
		panic(fmt.Sprintf("Error loading configuration file: %v", err))
	}

	dsn := os.Getenv("TEST_DB_URL")

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	code := m.Run()

	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	os.Exit(code)
}

func TestGetUsers(t *testing.T) {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}

	router := gin.Default()
	routes.SetupRoutes(router)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "4231c0e9-e7d7-4633-bf4a-a5b196f4ff7d",
	})

	jwtSecret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/users/self", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()

	if resp == nil {
		t.Fatal("Response is nil")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Response Body:", string(body))

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imnotjin/fetch-backend-challenge/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

// TestMain sets up the testing environment, including initializing the database
// and running any necessary migrations. It ensures the test database is properly
// cleaned up after the tests are executed.
func TestMain(m *testing.M) {
	// Setup test database
	var err error
	testDB, err = setupTestDB()
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the schema
	err = testDB.AutoMigrate(&models.Transaction{})
	if err != nil {
		fmt.Printf("Failed to run migrations on test database: %v\n", err)
		os.Exit(1)
	}

	// Run the tests
	code := m.Run()

	// Clean up the database connection
	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	os.Exit(code)
}

// setupTestDB initializes a connection to the test PostgreSQL database using the
// connection details from the .env.test file. It returns a pointer to the gorm.DB instance.
func setupTestDB() (*gorm.DB, error) {
	// Load .env.test file
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}

	// Build DSN (Data Source Name) string for the test database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Open the database connection using gorm
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// setupTestRouter creates and configures the Gin router for handling API requests
// during testing. It registers the relevant routes for adding, spending, and getting
// balance points.
func setupTestRouter() *gin.Engine {
	r := gin.Default()
	h := &Handler{DB: testDB}

	r.POST("/add", h.AddPoints)
	r.POST("/spend", h.SpendPoints)
	r.GET("/balance", h.GetBalance)

	return r
}

// clearTestData truncates the transactions table in the test database to ensure
// a clean state before running each test case.
func clearTestData() {
	if testDB != nil {
		testDB.Exec("TRUNCATE TABLE transactions")
	} else {
		log.Fatal("testDB is not initialized")
	}
}

// TestIntegration tests the entire integration flow of adding points, spending points,
// and retrieving the final balance. It is divided into three subtests that validate the
// correctness of each endpoint.
func TestIntegration(t *testing.T) {
	clearTestData()
	router := setupTestRouter()

	// Test adding points
	t.Run("TestAddPoints", func(t *testing.T) {
		addPointsRequests := []AddPointsRequest{
			{Payer: "DANNON", Points: 300, Timestamp: time.Now().Add(-5 * time.Minute)},
			{Payer: "UNILEVER", Points: 200, Timestamp: time.Now().Add(-4 * time.Minute)},
			{Payer: "DANNON", Points: -200, Timestamp: time.Now().Add(-3 * time.Minute)},
			{Payer: "MILLER COORS", Points: 10000, Timestamp: time.Now().Add(-2 * time.Minute)},
			{Payer: "DANNON", Points: 1000, Timestamp: time.Now().Add(-1 * time.Minute)},
		}

		// Run subtests for each AddPointsRequest
		for _, req := range addPointsRequests {
			t.Run(fmt.Sprintf("Add: %s %+d", req.Payer, req.Points), func(t *testing.T) {
				jsonValue, _ := json.Marshal(req)
				httpReq, _ := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonValue))
				w := httptest.NewRecorder()
				router.ServeHTTP(w, httpReq)
				assert.Equal(t, http.StatusOK, w.Code)
			})
		}
	})

	// Test spending points
	t.Run("TestSpendPoints", func(t *testing.T) {
		spendRequest := SpendPointsRequest{Points: 5000}
		t.Run(fmt.Sprintf("Spend: %d", spendRequest.Points), func(t *testing.T) {
			jsonValue, _ := json.Marshal(spendRequest)
			req, _ := http.NewRequest("POST", "/spend", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			// Validate the response matches expected spending breakdown
			var response []SpendPointsResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			expectedResponse := []SpendPointsResponse{
				{Payer: "DANNON", Points: -100},
				{Payer: "UNILEVER", Points: -200},
				{Payer: "MILLER COORS", Points: -4700},
			}

			assert.ElementsMatch(t, expectedResponse, response)
		})
	})

	// Test retrieving balance
	t.Run("TestGetBalance", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/balance", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Validate the final balance
		var balance map[string]int
		err := json.Unmarshal(w.Body.Bytes(), &balance)
		assert.NoError(t, err)

		expectedBalance := map[string]int{
			"DANNON":       1000,
			"UNILEVER":     0,
			"MILLER COORS": 5300,
		}

		assert.Equal(t, expectedBalance, balance)
	})
}

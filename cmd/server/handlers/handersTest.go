package handlers

import (
	"arcRaidersAPI/cmd/sqlfuncs"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetAllItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	conn, err := sqlfuncs.ConnectToDB()

	if err != nil {
		log.Fatalf("Error connecting to the database")
	}

	defer func() {
		if err := sqlfuncs.DisconnectDB(conn); err != nil {
			log.Printf("Unable to disconnect from database: %v", err)
		}
	}()

	router := gin.New()
	router.GET("items", func(c *gin.Context) {
		GetItems(c, conn)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var items []sqlfuncs.Item
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

}

package controllers

import (
	"Footware-Ecommerce/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestSignup(t *testing.T) {
	r := SetUpRouter()
	r.POST("/users/signup", Signup())
	id := primitive.NewObjectID()
	// var user models.User
	userDetails := models.User{
		ID:         id,
		First_Name: "Roopam",
		Last_Name:  "Hedoo",
		Email:      "roopam@gmail.com",
		Password:   "Roopam9029@",
		Phone:      "+919234727484",
		User_type:  "USER",
	}
	jsonValue, _ := json.Marshal(userDetails)
	req, _ := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	println(w.Code)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin(t *testing.T) {
	r := SetUpRouter()
	r.POST("/users/login", Login())
	userDetails := models.User{
		Email:    "roopam@gmail.com",
		Password: "Roopam9029@",
	}
	jsonValue, _ := json.Marshal(userDetails)
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	println(w.Code)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductSearch(t *testing.T) {
	r := SetUpRouter()
	r.GET("/users/productview", AllProductList())
	req, _ := http.NewRequest("GET", "/users/productview", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var productlist []models.ProductUser
	json.Unmarshal(w.Body.Bytes(), &productlist)

	assert.Equal(t, http.StatusOK, w.Code)
}

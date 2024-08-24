package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UpdateTodo struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      bool   `json:"status" validate:"required"`
	UpdateOn 	string `json:"updatedOn"`
	CreatedOn 	string `json:"createdOn"`
}

type CreateTodo struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	CreatedOn 	string `json:"createdOn"`
}

var validate *validator.Validate = validator.New()

func CreateHandler(c *gin.Context) {
	var createTodo CreateTodo

	err := c.BindJSON(&createTodo)
	if err != nil {
		clientError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = validate.Struct(&createTodo)
	if err != nil {
		clientError(c, http.StatusBadRequest, err)
		return
	}

	newTodo, err := insertItem(c, createTodo)

	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(201, newTodo)
	
}

func GetAllHandler(c *gin.Context) {
	todos, err := listItems(c)
	if err != nil {
		serverError(c, err)
	}

	c.JSON(200, todos)

}

func GetHandler(c *gin.Context) {
	log.Println("Get todo")

	id := c.Param("id")
	if len(id) == 0 {
		clientError(c, http.StatusBadRequest, fmt.Errorf("id is required"))
	}
	
	todo, err := getItem(c, id)
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(200, todo)
}

func UpdateHandler(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		clientError(c, http.StatusBadRequest, fmt.Errorf("id is required"))
	}

	var updateTodo UpdateTodo

	err := c.BindJSON(&updateTodo)
	if err != nil {
		clientError(c, http.StatusUnprocessableEntity, err)
		return
	}

	err = validate.Struct(updateTodo)
	if err != nil {
		clientError(c, http.StatusBadRequest, fmt.Errorf("invalid body"))
	}

	todo, err := updateItem(c, id, updateTodo)

	if err != nil {
		serverError(c, err)
		return
	}

	if todo == nil {
		clientError(c, http.StatusNotFound, fmt.Errorf("todo not found"))
		return
	}

	c.JSON(200, todo)

}

func DeleteHandler(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		clientError(c, http.StatusBadRequest, fmt.Errorf("id is required"))
	}

	todo, err := deleteItem(c, id)
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(200, todo)

}

func serverError(c *gin.Context, err error) {
	c.AbortWithError(http.StatusInternalServerError, err)
}

func clientError(c *gin.Context, status int, err error) {
	c.AbortWithError(status, err)
}
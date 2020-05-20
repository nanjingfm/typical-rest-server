package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/typical-go/typical-rest-server/internal/server/repository"
	"github.com/typical-go/typical-rest-server/internal/server/service"
	"go.uber.org/dig"
)

type (

	// BookCntrl is controller to book entity
	BookCntrl struct {
		dig.In
		service.BookService
	}
)

// SetRoute to define API Route
func (c *BookCntrl) SetRoute(e *echo.Echo) {
	e.GET("books", c.Find)
	e.POST("books", c.Create)
	e.GET("books/:id", c.FindOne)
	e.PUT("books/:id", c.Update)
	e.DELETE("books/:id", c.Delete)
}

// Create book
func (c *BookCntrl) Create(ec echo.Context) (err error) {
	var book repository.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}

	ctx := ec.Request().Context()
	id, err := c.BookService.Create(ctx, &book)

	if err != nil {
		return httpError(err)
	}

	ec.Response().Header().Set(echo.HeaderLocation, fmt.Sprintf("/books/%d", id))
	return ec.NoContent(http.StatusCreated)
}

// Find books
func (c *BookCntrl) Find(ec echo.Context) (err error) {
	var books []*repository.Book
	if books, err = c.BookService.Find(
		ec.Request().Context(),
	); err != nil {
		return httpError(err)
	}
	return ec.JSON(http.StatusOK, books)
}

// FindOne book
func (c *BookCntrl) FindOne(ec echo.Context) error {
	book, err := c.BookService.FindOne(
		ec.Request().Context(),
		ec.Param("id"),
	)

	if err != nil {
		return httpError(err)
	}

	return ec.JSON(http.StatusOK, book)
}

// Delete book
func (c *BookCntrl) Delete(ec echo.Context) (err error) {
	if err = c.BookService.Delete(
		ec.Request().Context(),
		ec.Param("id"),
	); err != nil {
		return httpError(err)
	}
	return ec.NoContent(http.StatusNoContent)
}

// Update book
func (c *BookCntrl) Update(ec echo.Context) (err error) {
	var book repository.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}

	if err = c.BookService.Update(
		ec.Request().Context(),
		ec.Param("id"),
		&book,
	); err != nil {
		return httpError(err)
	}

	return ec.NoContent(http.StatusOK)
}
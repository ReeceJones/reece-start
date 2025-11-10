package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

// Test structs for Validated function
type TestRequestBody struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"min=0,max=150"`
}

type TestRequestBodyOptional struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"omitempty,email"`
}

// Test structs for ValidatedQuery function
type TestQueryParams struct {
	Page   int    `query:"page" validate:"min=1"`
	Limit  int    `query:"limit" validate:"min=1,max=100"`
	Search string `query:"search"`
	SortBy string `query:"sort_by" validate:"oneof=name email created_at"`
}

type TestQueryParamsRequired struct {
	Page int `query:"page" validate:"required,min=1"`
}

// Test structs for ValidatedWithQuery function
type TestBodyWithQuery struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

type TestQueryWithBody struct {
	Filter string `query:"filter" validate:"required"`
}

func TestValidated(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		e := echo.New()
		reqBody := TestRequestBody{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   30,
		}
		bodyJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := Validated(func(c echo.Context, req TestRequestBody) error {
			handlerCalled = true
			require.Equal(t, "John Doe", req.Name)
			require.Equal(t, "john@example.com", req.Email)
			require.Equal(t, 30, req.Age)
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})

		err := handler(c)
		require.NoError(t, err)
		require.True(t, handlerCalled)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := Validated(func(c echo.Context, req TestRequestBody) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		e := echo.New()
		reqBody := TestRequestBody{
			Email: "john@example.com",
			Age:   30,
			// Name is missing
		}
		bodyJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := Validated(func(c echo.Context, req TestRequestBody) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		e := echo.New()
		reqBody := TestRequestBody{
			Name:  "John Doe",
			Email: "invalid-email",
			Age:   30,
		}
		bodyJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := Validated(func(c echo.Context, req TestRequestBody) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("InvalidAgeRange", func(t *testing.T) {
		e := echo.New()
		reqBody := TestRequestBody{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   200, // Exceeds max of 150
		}
		bodyJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := Validated(func(c echo.Context, req TestRequestBody) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("EmptyBody", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte("{}")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := Validated(func(c echo.Context, req TestRequestBody) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("OptionalFields", func(t *testing.T) {
		e := echo.New()
		reqBody := TestRequestBodyOptional{
			Name: "John Doe",
			// Email is optional and omitted
		}
		bodyJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := Validated(func(c echo.Context, req TestRequestBodyOptional) error {
			handlerCalled = true
			require.Equal(t, "John Doe", req.Name)
			require.Equal(t, "", req.Email)
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})

		err := handler(c)
		require.NoError(t, err)
		require.True(t, handlerCalled)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("OptionalFieldWithInvalidEmail", func(t *testing.T) {
		e := echo.New()
		reqBody := TestRequestBodyOptional{
			Name:  "John Doe",
			Email: "invalid-email", // Email is provided but invalid
		}
		bodyJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := Validated(func(c echo.Context, req TestRequestBodyOptional) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})
}

func TestValidatedQuery(t *testing.T) {
	t.Run("ValidQuery", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test?page=1&limit=10&sort_by=name", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := ValidatedQuery(func(c echo.Context, q TestQueryParams) error {
			handlerCalled = true
			require.Equal(t, 1, q.Page)
			require.Equal(t, 10, q.Limit)
			require.Equal(t, "name", q.SortBy)
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})

		err := handler(c)
		require.NoError(t, err)
		require.True(t, handlerCalled)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidPageMin", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test?page=0&limit=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := ValidatedQuery(func(c echo.Context, q TestQueryParams) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("InvalidLimitMax", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test?page=1&limit=200", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := ValidatedQuery(func(c echo.Context, q TestQueryParams) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("InvalidOneOf", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test?page=1&sort_by=invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := ValidatedQuery(func(c echo.Context, q TestQueryParams) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := ValidatedQuery(func(c echo.Context, q TestQueryParamsRequired) error {
			handlerCalled = true
			return nil
		})

		err := handler(c)
		require.Error(t, err)
		require.False(t, handlerCalled)

		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("EmptyQuery", func(t *testing.T) {
		// Test struct without min validation for optional fields
		type OptionalQueryParams struct {
			Page   int    `query:"page"`
			Limit  int    `query:"limit"`
			Search string `query:"search"`
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := ValidatedQuery(func(c echo.Context, q OptionalQueryParams) error {
			handlerCalled = true
			require.Equal(t, 0, q.Page)
			require.Equal(t, 0, q.Limit)
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})

		err := handler(c)
		require.NoError(t, err)
		require.True(t, handlerCalled)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}

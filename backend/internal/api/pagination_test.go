package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestEncodeCursor(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		cursor := map[string]string{
			"user_id":   "1",
			"direction": "next",
		}
		encoded, err := EncodeCursor(cursor)
		require.NoError(t, err)
		require.Equal(t, "eyJkaXJlY3Rpb24iOiJuZXh0IiwidXNlcl9pZCI6IjEifQ==", encoded)
	})

	t.Run("String", func(t *testing.T) {
		cursor := "test"
		encoded, err := EncodeCursor(cursor)
		require.NoError(t, err)
		require.Equal(t, "InRlc3Qi", encoded)
	})
}

func TestParseCursor(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		cursor := "eyJkaXJlY3Rpb24iOiJuZXh0IiwidXNlcl9pZCI6IjEifQ=="
		var result map[string]string
		err := ParseCursor(cursor, &result)
		require.NoError(t, err)
		require.Equal(t, map[string]string{
			"direction": "next",
			"user_id":   "1",
		}, result)
	})

	t.Run("String", func(t *testing.T) {
		cursor := "InRlc3Qi"
		var result string
		err := ParseCursor(cursor, &result)
		require.NoError(t, err)
		require.Equal(t, "test", result)
	})
}

func TestBuildPaginationLinks(t *testing.T) {
	t.Run("NoCursors", func(t *testing.T) {
		requiest := httptest.NewRequest(http.MethodGet, "/users", nil)
		requiest.URL.RawQuery = "page[cursor]=1"
		params := BuildPaginationLinksParams{
			Context: echo.New().NewContext(requiest, nil),
		}
		links := BuildPaginationLinks(params)
		require.Equal(t, PaginationLinks{}, links)
	})

	t.Run("WithPrevCursor", func(t *testing.T) {
		requiest := httptest.NewRequest(http.MethodGet, "/users", nil)
		requiest.URL.RawQuery = "page[cursor]=0"
		params := BuildPaginationLinksParams{
			PrevCursor: "1",
			Context:    echo.New().NewContext(requiest, nil),
		}
		links := BuildPaginationLinks(params)
		require.Equal(t, PaginationLinks{
			Prev: "/users?page%5Bcursor%5D=1",
		}, links)
	})

	t.Run("WithNextCursor", func(t *testing.T) {
		requiest := httptest.NewRequest(http.MethodGet, "/users", nil)
		requiest.URL.RawQuery = "page[cursor]=0"
		params := BuildPaginationLinksParams{
			NextCursor: "1",
			Context:    echo.New().NewContext(requiest, nil),
		}
		links := BuildPaginationLinks(params)
		require.Equal(t, PaginationLinks{
			Next: "/users?page%5Bcursor%5D=1",
		}, links)
	})
}

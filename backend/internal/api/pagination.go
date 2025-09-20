package api

import (
	"encoding/base64"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type PaginationLinks struct {
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}

type BuildPaginationLinksParams struct {
	PrevCursor string
	NextCursor string
	Context echo.Context
}

func EncodeCursor[T any](cursor T) (string, error) {
	// json marshal first
	encoded, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encoded), nil
}

func ParseCursor[T any](cursor string, result *T) error {
	if cursor == "" {
		return nil
	}

	// base64 decode first
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return err
	}

	// then json unmarshal
	err = json.Unmarshal(decoded, result)
	if err != nil {
		return err
	}
	return nil
}

func BuildPaginationLinks(params BuildPaginationLinksParams) PaginationLinks {
	var next string
	var prev string

	queryParams := params.Context.Request().URL.Query()

	if params.PrevCursor != "" {
		queryParams.Set("page[cursor]", params.PrevCursor)
		prev = params.Context.Request().URL.Path + "?" + queryParams.Encode()
	}

	if params.NextCursor != "" {
		queryParams.Set("page[cursor]", params.NextCursor)
		next = params.Context.Request().URL.Path + "?" + queryParams.Encode()
	}

	links := PaginationLinks{
		Prev: prev,
		Next: next,
	}

	return links
}

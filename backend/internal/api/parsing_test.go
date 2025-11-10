package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestParseOrganizationIDFromParams(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/organizations/123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("123")

		orgID, err := ParseOrganizationIDFromParams(c)
		require.NoError(t, err)
		require.Equal(t, uint(123), orgID)
	})

	t.Run("InvalidNonNumeric", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/organizations/abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("abc")

		orgID, err := ParseOrganizationIDFromParams(c)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidOrganizationID))
		require.Equal(t, uint(0), orgID)
	})

	t.Run("InvalidEmpty", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/organizations/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		orgID, err := ParseOrganizationIDFromParams(c)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidOrganizationID))
		require.Equal(t, uint(0), orgID)
	})

	t.Run("InvalidNegative", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/organizations/-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("-1")

		orgID, err := ParseOrganizationIDFromParams(c)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidOrganizationID))
		require.Equal(t, uint(0), orgID)
	})

	t.Run("ValidZero", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/organizations/0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("0")

		orgID, err := ParseOrganizationIDFromParams(c)
		require.NoError(t, err)
		require.Equal(t, uint(0), orgID)
	})
}

func TestParseOrganizationIDFromString(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		orgID, err := ParseOrganizationIDFromString("456")
		require.NoError(t, err)
		require.Equal(t, uint(456), orgID)
	})

	t.Run("InvalidNonNumeric", func(t *testing.T) {
		orgID, err := ParseOrganizationIDFromString("xyz")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidOrganizationID))
		require.Equal(t, uint(0), orgID)
	})

	t.Run("InvalidEmpty", func(t *testing.T) {
		orgID, err := ParseOrganizationIDFromString("")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidOrganizationID))
		require.Equal(t, uint(0), orgID)
	})

	t.Run("InvalidNegative", func(t *testing.T) {
		orgID, err := ParseOrganizationIDFromString("-5")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidOrganizationID))
		require.Equal(t, uint(0), orgID)
	})

	t.Run("ValidZero", func(t *testing.T) {
		orgID, err := ParseOrganizationIDFromString("0")
		require.NoError(t, err)
		require.Equal(t, uint(0), orgID)
	})

	t.Run("ValidMaxUint32", func(t *testing.T) {
		orgID, err := ParseOrganizationIDFromString("4294967295")
		require.NoError(t, err)
		require.Equal(t, uint(4294967295), orgID)
	})
}

func TestParseUserIDFromString(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		userID, err := ParseUserIDFromString("789")
		require.NoError(t, err)
		require.Equal(t, uint(789), userID)
	})

	t.Run("InvalidNonNumeric", func(t *testing.T) {
		userID, err := ParseUserIDFromString("invalid")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidUserID))
		require.Equal(t, uint(0), userID)
	})

	t.Run("InvalidEmpty", func(t *testing.T) {
		userID, err := ParseUserIDFromString("")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidUserID))
		require.Equal(t, uint(0), userID)
	})

	t.Run("InvalidNegative", func(t *testing.T) {
		userID, err := ParseUserIDFromString("-10")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidUserID))
		require.Equal(t, uint(0), userID)
	})

	t.Run("ValidZero", func(t *testing.T) {
		userID, err := ParseUserIDFromString("0")
		require.NoError(t, err)
		require.Equal(t, uint(0), userID)
	})
}

func TestParseMembershipIDFromParams(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/memberships/42", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("42")

		membershipID, err := ParseMembershipIDFromParams(c)
		require.NoError(t, err)
		require.Equal(t, uint(42), membershipID)
	})

	t.Run("InvalidNonNumeric", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/memberships/bad", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")

		membershipID, err := ParseMembershipIDFromParams(c)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidMembershipID))
		require.Equal(t, uint(0), membershipID)
	})

	t.Run("InvalidEmpty", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/memberships/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		membershipID, err := ParseMembershipIDFromParams(c)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidMembershipID))
		require.Equal(t, uint(0), membershipID)
	})

	t.Run("InvalidNegative", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/memberships/-99", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("-99")

		membershipID, err := ParseMembershipIDFromParams(c)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidMembershipID))
		require.Equal(t, uint(0), membershipID)
	})

	t.Run("ValidZero", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/memberships/0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("0")

		membershipID, err := ParseMembershipIDFromParams(c)
		require.NoError(t, err)
		require.Equal(t, uint(0), membershipID)
	})
}

func TestParseOrganizationInvitationIDFromParams(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		validUUID := uuid.New()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/invitations/"+validUUID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(validUUID.String())

		invitationID, err := ParseOrganizationInvitationIDFromParams(c)
		require.NoError(t, err)
		require.Equal(t, validUUID, invitationID)
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/invitations/not-a-uuid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("not-a-uuid")

		invitationID, err := ParseOrganizationInvitationIDFromParams(c)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidInvitationID))
		require.Equal(t, uuid.Nil, invitationID)
	})

	t.Run("InvalidEmpty", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/invitations/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		invitationID, err := ParseOrganizationInvitationIDFromParams(c)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidInvitationID))
		require.Equal(t, uuid.Nil, invitationID)
	})

	t.Run("InvalidPartialUUID", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/invitations/123e4567-e89b-12d3", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("123e4567-e89b-12d3")

		invitationID, err := ParseOrganizationInvitationIDFromParams(c)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidInvitationID))
		require.Equal(t, uuid.Nil, invitationID)
	})

	t.Run("ValidNilUUID", func(t *testing.T) {
		nilUUID := uuid.Nil
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/invitations/"+nilUUID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(nilUUID.String())

		invitationID, err := ParseOrganizationInvitationIDFromParams(c)
		require.NoError(t, err)
		require.Equal(t, uuid.Nil, invitationID)
	})
}

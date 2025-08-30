package constants

type ErrorCode string

const (
	ErrorCodeInvalidRequest ErrorCode = "invalid_request"
	ErrorCodeValidationFailed ErrorCode = "validation_failed"
	ErrorCodeInternalServerError ErrorCode = "internal_server_error"
	ErrorCodeBadRequest ErrorCode = "bad_request"
	ErrorCodeUnauthorized ErrorCode = "unauthorized"
	ErrorCodeNotFound ErrorCode = "not_found"
	ErrorCodeForbidden ErrorCode = "forbidden"
	ErrorCodeConflict ErrorCode = "conflict"
)

type ApiType string

const (
	ApiTypeUser ApiType = "user"
	ApiTypeOrganization ApiType = "organization"
	ApiTypeOrganizationMembership ApiType = "organization-membership"
	ApiTypeOrganizationInvitation ApiType = "organization-invitation"
)

package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUniqueConstraintViolation(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		err := errors.New("pq: duplicate key value violates unique constraint \"users_email_key\"")
		assert.True(t, IsUniqueConstraintViolation(err))
	})

	t.Run("False", func(t *testing.T) {
		err := errors.New("pq: insert or update on table \"users\" violates foreign key constraint \"users_organization_id_fkey\"")
		assert.False(t, IsUniqueConstraintViolation(err))
	})
}

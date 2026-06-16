package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateProjectStatus_Validation(t *testing.T) {
	// We only test the validation part which doesn't hit the DB for invalid statuses
	t.Run("invalid status", func(t *testing.T) {
		p, err := UpdateProjectStatus("test-id", "InvalidStatus")
		
		assert.Error(t, err)
		assert.Equal(t, "invalid status: must be Backlog, Em andamento, or Go-Live", err.Error())
		assert.Nil(t, p)
	})
}

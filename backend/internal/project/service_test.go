package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateProjectStatus_Validation(t *testing.T) {
	t.Run("invalid status", func(t *testing.T) {
		p, err := UpdateProjectStatus("test-id", "InvalidStatus")
		
		assert.Error(t, err)
		assert.Equal(t, "invalid status: must be Backlog, Em andamento, or Go-Live", err.Error())
		assert.Nil(t, p)
	})

	t.Run("empty status", func(t *testing.T) {
		p, err := UpdateProjectStatus("test-id", "")
		
		assert.Error(t, err)
		assert.Equal(t, "invalid status: must be Backlog, Em andamento, or Go-Live", err.Error())
		assert.Nil(t, p)
	})

	t.Run("case sensitive status", func(t *testing.T) {
		// "backlog" != "Backlog" (case sensitive)
		p, err := UpdateProjectStatus("test-id", "backlog")
		
		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestCreateProject_EmptyName(t *testing.T) {
	p, err := CreateProject("client-1", "")

	assert.Error(t, err)
	assert.Equal(t, "project name cannot be empty", err.Error())
	assert.Nil(t, p)
}

package user

import (
	"database/sql"
	"errors"

	"onboardly-backend/internal/auth"
	"onboardly-backend/internal/db"
)

// ListUsers returns a list of all users.
func ListUsers() ([]auth.User, error) {
	query := `SELECT id, email, role, created_at FROM users ORDER BY created_at DESC`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []auth.User
	for rows.Next() {
		var u auth.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteUser removes a user by ID. It prevents deleting the last Admin.
func DeleteUser(id string, currentUserID string) error {
	// 1. Prevents deleting the current user
	if id == currentUserID {
		return errors.New("cannot delete your own account")
	}

	// 2. Check if the user is an Admin and if it's the last one
	var role string
	err := db.DB.QueryRow(`SELECT role FROM users WHERE id = $1`, id).Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	if role == "Admin" {
		var adminCount int
		err = db.DB.QueryRow(`SELECT count(*) FROM users WHERE role = 'Admin'`).Scan(&adminCount)
		if err != nil {
			return err
		}
		if adminCount <= 1 {
			return errors.New("cannot delete the last admin in the system")
		}
	}

	// 3. Delete the user
	result, err := db.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

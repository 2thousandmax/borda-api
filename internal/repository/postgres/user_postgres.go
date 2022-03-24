package postgres

import (
	"borda/internal/domain"
	"database/sql"
	"errors"
	"strings"

	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// TODO: pass user object when create user
func (r UserRepository) SaveUser(username, password, contact string) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return -1, err
	}

	isUserExistQuery := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE name=$1
			LIMIT 1
		)`,
		userTable,
	)

	var isUserExist bool
	if err := tx.Get(&isUserExist, isUserExistQuery, username); err != nil {
		return -1, err
	}

	if isUserExist {
		return -1, NewErrAlreadyExist("user", "username", username)
	}

	// Save user to the database
	createUserQuery := fmt.Sprintf(`
		INSERT INTO public.%s (
			name,
			password,
			contact
		)
		VALUES($1, $2, $3)
		RETURNING id`,
		userTable,
	)

	var userId int
	if err := tx.Get(&userId, createUserQuery, username, password, contact); err != nil {
		return -1, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return -1, err
	}

	return userId, nil
}

func (r UserRepository) GetUserByCredentials(username, password string) (*domain.User, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM public.%s
		WHERE name=$1 AND password=$2`,
		userTable,
	)

	var user domain.User
	if err := r.db.Get(&user, query, username, password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrNotFound("user", "username, password",
				strings.Join([]string{username, password}, ", "))
		}
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) UpdatePassword(userId int, newPassword string) error {
	query := fmt.Sprintf(`
		UPDATE public.%s
		SET password = $1
		WHERE id = $2`,
		userTable,
	)

	if _, err := r.db.Exec(query, newPassword, userId); err != nil {
		return err
	}

	return nil
}

func (r UserRepository) AssignRole(userId, roleId int) error {
	query := fmt.Sprintf(`
		INSERT INTO public.%s (
			user_id,
			role_id
		)
		VALUES($1, $2)`,
		userRolesTable)

	if _, err := r.db.Exec(query, userId, roleId); err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetUserRole(userId int) (*domain.Role, error) {
	query := fmt.Sprintf(`
		SELECT r.id, r.name
		FROM public.%s AS r
		INNER JOIN public.%s AS ur ON r.id=ur.role_id
		WHERE ur.user_id = $1;`,
		roleTable, userRolesTable)

	var role domain.Role
	if err := r.db.Get(&role, query, userId); err != nil {
		return nil, err
	}

	return &role, nil
}

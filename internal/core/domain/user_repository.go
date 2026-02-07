package domain

import "context"

// UserRow represents a user record returned from the database.
// It includes the password hash so the Logic layer can verify credentials.
type UserRow struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
}

// UserRepository defines the data-access contract for user operations.
// Implementations live in internal/core/repository (Core layer).
// The Logic layer depends on this interface only — never on SQL or pgx directly.
type UserRepository interface {
	// GetByUsername returns the user matching the given username.
	// Returns (nil, nil) when no user is found.
	GetByUsername(ctx context.Context, username string) (*UserRow, error)

	// ExistsByUsernameOrEmail returns true when a user with the given
	// username or email already exists.
	ExistsByUsernameOrEmail(ctx context.Context, username, email string) (bool, error)

	// Create inserts a new user and returns the generated user ID.
	Create(ctx context.Context, username, email, passwordHash string) (int, error)

	// UpdateLastLogin sets the last_login timestamp to now for the given user.
	UpdateLastLogin(ctx context.Context, userID int) error
}

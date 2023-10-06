package password

// Hasher provides logic for password hashing and verifying.
type Hasher interface {
	// GenerateHashFromPassword is used to generate a hash for passed password.
	GenerateHashFromPassword(password string) (string, error)
	// CompareHashAndPassword is used for comparing hashed password with not hashed password.
	// If passwords are equal - returns true, otherwise - false.
	CompareHashAndPassword(opts *CompareHashAndPasswordOptions) (bool, error)
}

type CompareHashAndPasswordOptions struct {
	Hashed   string
	Password string
}

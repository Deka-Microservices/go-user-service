package bcrypt

import "context"

type BcryptClient interface {
	HashPassword(ctx context.Context, password string) (string, error)
	CheckPassword(ctx context.Context, password string, hash string) (bool, error)
	Close()
}

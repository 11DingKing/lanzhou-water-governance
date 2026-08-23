package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
)

type Auth struct {
	Users      repository.Users
	SessionTTL time.Duration
}

func NewAuth(users repository.Users, ttl time.Duration) Auth {
	return Auth{Users: users, SessionTTL: ttl}
}
func (s Auth) CreateUser(ctx context.Context, username, password string, role domain.Role, region string) (domain.User, error) {
	if username == "" || password == "" {
		return domain.User{}, fmt.Errorf("credentials required")
	}
	return s.Users.Create(ctx, username, hash(password), role, region)
}
func (s Auth) Login(ctx context.Context, username, password string) (string, domain.User, error) {
	stored, err := s.Users.Password(ctx, username)
	if err != nil {
		return "", domain.User{}, err
	}
	if stored != hash(password) {
		return "", domain.User{}, domain.ErrForbidden
	}
	user, err := s.Users.Find(ctx, username)
	if err != nil {
		return "", user, err
	}
	token := fmt.Sprintf("%x-%d", sha256.Sum256([]byte(fmt.Sprintf("%s-%d", username, time.Now().UnixNano()))), time.Now().UnixNano())
	if err = s.Users.CreateSession(ctx, user.ID, token, time.Now().UTC().Add(s.SessionTTL)); err != nil {
		return "", user, err
	}
	return token, user, nil
}
func (s Auth) Current(ctx context.Context, token string) (domain.User, error) {
	if token == "" {
		return domain.User{}, domain.ErrForbidden
	}
	user, err := s.Users.Session(ctx, token, time.Now().UTC())
	if err != nil { return domain.User{}, fmt.Errorf("session lookup: %w", err) }
	return user, nil
}
func (s Auth) Logout(ctx context.Context, token string) error {
	return s.Users.RevokeSession(ctx, token)
}
func hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

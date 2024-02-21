package session

import (
	"context"

	"github.com/google/uuid"
)

type (
	Repository interface {
		GetSession(ctx context.Context, sID uuid.UUID) (*Session, error)

		GetSessionByToken(ctx context.Context, token string) (*Session, error)

		UpsertSession(ctx context.Context, s *Session) error

		DeactivateSession(ctx context.Context, iID, sID uuid.UUID) error

		DeactivateSessionsFromIdentityExcept(ctx context.Context, iID, sID uuid.UUID) (int, error)
	}
)

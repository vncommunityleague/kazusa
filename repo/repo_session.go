package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/vncommunityleague/kazusa/session"
)

func (r *repositoryImpl) GetSession(ctx context.Context, sID uuid.UUID) (*session.Session, error) {
	var s session.Session
	if result := r.d.DB.WithContext(ctx).First(&s, sID); result.Error != nil {
		return nil, result.Error
	}

	return &s, nil
}

func (r *repositoryImpl) GetSessionByToken(ctx context.Context, token string) (*session.Session, error) {
	var s session.Session
	if result := r.d.DB.WithContext(ctx).InnerJoins("Identity", &session.Session{
		Token: token,
	}).First(&s); result.Error != nil {
		return nil, result.Error
	}

	return &s, nil
}

func (r *repositoryImpl) UpsertSession(ctx context.Context, s *session.Session) error {
	if result := r.d.DB.WithContext(ctx).Save(s); result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repositoryImpl) DeactivateSession(ctx context.Context, iID, sID uuid.UUID) error {
	var s session.Session
	if result := r.d.DB.WithContext(ctx).Model(&s).Where(&session.Session{
		ID:         sID,
		IdentityID: iID,
	}).Updates(session.Session{
		Active: false,
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repositoryImpl) DeactivateSessionsFromIdentityExcept(ctx context.Context, iID, sID uuid.UUID) (int, error) {
	var s session.Session
	result := r.d.DB.WithContext(ctx).Model(&s).Where("identity_id = ? AND id != ?", iID, sID).Updates(session.Session{
		Active: false,
	})
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

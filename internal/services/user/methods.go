package user

import (
	"context"
	"time"
	"tokeon-test-task/internal/models"
	"tokeon-test-task/internal/repos/users"
	"tokeon-test-task/pkg/utils"

	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user, ", err)
		return nil, err
	}

	return user, nil
}

func (s *service) Delete(ctx context.Context, user *models.User) (*models.User, error) {
	var updatePrams users.UpdateParams
	if err := utils.JsonToStruct(user, &updatePrams); err != nil {
		s.logger.Error("failed to update user, ", err)
		return nil, err
	}

	updatePrams.DeletedAt = utils.Pointer(time.Now())

	user, err := s.repo.Update(ctx, updatePrams)
	if err != nil {
		s.logger.Error("failed to update user, ", err)
		return nil, err
	}

	return user, nil
}

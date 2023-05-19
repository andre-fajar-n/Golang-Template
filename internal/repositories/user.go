package repositories

import (
	"context"
	"fmt"
	"go_template/gen/models"
	"go_template/runtime"
	"net/http"

	"gorm.io/gorm"
)

type (
	user struct {
		runtime runtime.Runtime
	}

	User interface {
		Create(ctx context.Context, tx *gorm.DB, data *models.User) (*models.User, error)
		FindBySingleColumn(ctx context.Context, column string, value interface{}, isDeleted bool) (*models.User, error)
		UsernameExist(ctx context.Context, username string) (bool, error)
	}
)

func Newuser(rt runtime.Runtime) User {
	return &user{
		rt,
	}
}

func (r *user) Create(ctx context.Context, tx *gorm.DB, data *models.User) (*models.User, error) {
	logger := r.runtime.Logger.With().
		Interface("data", data).
		Logger()

	if tx == nil {
		tx = r.runtime.Db
	}

	err := tx.Model(&data).Select("*").Create(&data).Error
	if err != nil {
		logger.Error().Err(err).Msg("error query")
		return nil, err
	}

	return data, nil
}

func (r *user) FindBySingleColumn(ctx context.Context, column string, value interface{}, isDeleted bool) (*models.User, error) {
	logger := r.runtime.Logger.With().
		Str("column", column).
		Interface("value", value).
		Logger()

	userModel := models.User{}
	db := r.runtime.Db.Model(&userModel).Where(fmt.Sprintf("%s = ?", column), value)

	if isDeleted {
		db = db.Where("deleted_at IS NOT NULL")
	} else {
		db = db.Where("deleted_at IS NULL")
	}

	err := db.First(&userModel).Error
	if err == gorm.ErrRecordNotFound {
		return nil, r.runtime.SetError(http.StatusNotFound, "user not found")
	}
	if err != nil {
		logger.Error().Err(err).Msg("error query")
		return nil, err
	}

	return &userModel, nil
}

func (r *user) UsernameExist(ctx context.Context, username string) (bool, error) {
	logger := r.runtime.Logger.With().
		Str("username", username).
		Logger()

	userModel := models.User{}

	db := r.runtime.Db.Model(&userModel).Where("username", username)

	err := db.First(&userModel).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		logger.Error().Err(err).Msg("error query")
		return false, err
	}

	return true, nil
}

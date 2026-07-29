package example

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, item Example) (Example, error) {
	err := r.db.WithContext(ctx).Create(&item).Error
	return item, err
}

func (r *Repository) FindAll(ctx context.Context, search string, limit, offset int32) ([]Example, int64, error) {
	var items []Example
	var total int64

	query := r.db.WithContext(ctx).Model(&Example{})
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(int(limit)).Offset(int(offset)).Find(&items).Error
	return items, total, err
}

func (r *Repository) FindByID(ctx context.Context, id string) (Example, error) {
	var item Example
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	return item, err
}

func (r *Repository) Update(ctx context.Context, item Example) (Example, error) {
	err := r.db.WithContext(ctx).Save(&item).Error
	return item, err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Example{}, "id = ?", id).Error
}

package postgres

import (
	"context"
	"errors"

	"veltiq/internal/core/domain"

	"gorm.io/gorm"
)

type ImportRepository struct {
	db *gorm.DB
}

func NewImportRepository(db *gorm.DB) *ImportRepository {
	return &ImportRepository{db: db}
}

func (r *ImportRepository) GetByID(ctx context.Context, id string) (domain.Import, error) {
	var record importRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Import{}, domain.ErrImportNotFound
		}
		return domain.Import{}, err
	}
	return importToDomain(record), nil
}

func (r *ImportRepository) GetStatusByID(ctx context.Context, id string) (domain.ImportStatus, error) {
	imp, err := r.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	return imp.Status, nil
}

func (r *ImportRepository) CreateIfAbsent(ctx context.Context, imp domain.Import) (bool, error) {
	var existing importRecord
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND document_id = ?", imp.TenantID, imp.DocumentID).
		First(&existing).Error
	if err == nil {
		return true, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	record := importFromDomain(imp)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		if isDuplicateKey(err) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (r *ImportRepository) SetStatus(ctx context.Context, importID string, st domain.ImportStatus) error {
	res := r.db.WithContext(ctx).
		Model(&importRecord{}).
		Where("id = ?", importID).
		Update("status", string(st))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrImportNotFound
	}
	return nil
}

func (r *ImportRepository) ListByTenant(ctx context.Context, tenantID string, limit int) ([]domain.Import, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}

	var total int64
	if err := r.db.WithContext(ctx).
		Model(&importRecord{}).
		Where("tenant_id = ?", tenantID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []importRecord
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	out := make([]domain.Import, 0, len(records))
	for _, rec := range records {
		out = append(out, importToDomain(rec))
	}
	return out, total, nil
}

func (r *ImportRepository) GetByDocumentID(ctx context.Context, tenantID, documentID string) (domain.Import, error) {
	var record importRecord
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND document_id = ?", tenantID, documentID).
		First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Import{}, domain.ErrImportNotFound
		}
		return domain.Import{}, err
	}
	return importToDomain(record), nil
}

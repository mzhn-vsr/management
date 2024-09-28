package faqservice

import (
	"context"
	"errors"
	"mzhn/management/internal/dto"
	"mzhn/management/internal/entity"
	"mzhn/management/internal/storage"
)

type FaqStore interface {
	Create(ctx context.Context, entry *dto.FaqEntryCreate) (string, error)
	Find(ctx context.Context, id string) (*entity.FaqEntry, error)
	List(ctx context.Context, filters dto.FaqEntryList) ([]*entity.FaqEntry, uint64, error)
	Update(ctx context.Context, entry *dto.FaqEntryUpdate) error
	Delete(ctx context.Context, id string) error
}

type FaissRepo interface {
	Save(ctx context.Context, e []*dto.FaqFaissCreate) error
	Delete(ctx context.Context, ids []string) error
}

type FaqService struct {
	FaqStore
	faiss FaissRepo
}

func New(faqstore FaqStore, faiss FaissRepo) *FaqService {
	return &FaqService{
		FaqStore: faqstore,
		faiss:    faiss,
	}
}

func (svc *FaqService) Create(ctx context.Context, entry *dto.FaqEntryCreate) (string, error) {
	id, err := svc.FaqStore.Create(ctx, entry)
	if err != nil {
		return "", err
	}

	if err := svc.faiss.Save(ctx, []*dto.FaqFaissCreate{
		{
			Id:             id,
			FaqEntryCreate: *entry,
		},
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (svc *FaqService) Delete(ctx context.Context, id string) error {
	if err := svc.FaqStore.Delete(ctx, id); err != nil {
		return err
	}

	if err := svc.faiss.Delete(ctx, []string{id}); err != nil {
		return err
	}

	return nil
}

func (s *FaqService) Find(ctx context.Context, id string) (*entity.FaqEntry, error) {
	e, err := s.FaqStore.Find(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNoFaqEntry) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e, nil
}

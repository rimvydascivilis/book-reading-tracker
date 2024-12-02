package list

import (
	"context"
	"fmt"
	"time"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/dto"
)

type listService struct {
	listRepo      domain.ListRepository
	listItemRepo  domain.ListItemRepository
	bookRepo      domain.BookRepository
	validationSvc domain.ValidationService
}

func NewListService(repo domain.ListRepository, listItemRepo domain.ListItemRepository,
	bookRepo domain.BookRepository, validationSvc domain.ValidationService) domain.ListService {
	return &listService{
		listRepo:      repo,
		listItemRepo:  listItemRepo,
		bookRepo:      bookRepo,
		validationSvc: validationSvc,
	}
}

func (s *listService) ListLists(ctx context.Context, userID int64) ([]dto.ListListsResponse, error) {
	lists, err := s.listRepo.GetListsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var listsResp []dto.ListListsResponse = make([]dto.ListListsResponse, 0, len(lists))
	for _, list := range lists {
		listsResp = append(listsResp, dto.ListListsResponse{
			ID:    list.ID,
			Title: list.Title,
		})
	}

	return listsResp, nil
}

func (s *listService) GetList(ctx context.Context, userID int64, listID int64) (dto.ListResponse, error) {
	list, err := s.listRepo.GetListByID(ctx, listID)
	if err != nil {
		return dto.ListResponse{}, err
	}

	listItems, err := s.listItemRepo.GetListItemsByListID(ctx, listID)
	if err != nil {
		return dto.ListResponse{}, err
	}

	var listItemsResp []dto.ListItemsResponse = make([]dto.ListItemsResponse, 0, len(listItems))
	for _, item := range listItems {
		book, err := s.bookRepo.GetBookByUserID(ctx, userID, item.BookID)
		if err != nil {
			return dto.ListResponse{}, err
		}

		listItemsResp = append(listItemsResp, dto.ListItemsResponse{
			ID:       item.ID,
			ListID:   item.ListID,
			BookName: book.Title,
		})
	}

	return dto.ListResponse{
		ID:        list.ID,
		Title:     list.Title,
		ListItems: listItemsResp,
	}, nil
}

func (s *listService) CreateList(ctx context.Context, userID int64, req dto.ListRequest) (dto.ListResponse, error) {
	list := domain.List{
		UserID:    userID,
		Title:     req.Title,
		CreatedAt: time.Now(),
	}

	if err := s.validationSvc.ValidateStruct(list); err != nil {
		return dto.ListResponse{}, err
	}

	list, err := s.listRepo.CreateList(ctx, list)
	if err != nil {
		return dto.ListResponse{}, err
	}

	return dto.ListResponse{
		ID:        list.ID,
		Title:     list.Title,
		ListItems: []dto.ListItemsResponse{},
	}, nil
}

func (s *listService) AddBookToList(ctx context.Context, userID, listID, bookID int64) error {
	list, err := s.listRepo.GetListByID(ctx, listID)
	if err != nil {
		return err
	}

	if list.UserID != userID {
		return fmt.Errorf("%w: %s", domain.ErrForbidden, "list does not belong to user")
	}

	book, err := s.bookRepo.GetBookByUserID(ctx, userID, bookID)
	if err != nil {
		return err
	}

	if book.UserID != userID {
		return fmt.Errorf("%w: %s", domain.ErrForbidden, "book does not belong to user")
	}

	listItem := domain.ListItem{
		ListID: listID,
		BookID: bookID,
	}

	_, err = s.listItemRepo.CreateListItem(ctx, listItem)
	if err != nil {
		return err
	}

	return nil
}

func (s *listService) RemoveBookFromList(ctx context.Context, userID, listID, itemID int64) error {
	list, err := s.listRepo.GetListByID(ctx, listID)
	if err != nil {
		return err
	}

	if list.UserID != userID {
		return fmt.Errorf("%w: %s", domain.ErrForbidden, "list does not belong to user")
	}

	err = s.listItemRepo.DeleteListItem(ctx, itemID)
	if err != nil {
		return err
	}

	return nil
}

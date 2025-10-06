package service

import (
	"context"

	"booknest/internal/domain"
)

type bookServiceImpl struct {
	r domain.BookRepository
}

func NewBookServiceImpl(r domain.BookRepository) domain.BookService {
	return &bookServiceImpl{r: r}
}

func (s *bookServiceImpl) CreateBook(in domain.BookInput) (book domain.Book, err error) {
	book.Author = in.Author
	book.Price = in.Price
	book.Stock = in.Stock
	book.Title = in.Title

	err = s.r.CreateBook(context.TODO(), &book)
	if err != nil {
		return book, err
	}

	return book, err
}

package service

import (
	"context"

	"github.com/joaovds/learn-grpc-go/internal/db"
	"github.com/joaovds/learn-grpc-go/internal/pb"
)

type BookService struct {
	pb.UnimplementedBookServiceServer
	BookDB db.Book
}

func NewBookService(bookDB db.Book) *BookService {
  return &BookService{
    BookDB: bookDB,
  }
}

func (b *BookService) CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.Book, error) {
  book, err := b.BookDB.Create(in.Name, in.Isbn, in.AuthorId)
  if err != nil {
    return nil, err
  }

  return &pb.Book{
    Id:          book.ID,
    Name:        book.Name,
    Isbn:        book.ISBN,
    AuthorId:    book.AuthorID,
  }, nil
}

func (b *BookService) GetBooks(ctx context.Context, in *pb.BlankBook) (*pb.BookList, error) {
	books, err := b.BookDB.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Book, len(books))
	for i, book := range books {
		result[i] = &pb.Book{
      Id:          book.ID,
      Name:        book.Name,
      Isbn:        book.ISBN,
      AuthorId:    book.AuthorID,
		}
	}

	return &pb.BookList{
    Books: result,
	}, nil
}

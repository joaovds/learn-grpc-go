package service

import (
	"context"
	"io"

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

func (b *BookService) GetBook(ctx context.Context, in *pb.GetByIdBook) (*pb.Book, error) {
  book, err := b.BookDB.GetById(in.Id)
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

func (b *BookService) CreateBookStream(stream pb.BookService_CreateBookStreamServer) error {
  books := make([]*pb.Book, 4)

  for {
    book, err := stream.Recv()
    if err == io.EOF {
      return stream.SendAndClose(&pb.BookList{Books: books})
    }
    if err != nil {
      return err
    }

    result, err := b.BookDB.Create(book.Name, book.Isbn, book.AuthorId)
    if err != nil {
      return err
    }

    books = append(books, &pb.Book{
      Id:          result.ID,
      Name:        result.Name,
      Isbn:        result.ISBN,
      AuthorId:    result.AuthorID,
    })
  }
}

func (b *BookService) CreateBookStreamBidirectional(stream pb.BookService_CreateBookStreamBidirectionalServer) error {
  for {
    book, err := stream.Recv()
    if err == io.EOF {
      return nil
    }
    if err != nil {
      return err
    }

    result, err := b.BookDB.Create(book.Name, book.Isbn, book.AuthorId)
    if err != nil {
      return err
    }

    err = stream.Send(&pb.Book{
      Id:          result.ID,
      Name:        result.Name,
      Isbn:        result.ISBN,
      AuthorId:    result.AuthorID,
    })
    if err != nil {
      return err
    }
  }
}

package service

import (
	"context"
	"io"

	"github.com/joaovds/learn-grpc-go/internal/db"
	"github.com/joaovds/learn-grpc-go/internal/pb"
)

type AuthorService struct {
	pb.UnimplementedAuthorServiceServer
	AuthorDB db.Author
}

func NewAuthorService(authorDB db.Author) *AuthorService {
	return &AuthorService{
		AuthorDB: authorDB,
	}
}

func (a *AuthorService) CreateAuthor(ctx context.Context, in *pb.CreateAuthorRequest) (*pb.Author, error) {
	author, err := a.AuthorDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	return &pb.Author{
		Id:          author.ID,
		Name:        author.Name,
		Description: author.Description,
	}, nil
}

func (a *AuthorService) GetAuthors(ctx context.Context, in *pb.Blank) (*pb.AuthorList, error) {
	authors, err := a.AuthorDB.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Author, len(authors))
	for i, author := range authors {
		result[i] = &pb.Author{
			Id:          author.ID,
			Name:        author.Name,
			Description: author.Description,
		}
	}

	return &pb.AuthorList{
		Authors: result,
	}, nil
}

func (a *AuthorService) GetAuthor(ctx context.Context, in *pb.GetById) (*pb.Author, error) {
  author, err := a.AuthorDB.GetById(in.Id)
  if err != nil {
    return nil, err
  }

  return &pb.Author{
    Id:          author.ID,
    Name:        author.Name,
    Description: author.Description,
  }, nil
}

func (a *AuthorService) CreateAuthorStream(stream pb.AuthorService_CreateAuthorStreamServer) error {
  authors := make([]*pb.Author, 4)

  for {
    author, err := stream.Recv()
    if err == io.EOF {
      return stream.SendAndClose(&pb.AuthorList{Authors: authors})
    }
    if err != nil {
      return err
    }

    result, err := a.AuthorDB.Create(author.Name, author.Description)
    if err != nil {
      return err
    }

    authors = append(authors, &pb.Author{
      Id:          result.ID,
      Name:        result.Name,
      Description: result.Description,
    })
  }
}


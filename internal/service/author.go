package service

import (
	"context"

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

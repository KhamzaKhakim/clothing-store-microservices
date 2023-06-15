package server

import (
	"bytes"
	"clothing-store-clothes/internal/data"
	"clothing-store-clothes/internal/validator"
	"clothing-store-clothes/pkg/jsonlog"
	pb "clothing-store-clothes/pkg/pb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type Server struct {
	pb.UnimplementedClothesServiceServer
	Models data.Models
	Wg     sync.WaitGroup
	Logger *jsonlog.Logger
}

func (s *Server) CreateClothe(ctx context.Context, r *pb.Clothe) (*pb.Clothe, error) {

	clothe := &data.Clothe{
		Name:     r.Name,
		Price:    r.Price,
		Brand:    r.Brand,
		Color:    r.Color,
		Sizes:    r.Sizes,
		Sex:      r.Sex,
		Type:     r.Type,
		ImageURL: r.ImageUrl,
	}
	v := validator.New()

	if data.ValidateClothe(v, clothe); !v.Valid() {
		return nil, status.Errorf(codes.InvalidArgument, createKeyValuePairs(v.Errors))
	}

	err := s.Models.Clothes.Insert(clothe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while persisting the clothe")
	}

	r.Id = clothe.ID

	return r, nil
}

func (s *Server) ShowClothe(ctx context.Context, r *pb.ShowClotheRequest) (*pb.Clothe, error) {
	clothe, err := s.Models.Clothes.Get(r.GetId())

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, status.Errorf(codes.NotFound, "Clothe not found")
		default:
			return nil, status.Errorf(codes.Internal, "Error while getting clothe by id")
		}
	}

	return &pb.Clothe{
		Id:       clothe.ID,
		Name:     clothe.Name,
		Price:    clothe.Price,
		Brand:    clothe.Brand,
		Color:    clothe.Color,
		Sizes:    clothe.Sizes,
		Sex:      clothe.Sex,
		Type:     clothe.Type,
		ImageUrl: clothe.ImageURL,
	}, nil
}

func (s *Server) ListClothe(ctx context.Context, r *pb.ListClotheRequest) (*pb.ClotheList, error) {
	filter := data.Filters{
		Page:         r.Filter.Page,
		PageSize:     r.Filter.PageSize,
		SortSafelist: r.Filter.SortSafeList,
		Sort:         r.Filter.Sort,
	}
	clothes, err := s.Models.Clothes.GetAll(r.Name, r.Brand, r.PriceMax, r.PriceMin,
		r.Sizes, r.Color, r.Type, r.Sex, filter)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while getting clothes")
	}

	var clotheList []*pb.Clothe

	for i := 0; i < len(clothes); i++ {
		clotheList = append(clotheList, &pb.Clothe{
			Id:       clothes[i].ID,
			Name:     clothes[i].Name,
			Price:    clothes[i].Price,
			Brand:    clothes[i].Brand,
			Color:    clothes[i].Color,
			Sizes:    clothes[i].Sizes,
			Sex:      clothes[i].Sex,
			Type:     clothes[i].Sex,
			ImageUrl: clothes[i].ImageURL,
		})
	}

	return &pb.ClotheList{ClotheList: clotheList}, nil
}

func (s *Server) UpdateClothe(ctx context.Context, r *pb.UpdateClotheRequest) (*pb.Clothe, error) {
	clothe, err := s.Models.Clothes.Get(r.GetId())
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, status.Errorf(codes.NotFound, "Clothe not found")
		default:
			return nil, status.Errorf(codes.Internal, "Error while getting clothe by id")
		}
	}

	if r.Clothe.Name != "" {
		clothe.Name = r.Clothe.Name
	}
	if r.Clothe.Price != 0 {
		clothe.Price = r.Clothe.Price
	}
	if r.Clothe.Brand != "" {
		clothe.Brand = r.Clothe.Brand
	}
	if r.Clothe.Color != "" {
		clothe.Color = r.Clothe.Color
	}
	if r.Clothe.Sex != "" {
		clothe.Sex = r.Clothe.Sex
	}
	if r.Clothe.Type != "" {
		clothe.Type = r.Clothe.Type
	}
	if r.Clothe.ImageUrl != "" {
		clothe.ImageURL = r.Clothe.ImageUrl
	}

	err = s.Models.Clothes.Update(clothe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while updating clothe")
	}

	return &pb.Clothe{
		Id:       clothe.ID,
		Name:     clothe.Name,
		Price:    clothe.Price,
		Brand:    clothe.Brand,
		Color:    clothe.Color,
		Sizes:    clothe.Sizes,
		Sex:      clothe.Sex,
		Type:     clothe.Type,
		ImageUrl: clothe.ImageURL,
	}, nil

}

func (s *Server) DeleteClothe(ctx context.Context, r *pb.DeleteClotheRequest) (*pb.Clothe, error) {
	err := s.Models.Clothes.Delete(r.Id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, status.Errorf(codes.NotFound, "Clothe not found")
		default:
			return nil, status.Errorf(codes.Internal, "Error while deleting the clothe")
		}
	}
	return &pb.Clothe{}, nil
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

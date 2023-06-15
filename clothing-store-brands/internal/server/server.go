package server

import (
	"bytes"
	"clothing-store-brands/internal/data"
	"clothing-store-brands/internal/validator"
	"clothing-store-brands/pkg/jsonlog"
	pb "clothing-store-brands/pkg/pb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type Server struct {
	pb.UnimplementedBrandsServiceServer
	Models data.Models
	Wg     sync.WaitGroup
	Logger *jsonlog.Logger
}

func (s *Server) CreateBrand(ctx context.Context, r *pb.Brand) (*pb.Brand, error) {

	brand := &data.Brand{
		Name:        r.Name,
		Country:     r.Country,
		Description: r.Description,
		ImageURL:    r.ImageUrl,
	}
	v := validator.New()

	if data.ValidateBrand(v, brand); !v.Valid() {
		return nil, status.Errorf(codes.InvalidArgument, createKeyValuePairs(v.Errors))
	}

	err := s.Models.Brands.Insert(brand)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while persisting the brand")
	}

	r.Id = brand.ID

	return r, nil
}

func (s *Server) ShowBrand(ctx context.Context, r *pb.ShowBrandRequest) (*pb.Brand, error) {
	clothe, err := s.Models.Brands.Get(r.GetId())

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, status.Errorf(codes.NotFound, "Brand not found")
		default:
			return nil, status.Errorf(codes.Internal, "Error while getting clothe by id")
		}
	}

	return &pb.Brand{
		Id:          clothe.ID,
		Name:        clothe.Name,
		Country:     clothe.Country,
		Description: clothe.Description,
		ImageUrl:    clothe.ImageURL,
	}, nil
}

func (s *Server) ListBrand(ctx context.Context, r *pb.ListBrandRequest) (*pb.BrandList, error) {
	brands, err := s.Models.Brands.GetAll()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while getting brands")
	}

	var brandList []*pb.Brand

	for i := 0; i < len(brands); i++ {
		brandList = append(brandList, &pb.Brand{
			Id:          brands[i].ID,
			Name:        brands[i].Name,
			Country:     brands[i].Country,
			Description: brands[i].Description,
			ImageUrl:    brands[i].ImageURL,
		})
	}

	return &pb.BrandList{BrandList: brandList}, nil
}

func (s *Server) UpdateBrand(ctx context.Context, r *pb.UpdateBrandRequest) (*pb.Brand, error) {
	brand, err := s.Models.Brands.Get(r.GetId())
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, status.Errorf(codes.NotFound, "Brand not found")
		default:
			return nil, status.Errorf(codes.Internal, "Error while getting brand by id")
		}
	}

	if r.Brand.Name != "" {
		brand.Name = r.Brand.Name
	}
	if r.Brand.Country != "" {
		brand.Country = r.Brand.Country
	}
	if r.Brand.Description != "" {
		brand.Description = r.Brand.Description
	}

	if r.Brand.ImageUrl != "" {
		brand.ImageURL = r.Brand.ImageUrl
	}

	err = s.Models.Brands.Update(brand)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while updating brand")
	}

	return &pb.Brand{
		Id:          brand.ID,
		Name:        brand.Name,
		Country:     brand.Country,
		Description: brand.Description,
		ImageUrl:    brand.ImageURL,
	}, nil

}

func (s *Server) DeleteBrand(ctx context.Context, r *pb.DeleteBrandRequest) (*pb.Brand, error) {
	err := s.Models.Brands.Delete(r.Id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, status.Errorf(codes.NotFound, "Clothe not found")
		default:
			return nil, status.Errorf(codes.Internal, "Error while deleting the clothe")
		}
	}
	return &pb.Brand{}, nil
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

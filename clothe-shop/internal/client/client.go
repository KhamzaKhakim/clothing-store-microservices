package client

import (
	auth "clothing-store/pkg/pb/auth"
	brand "clothing-store/pkg/pb/brand"
	clothe "clothing-store/pkg/pb/clothe"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetAuthClient() auth.AuthServiceClient {
	dial, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//dial, err := grpc.Dial("clothing-store-auth:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
	}
	return auth.NewAuthServiceClient(dial)
}

func GetClotheClient() clothe.ClothesServiceClient {
	dial, err := grpc.Dial("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//dial, err := grpc.Dial("clothing-store-clothes:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
	}
	return clothe.NewClothesServiceClient(dial)
}

func GetBrandClient() brand.BrandsServiceClient {
	dial, err := grpc.Dial("localhost:8084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//dial, err := grpc.Dial("clothing-store-brands:8084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
	}
	return brand.NewBrandsServiceClient(dial)
}

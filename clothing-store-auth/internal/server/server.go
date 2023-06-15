package server

import (
	"bytes"
	"clothing-store-auth/internal/data"
	"clothing-store-auth/internal/validator"
	"clothing-store-auth/pkg/jsonlog"
	"clothing-store-auth/pkg/mailer"
	pb "clothing-store-auth/pkg/pb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"sync"
	"time"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	Models data.Models
	Mailer mailer.Mailer
	Wg     sync.WaitGroup
	Logger *jsonlog.Logger
}

func (s *Server) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := &data.User{
		Name:      r.GetName(),
		Email:     r.GetEmail(),
		Activated: false,
		Money:     100000,
	}
	err := user.Password.Set(r.GetPassword())
	if err != nil {
		return nil, err
	}
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		return nil, status.Errorf(codes.InvalidArgument, createKeyValuePairs(v.Errors))
	}

	err = s.Models.Users.Insert(user)

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "User with this email already exists")
	}

	err = s.Models.Roles.AddRolesForUser(user.ID, "USER")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create a role for user")
	}

	//TODO: have to make request to CartService
	//err = s.Models.Carts.CreateCartForUser(user.ID)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}

	token, err := s.Models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create a token")
	}
	s.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
		}
		err = s.Mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			s.Logger.PrintError(err, nil)
		}
	})
	return &pb.RegisterResponse{Name: user.Name, Email: user.Email, Activated: user.Activated, Money: user.Money}, nil
}

func (s *Server) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	v := validator.New()
	data.ValidateEmail(v, r.GetEmail())
	data.ValidatePasswordPlaintext(v, r.GetPassword())
	if !v.Valid() {
		return nil, status.Errorf(codes.InvalidArgument, createKeyValuePairs(v.Errors))
	}

	user, err := s.Models.Users.GetByEmail(r.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User with that email was not found")
	}

	match, err := user.Password.Matches(r.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while checking password")
	}
	if !match {
		return nil, status.Errorf(codes.Canceled, "Password is incorrect")
	}

	token, err := s.Models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while creating token")
	}

	return &pb.LoginResponse{Token: token.Plaintext}, nil
}

func (s *Server) Activate(ctx context.Context, r *pb.ActivateRequest) (*pb.ActivateResponse, error) {
	var input struct {
		TokenPlaintext string
	}

	input.TokenPlaintext = r.Token

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		return nil, status.Errorf(codes.InvalidArgument, createKeyValuePairs(v.Errors))
	}
	user, err := s.Models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Invalid or expired token")
	}

	user.Activated = true
	err = s.Models.Users.Update(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while activating the user")
	}
	err = s.Models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while deleting activation token")
	}

	return &pb.ActivateResponse{Name: user.Name, Email: user.Email, Activated: user.Activated, Money: user.Money}, nil
}

func (s *Server) Authenticate(ctx context.Context, r *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {

	headerParts := strings.Split(r.Token, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid authorization header")
	}

	token := headerParts[1]
	v := validator.New()
	if data.ValidateTokenPlaintext(v, token); !v.Valid() {
		return nil, status.Errorf(codes.InvalidArgument, createKeyValuePairs(v.Errors))
	}

	user, err := s.Models.Users.GetForToken(data.ScopeAuthentication, token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, status.Errorf(codes.PermissionDenied, "Invalid token")
		default:
			return nil, status.Errorf(codes.Internal, "Error while getting token from the database")
		}
	}
	return &pb.AuthenticateResponse{Id: user.ID, Activated: user.Activated, Money: user.Money}, nil
}

func (s *Server) Authorize(ctx context.Context, r *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	roles, err := s.Models.Roles.GetAllRolesForUser(r.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while getting roles from database")
	}
	if !roles.Include(r.GetRole()) {
		return nil, status.Errorf(codes.PermissionDenied, "User does not have an expected role")
	}
	return &pb.AuthorizeResponse{Msg: "User can access"}, nil
}

func (s *Server) DeleteUser(ctx context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id := r.Id
	err := s.Models.Users.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, status.Errorf(codes.NotFound, "User not found")
		default:
			return nil, status.Errorf(codes.Internal, "Error while deleting user")
		}
	}
	return &pb.DeleteUserResponse{Message: "User successfully deleted"}, nil
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func (s *Server) background(fn func()) {
	s.Wg.Add(1)
	go func() {
		defer s.Wg.Done()
		defer func() {
			if err := recover(); err != nil {
				s.Logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		fn()
	}()
}

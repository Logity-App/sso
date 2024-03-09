package auth

import (
	"context"
	ssov1 "github.com/Logity-App/contracts/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPCServer *grpc.Server) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{})
}

func (s *serverAPI) VerifyNewPhoneNumber(ctx context.Context, req *ssov1.VerifyNewPhoneNumberRequest) (*ssov1.VerifyNewPhoneNumberResponse, error) {
	return &ssov1.VerifyNewPhoneNumberResponse{
		Status:         "Status",
		SmsCode:        "SmsCode",
		ExpirationTime: 1000,
	}, nil
}
func (s *serverAPI) SendSmsCode(ctx context.Context, req *ssov1.Empty) (*ssov1.Empty, error) {
	return &ssov1.Empty{}, nil
}
func (s *serverAPI) SignUpByPhone(ctx context.Context, req *ssov1.SignUpByPhoneRequest) (*ssov1.SignUpByPhoneResponse, error) {
	return &ssov1.SignUpByPhoneResponse{
		Status: "Status",
	}, nil
}
func (s *serverAPI) VerifyPhoneNumber(ctx context.Context, req *ssov1.VerifyPhoneNumberRequest) (*ssov1.VerifyPhoneNumberResponse, error) {
	return &ssov1.VerifyPhoneNumberResponse{
		Status:         "Status",
		SmsCode:        "SmsCode",
		ExpirationTime: 1000,
	}, nil
}
func (s *serverAPI) SignInByPhone(ctx context.Context, req *ssov1.SignInByPhoneRequest) (*ssov1.SignInByPhoneResponse, error) {
	return &ssov1.SignInByPhoneResponse{
		AccessToken:  "AccessToken",
		RefreshToken: "RefreshToken",
	}, nil
}

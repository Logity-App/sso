package auth

import (
	"context"
	ssov1 "github.com/Logity-App/contracts/gen/go/sso"
	"github.com/Logity-App/sso/internal/domain/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	StatusSuccess = "Success"
	StatusFail    = "Fail"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	VerifyNewPhoneNumber(ctx context.Context, phone string) (*models.Code, error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{
		auth: auth,
	})
}

func (s *serverAPI) VerifyNewPhoneNumber(ctx context.Context, req *ssov1.VerifyNewPhoneNumberRequest) (*ssov1.VerifyNewPhoneNumberResponse, error) {
	code, err := s.auth.VerifyNewPhoneNumber(ctx, req.Phone)

	if err != nil {
		return &ssov1.VerifyNewPhoneNumberResponse{
			Status: StatusFail,
		}, status.Error(codes.InvalidArgument, "phone is required") // TODO
	}

	return &ssov1.VerifyNewPhoneNumberResponse{
		Status:         StatusSuccess,
		SmsCode:        code.Value,
		ExpirationTime: code.ExpirationAt,
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

package auth

import (
	"context"
	ssov1 "github.com/Logity-App/contracts/gen/go/sso"
	"github.com/Logity-App/sso/internal/domain/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: вынести enums
const (
	StatusSuccess = "Success"
	StatusFail    = "Fail"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	VerifyNewPhoneNumber(ctx context.Context, phone string) (*models.TempUser, error)
	SendSmsCode(ctx context.Context, phone string, code string) error
	SignUpByPhone(ctx context.Context, dto models.SignUpByPhoneDTO) (*models.User, error)
	VerifyPhoneNumber(ctx context.Context, phone string) (*models.User, error)
	SignInByPhone(ctx context.Context, phone string, code string) (*models.User, error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{
		auth: auth,
	})
}

func (s *serverAPI) VerifyNewPhoneNumber(ctx context.Context, req *ssov1.VerifyNewPhoneNumberRequest) (*ssov1.VerifyNewPhoneNumberResponse, error) {
	tempUser, err := s.auth.VerifyNewPhoneNumber(ctx, req.Phone)

	if err != nil {
		return &ssov1.VerifyNewPhoneNumberResponse{
			Status: StatusFail,
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.VerifyNewPhoneNumberResponse{
		Status:  StatusSuccess,
		SmsCode: tempUser.Code,
	}, nil
}

// SendSmsCode TODO : хз зачем этот метод, можно сделать кнопку отправить код ещё раз
func (s *serverAPI) SendSmsCode(ctx context.Context, req *ssov1.SendSmsCodeRequest) (*ssov1.SendSmsCodeResponse, error) {
	err := s.auth.SendSmsCode(ctx, req.Phone, req.SmsCode)

	if err != nil {
		return &ssov1.SendSmsCodeResponse{
			Status: StatusFail,
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.SendSmsCodeResponse{
		Status: StatusSuccess,
	}, nil
}
func (s *serverAPI) SignUpByPhone(ctx context.Context, req *ssov1.SignUpByPhoneRequest) (*ssov1.SignUpByPhoneResponse, error) {
	dto := models.SignUpByPhoneDTO{
		Phone:        req.GetPhone(),
		BirthdayDate: req.GetBirthdayDate(),
		DefaultTag:   req.GetDefaultTag(),
	}

	user, err := s.auth.SignUpByPhone(ctx, dto)

	if err != nil {
		return &ssov1.SignUpByPhoneResponse{
			Status: StatusFail,
		}, status.Error(codes.Unknown, err.Error()) // TODO: add validation
	}

	return &ssov1.SignUpByPhoneResponse{
		Status:      StatusSuccess,
		AccessToken: user.Token.AccessToken,
	}, nil
}
func (s *serverAPI) VerifyPhoneNumber(ctx context.Context, req *ssov1.VerifyPhoneNumberRequest) (*ssov1.VerifyPhoneNumberResponse, error) {
	if req.Phone == "" {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	user, err := s.auth.VerifyPhoneNumber(ctx, req.Phone)

	if err != nil {
		return &ssov1.VerifyPhoneNumberResponse{
			Status: StatusFail,
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.VerifyPhoneNumberResponse{
		Status:         StatusSuccess,
		SmsCode:        user.SmsCode.Value,
		ExpirationTime: user.SmsCode.ExpirationAt.Unix(),
	}, nil
}
func (s *serverAPI) SignInByPhone(ctx context.Context, req *ssov1.SignInByPhoneRequest) (*ssov1.SignInByPhoneResponse, error) {
	if req.Phone == "" {
		return &ssov1.SignInByPhoneResponse{
			Status: StatusFail,
		}, status.Error(codes.InvalidArgument, "phone is required")
	}

	if req.SmsCode == "" {
		return &ssov1.SignInByPhoneResponse{
			Status: StatusFail,
		}, status.Error(codes.InvalidArgument, "sms code is required")
	}

	user, err := s.auth.SignInByPhone(ctx, req.GetPhone(), req.GetSmsCode())

	if err != nil {
		return &ssov1.SignInByPhoneResponse{
			Status: StatusFail,
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.SignInByPhoneResponse{
		Status:       StatusSuccess,
		AccessToken:  user.Token.AccessToken,
		RefreshToken: user.Token.RefreshToken,
	}, nil
}

package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Fuchsoria/banners-rotation/internal/app"
	gw "github.com/Fuchsoria/banners-rotation/internal/server/pb/api"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	app    app.App
	server *http.Server
}

type grpcserver struct {
	gw.UnimplementedBannersRotationServer
	app app.App
}

var ErrBadRequest = errors.New("bad request")

func NewServer(app *app.App, address string, port string, grpcPort string) (*Server, error) {
	grpcServerEndpoint := net.JoinHostPort(address, grpcPort)

	lis, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to listen, %w", err)
	}

	logger := app.GetLogger().GetInstance()

	s := grpc.NewServer(grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_zap.StreamServerInterceptor(logger),
	)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
		)))

	gw.RegisterBannersRotationServer(s, &grpcserver{app: *app})

	go func() {
		err := s.Serve(lis)
		if err != nil {
			app.GetLogger().Error(fmt.Errorf("cannot serve grpc, %w", err).Error())
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		grpcServerEndpoint,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server, %w", err)
	}

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	gwmux := runtime.NewServeMux()
	err = gw.RegisterBannersRotationHandler(ctx, gwmux, conn)
	if err != nil {
		return nil, fmt.Errorf("cannot register app handler, %w", err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	server := &http.Server{
		Addr:         net.JoinHostPort(address, port),
		Handler:      gwmux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{*app, server}, nil
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("cannot start gateway server, %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("cannot shutdown gateway server, %w", err)
	}

	return nil
}

func (s *grpcserver) AddBanner(ctx context.Context, in *gw.AddBannerRequest) (*gw.MessageResponse, error) {
	if in.BannerId == "" || in.SlotId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "cannot add banner in rotation, %s", ErrBadRequest)
	}

	err := s.app.AddBannerRotation(in.BannerId, in.SlotId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot add banner in rotation, %s", err)
	}

	return &gw.MessageResponse{Message: "added"}, nil
}

func (s *grpcserver) RemoveBanner(ctx context.Context, in *gw.RemoveBannerRequest) (*gw.MessageResponse, error) {
	if in.BannerId == "" || in.SlotId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "cannot remove banner from rotation, %s", ErrBadRequest)
	}

	err := s.app.RemoveBannerRotation(in.BannerId, in.SlotId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot remove banner from rotation, %s", err)
	}

	return &gw.MessageResponse{Message: "removed"}, nil
}

func (s *grpcserver) ClickEvent(ctx context.Context, in *gw.ClickEventRequest) (*gw.MessageResponse, error) {
	if in.BannerId == "" || in.SlotId == "" || in.SocialDemoId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "cannot click on banner, %s", ErrBadRequest)
	}

	err := s.app.AddClickEvent(in.BannerId, in.SlotId, in.SocialDemoId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot add click event, %s", err)
	}

	return &gw.MessageResponse{Message: "clicked"}, nil
}

func (s *grpcserver) GetBanner(ctx context.Context, in *gw.GetBannerRequest) (*gw.BannerResponse, error) {
	if in.SlotId == "" || in.SocialDemoId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "cannot get banner, %s", ErrBadRequest)
	}

	ID, err := s.app.GetBanner(in.SlotId, in.SocialDemoId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get banners, %s", err)
	}

	return &gw.BannerResponse{Id: ID}, nil
}

func (s *grpcserver) CreateBanner(ctx context.Context, in *gw.BannerRequest) (*gw.BannerResponse, error) {
	ID := in.Id

	if ID == "" {
		ID = uuid.NewString()
	}

	ID, err := s.app.CreateBanner(ID, in.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create banner, %s", err)
	}

	return &gw.BannerResponse{Id: ID}, nil
}

func (s *grpcserver) CreateSlot(ctx context.Context, in *gw.SlotRequest) (*gw.SlotResponse, error) {
	ID := in.Id

	if ID == "" {
		ID = uuid.NewString()
	}

	ID, err := s.app.CreateSlot(ID, in.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create slot, %s", err)
	}
	return &gw.SlotResponse{Id: ID}, nil
}

func (s *grpcserver) CreateSocialDemo(ctx context.Context, in *gw.SocialDemoRequest) (*gw.SocialDemoResponse, error) {
	ID := in.Id

	if ID == "" {
		ID = uuid.NewString()
	}

	ID, err := s.app.CreateSocialDemo(ID, in.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create social demo, %s", err)
	}

	return &gw.SocialDemoResponse{Id: ID}, nil
}

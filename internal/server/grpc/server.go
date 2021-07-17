package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Fuchsoria/banners-rotation/internal/app"
	gw "github.com/Fuchsoria/banners-rotation/internal/server/pb/BannersRotation"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
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
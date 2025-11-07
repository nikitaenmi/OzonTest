package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikitaenmi/OzonTest/internal/adapter"
	"github.com/nikitaenmi/OzonTest/internal/config"
	"github.com/nikitaenmi/OzonTest/internal/database"
	"github.com/nikitaenmi/OzonTest/internal/repositories"
	"github.com/nikitaenmi/OzonTest/internal/server/handler"
	"github.com/nikitaenmi/OzonTest/internal/services"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nikitaenmi/OzonTest/gen"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port, cfg.DB.SSLMode)

	db, err := database.ConnectGORM(dsn)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	repo := repositories.NewRepo(db)
	cbr := adapter.NewMockCBRClient()
	conventor := services.NewCurrencyConverter(cbr)
	paymentService := services.NewPaymentService(repo, conventor, cfg.Payment.MaxAmountRUB)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	grpcServer := startGRPCServer(paymentHandler, cfg.Server.GRPCPort)
	httpServer := startGateway(cfg.Server.HTTPPort, cfg.Server.GRPCPort)

	slog.Info("servers started",
		"grpc_port", cfg.Server.GRPCPort,
		"http_port", cfg.Server.HTTPPort,
		"swagger_ui", "http://localhost"+cfg.Server.HTTPPort+"/docs",
		"swagger_json", "http://localhost"+cfg.Server.HTTPPort+"/swagger.json")

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.GRPCTimeout)
	defer cancel()

	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("http gateway shutdown error", "error", err)
	}

	slog.Info("servers stopped")
}

func startGRPCServer(handler *handler.Payment, port string) *grpc.Server {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error("failed to listen", "port", port, "error", err)
		os.Exit(1)
	}

	server := grpc.NewServer()
	gen.RegisterPaymentServiceServer(server, handler)

	go func() {
		slog.Info("gRPC server started", "port", port)
		if err := server.Serve(lis); err != nil {
			slog.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	return server
}

func startGateway(httpPort, grpcPort string) *http.Server {
	mainMux := http.NewServeMux()

	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := gen.RegisterPaymentServiceHandlerFromEndpoint(
		context.Background(),
		gwMux,
		grpcPort,
		opts,
	)
	if err != nil {
		slog.Error("failed to register gateway", "error", err)
		os.Exit(1)
	}

	swaggerHandler := handler.NewSwaggerHandler()

	mainMux.Handle("/", gwMux)
	mainMux.HandleFunc("/docs", swaggerHandler.ServeSwaggerUI)
	mainMux.HandleFunc("/swagger.json", swaggerHandler.ServeSwaggerJSON)

	server := &http.Server{
		Addr:    httpPort,
		Handler: mainMux,
	}

	go func() {
		slog.Info("HTTP gateway started", "port", httpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP gateway failed", "error", err)
			os.Exit(1)
		}
	}()

	return server
}

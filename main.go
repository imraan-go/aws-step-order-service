package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	cfg "github.com/imraan-go/aws-step-order-service/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"
)

var dnc *dynamodb.Client

func main() {
	conf := cfg.NewConfig("config.env")

	// Setup database

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AwsKey, conf.AwsSecret, "")), config.WithDefaultRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	dnc = dynamodb.NewFromConfig(cfg)

	e := echo.New()

	// Recover from panics
	e.Use(middleware.Recover())

	// Allow requests from *
	e.Use(middleware.CORS())
	// Print http request and response log to stdout if debug is enabled
	if conf.Debug {
		e.Use(middleware.Logger())
	}

	setupRoutes(e)

	// Start HTTP Server

	go httpServer(e, conf.HTTP)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("Shutting down HTTP server...")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("HTTP server stopped!")
}

func httpServer(e *echo.Echo, httpConfig cfg.HTTP) {
	var err error
	log.Println("H2C Mode:", httpConfig.H2C)
	if httpConfig.H2C {
		s := &http2.Server{
			// setting MaxConcurrentStreams to a bigger number because this server will be behind load balancer
			MaxConcurrentStreams: 2500,
			MaxReadFrameSize:     1048576,
			IdleTimeout:          10 * time.Second,
		}

		if err = e.StartH2CServer(httpConfig.HTTPAddress, s); err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}

	} else {
		if err = e.Start(httpConfig.HTTPAddress); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}

}

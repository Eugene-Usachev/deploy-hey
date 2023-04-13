package main

import (
	"GoServer/internal/handler"
	"GoServer/internal/repository"
	"GoServer/internal/service"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func HandlerLamd(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return handler.GinLambda.ProxyWithContext(ctx, req)
}

func main() {

	pool, err := repository.NewPostgresDB(context.Background(), 12,
		repository.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			UserName: os.Getenv("DB_USERNAME"),
			UserPass: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("SSL_MODE"),
		})
	if err != nil {
		//log.Fatalf("error in connection to database: %s", err)
	}

	repositoryImpl := repository.NewRepository(pool)
	serviceImpl := service.NewService(repositoryImpl)
	handlerImpl := handler.NewHandler(serviceImpl)

	handlerImpl.InitRoutes()

	lambda.Start(HandlerLamd)

	//serverImpl := new(server.Server)
	//
	//if err = serverImpl.Run(os.Getenv("PORT"), handlerImpl.InitRoutes()); err != nil {
	//	log.Fatalf("error occured while running http server: %s", err.Error())
	//}
}

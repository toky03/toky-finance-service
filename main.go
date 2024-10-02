package main

import (
	"log"
	"os"

	"github.com/toky03/toky-finance-accounting-service/api"
	"github.com/toky03/toky-finance-accounting-service/handler"
	"github.com/toky03/toky-finance-accounting-service/repository"
	"github.com/toky03/toky-finance-accounting-service/service"
)

func main() {

	log.SetOutput(os.Stdout)

	bookRepository := repository.CreateRepository()

	bookService := service.CreateBookService(bookRepository)
	accountingService :=
		service.CreateAccountingService()
	userService := service.CreateApplicationUserService()

	bookHandler := handler.CreateBookRealmHandler(bookService, userService)
	monitoringHandler := handler.CreateMonitoringHandler()
	accountingHandler := handler.CreateAccountingHandler(accountingService, userService)
	authenticationHandler := handler.CreateAuthenticationHandler()

	server := api.CreateServer(
		bookHandler,
		monitoringHandler,
		accountingHandler,
		authenticationHandler,
	)

	err := handler.CreateAndRegisterUserBatchService()
	if err != nil {
		log.Fatal(err)
	}
	server.RegisterHandlers()
	server.ServeHTTP()

}

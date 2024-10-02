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

	accountingService := service.CreateBookService(bookRepository)
	userService := service.CreateApplicationUserService()

	bookHandler := handler.CreateBookRealmHandler(accountingService, userService)
	monitoringHandler := handler.CreateMonitoringHandler()
	accountingHandler := handler.CreateAccountingHandler()
	authenticationHandler := handler.CreateAuthenticationHandler()

	server := api.CreateServer(bookHandler, monitoringHandler, accountingHandler, authenticationHandler)

	err := handler.CreateAndRegisterUserBatchService()
	if err != nil {
		log.Fatal(err)
	}
	server.RegisterHandlers()
	server.ServeHTTP()

}

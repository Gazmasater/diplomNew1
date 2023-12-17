package main

import (
	"net/http"
	"time"

	"diplom.com/go-musthave-diploma-tpl/internal/authentication"
	"diplom.com/go-musthave-diploma-tpl/internal/config"
	"diplom.com/go-musthave-diploma-tpl/internal/handlers/addorder"
	"diplom.com/go-musthave-diploma-tpl/internal/handlers/authorization"
	"diplom.com/go-musthave-diploma-tpl/internal/handlers/charge"
	"diplom.com/go-musthave-diploma-tpl/internal/handlers/getbalance"
	"diplom.com/go-musthave-diploma-tpl/internal/handlers/getorders"
	"diplom.com/go-musthave-diploma-tpl/internal/handlers/getwithdraw"
	"diplom.com/go-musthave-diploma-tpl/internal/handlers/registration"
	"diplom.com/go-musthave-diploma-tpl/internal/interrogator"
	"diplom.com/go-musthave-diploma-tpl/internal/logger"
	"diplom.com/go-musthave-diploma-tpl/internal/storage/postgres"
	"diplom.com/go-musthave-diploma-tpl/internal/storage/redis"
	"diplom.com/go-musthave-diploma-tpl/internal/systemservice"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	log, err := logger.InitLogger()
	if err != nil {
		panic("couldn't init logger")
	}

	db, err := postgres.New(cfg.DefaultDBConnStr)
	if err != nil {
		log.LogWarning(err)
	}
	log.LogInfo("database connected")

	redisClient := redis.NewRedisClient(cfg.RedisAddress)
	pong, err := redisClient.Ping()
	if err != nil {
		log.LogWarning("redis connection error:", err)
	}
	log.LogInfo("Connection to redis established:", pong)

	JWTMiddleware := authentication.JWTMiddleware{
		RedisClient: redisClient,
		Log:         log,
	}

	app := systemservice.NewService(db)

	interrog := interrogator.NewInterrogator(db, log, cfg)
	go func() {
		for {
			interrog.OrderStatusWorker()
			time.Sleep(1 * time.Second)
		}
	}()
	r := chi.NewRouter()

	r.Use(log.MyLogger)

	r.Group(func(r chi.Router) {
		r.Use(JWTMiddleware.JWTMiddleware())
		r.Post("/api/user/orders", addorder.NewPutOrderNumberHandler(app, redisClient, log).AddOrderNumberHandler)
		r.Get("/api/user/orders", getorders.NewGetOrdersHandler(app, log, interrog).GetOrdersHandler)
		r.Get("/api/user/balance", getbalance.NewGetBalanceHandler(app, log).GetUserBalanceHandler)
		r.Post("/api/user/balance/withdraw", charge.NewChargeHandler(app, log).ChargeHandler)
		r.Get("/api/user/withdrawals", getwithdraw.NewGetWithdrawHandler(app, log).GetWithdrawHandler)
	})

	r.Post("/api/user/register", registration.NewRegistrationHandler(app, redisClient, log).RegistrationHandler)
	r.Post("/api/user/login", authorization.NewAuthorizationHandler(app, redisClient, log).AuthorizationHandler)

	log.LogInfo("starting server at localhost", cfg.Address)
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.LogWarning(err)
	}
}

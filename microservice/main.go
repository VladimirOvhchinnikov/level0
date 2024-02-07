package main

import (
	"context"
	"log"
	"microservice/internal/controller"
	"microservice/internal/infrastructures/cache"
	"microservice/internal/infrastructures/database"
	"microservice/internal/infrastructures/nats"
	"microservice/internal/infrastructures/route"
	"microservice/internal/infrastructures/server"
	"microservice/internal/usecase"
	"time"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

func main() {

	time.Sleep(5 * time.Second)
	//Инициализируем логгер
	logger, _ := zap.NewProduction()

	//Инициализируем Nats
	sc, err := stan.Connect("wbCluster", "client-subscriber", stan.NatsURL("nats://nats:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	//Инициализируем Redis
	redisCashe := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	redisClient := cache.NewRedisCashe(logger, redisCashe)

	//Инициализируем Postgres
	databaseURL := "postgres://postgres:0@postgres:5432/wb"
	dbPool, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	postgresCleint := database.NewDataBase(logger, dbPool)

	//инициализиурем чекер натса
	natsCheker := nats.NewNats(logger, sc, redisCashe)
	natsCheker.SubscribeToChannel("wbCluster")

	//Инициализируем слой бизнес логики
	usecaset := usecase.NewGetData(logger, redisClient, postgresCleint)

	//Инициализируем контроллер для ID
	controllerID := controller.NewIDSearch(logger, usecaset)
	controllerWeb := controller.NewWebRequest(logger, usecaset)

	//Инициализируем роутер
	router := route.NewRouter(logger, controllerID, controllerWeb)
	//Инициализируем минималный роутинг для задачи
	router.MyRouter()
	//Инициализируем сервер
	server := server.NewGoChi(logger, router.Router)
	//Запускаем сервер
	server.StartServer(":8080")
}

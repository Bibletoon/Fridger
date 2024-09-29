package main

import (
	"Fridger/internal/configuration"
	"Fridger/internal/handlers"
	"Fridger/internal/infrastructure/clients"
	"Fridger/internal/infrastructure/db"
	"Fridger/internal/infrastructure/repositories"
	"Fridger/internal/jobs"
	"Fridger/internal/services"
	"context"
	configuration_yaml_file "github.com/BoRuDar/configuration-yaml-file"
	configlib "github.com/BoRuDar/configuration/v4"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makiuchi-d/gozxing/datamatrix"
	_ "image/jpeg"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := configuration.Configuration{}
	configurator := configlib.New(
		&cfg,
		configlib.NewEnvProvider(),
		configuration_yaml_file.NewYAMLFileProvider("config.yaml"),
		configlib.NewDefaultProvider(),
	)
	err := configurator.InitValues()
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(ctx, cfg.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		panic(err)
	}

	if cfg.MigrationsConfiguration.RunOnStart {
		err = db.Migrate(pool, cfg.MigrationsConfiguration)
		if err != nil {
			panic(err)
		}
	}

	crptClient := clients.NewCrptClient()
	productRepo := repositories.NewProductRepo(pool)
	productService := services.NewProductService(productRepo, crptClient)

	dmReader := datamatrix.NewDataMatrixReader()
	photoService := services.NewPhotoService(dmReader)

	photoHandler := handlers.NewPhotoHandler(photoService, productService)
	productsListHandler := handlers.NewProductsListHandler(productService)
	expiringProductsHandler := handlers.NewExpiringProductsHandler(productService, cfg.JobsConfiguration.ExpiringProductsJobConfiguration)

	bot, err := services.NewBot(
		cfg.BotConfiguration,
		photoHandler,
		productsListHandler,
		expiringProductsHandler)
	if err != nil {
		panic(err)
	}

	s, err := gocron.NewScheduler()

	if err != nil {
		panic(err)
	}

	expiringProductsJob := jobs.NewExpiringProductsJob(cfg.JobsConfiguration.ExpiringProductsJobConfiguration, bot, productService)

	task := gocron.NewTask(expiringProductsJob.Run)
	schedule := gocron.CronJob(cfg.JobsConfiguration.ExpiringProductsJobConfiguration.Schedule, false)

	_, err = s.NewJob(schedule, task, gocron.WithName("ExpiringProducts"))

	if err != nil {
		panic(err)
	}

	s.Start()
	bot.Start(ctx)
}

package main

import (
    "log"
    "github.com/joho/godotenv"
    "github.com/hibiken/asynq"
    "github.com/rampenke/zosma-api-server/tasks"
    "github.com/kelseyhightower/envconfig"
)

type Config struct {
	RedisAddr   string `envconfig:"REDIS_ADDR" required:"true"`
	RedisPassword  string `envconfig:"REDIS_PASSWORD" required:"true"`
    SdApiHost string `envconfig:"SD_API_HOST" required:"true"`
}
var cfg Config

func main() {
    _ = godotenv.Overload()
    if err := envconfig.Process("", &cfg); err != nil {
	    log.Fatal(err.Error())
    }
    srv := asynq.NewServer(
        asynq.RedisClientOpt{Addr: cfg.RedisAddr, Password: cfg.RedisPassword,},
        asynq.Config{
            // Specify how many concurrent workers to use
            Concurrency: 10,
            // Optionally specify multiple queues with different priority.
            Queues: map[string]int{
                "critical": 6,
                "default":  3,
                "low":      1,
            },
            // See the godoc for other configuration options
        },
    )

    // mux maps a type to a handler
    mux := asynq.NewServeMux()
    mux.Handle(tasks.TypeTxt2img, tasks.NewTxt2imgProcessor(cfg.SdApiHost))
    // ...register other handlers...

    if err := srv.Run(mux); err != nil {
        log.Fatalf("could not run server: %v", err)
    }
}

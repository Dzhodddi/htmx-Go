package main

import (
	"expvar"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"project/internal/auth"
	"project/internal/db"
	"project/internal/env"
	ratelimiter "project/internal/ratelimiter"
	store2 "project/internal/store"
	cache "project/internal/store/cache"
	"runtime"
	"time"
)

//	@title			Social site API
//	@description	API for SocialSite
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath					/v1
//
// @securityDefinitions.apikey ApiKeyAuth
// @in							header
// @name						Authorization
// @description
const version = "0.0.10"

func main() {

	cfg := config{
		addr:   env.GetString("ADDR", ":3050"),
		apiURL: env.GetString("API_URL", "localhost:3050"),
		//frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://admin:adminpasswrod@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisConfig: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  false,
		},
		env: env.GetString("ENV", "development"),
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("TOKEN_SECRET", "secret"),
				expiry: time.Hour * 24 * 3,
				issuer: "Social API",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestPerTimeFrame: 20,
			TimeFrame:           time.Second * 5,
			Enabled:             true,
		},
		//mail: mailConfig{
		//	exp:       time.Hour * 24 * 3,
		//	fromEmail: env.GetString("FROM_EMAIL", "dima2006x@email.com"),
		//	//sendGrid: sendGridConfig{
		//	//	apiKey: env.GetString("SENDGRID_API_KEY", ""),
		//	//},
		//	mailTrap: mailTrapConfig{
		//		apiKey: env.GetString("MAIL_TRAP_API_KEY", ""),
		//	},
		//},
	}

	//Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database
	database, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer database.Close()
	logger.Info("database initialized")
	store := store2.NewStorage(database)

	// jwt auth
	jwtAuth := auth.NewJWTAuth(cfg.auth.token.secret, cfg.auth.token.issuer, cfg.auth.token.issuer)

	//redis
	var cacheRedis *redis.Client
	if cfg.redisConfig.enabled {
		cacheRedis = cache.NewRedisClient(cfg.redisConfig.addr, cfg.redisConfig.password, cfg.redisConfig.db)
		logger.Info("cache initialized")
	}
	cacheStorage := cache.NewRedisStorage(cacheRedis)

	// rate limiter
	fixedRateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	// email
	//mailerConfig := mailer.NewSendGridMailer(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	//mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	//if err != nil {
	//	logger.Fatal(err)
	//}

	app := &application{
		config:       cfg,
		store:        store,
		cacheStorage: cacheStorage,
		logger:       logger,
		//mailer: mailtrap,
		authenticator: jwtAuth,
		rateLimiter:   fixedRateLimiter,
	}
	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return database.Stats()
	}))
	expvar.Publish("Go routines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	mux := app.mount()
	logger.Fatal(app.run(mux))
}

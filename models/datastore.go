package models

import (
	"database/sql"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DataStore struct {
	User IUserStore

	Redis    *redis.Client
	Postgres *sqlx.DB
}

type DBConf struct {
	Host     string
	DBName   string
	User     string
	Password string
}

func BuildStore() (*DataStore, error) {
	db, err := InitSQLStore()
	if err != nil {
		panic(err)
	}

	red, err := InitRedisStore()
	if err != nil {
		panic(err)
	}

	return &DataStore{
		User: NewUserStore(db, red),

		Redis:    red,
		Postgres: db,
	}, nil
}

func InitRedisStore() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func InitSQLStore() (*sqlx.DB, error) {
	var conf DBConf

	if os.Getenv("DBNAME") == "" {
		f, err := os.Open("./dbconf.json")
		defer f.Close()
		if err != nil {
			if err == os.ErrNotExist {
				err = json.NewEncoder(f).Encode(DBConf{
					Host:     "localhost",
					DBName:   "db",
					User:     "postgres",
					Password: "",
				})
				if err != nil {
					log.Println("dbconf.json writing error:", err)
					return nil, err
				}
				log.Println("dbconf.json created, fill it now!")
				defer os.Exit(0)
				return nil, err
			}
			log.Println("dbconf.json open error:", err)
			return nil, err
		}

		err = json.NewDecoder(f).Decode(&conf)
		if err != nil {
			log.Println("dbconf.json decode error:", err)
			return nil, err
		}
	} else {
		conf = DBConf{
			Host:     os.Getenv("DBHOST"),
			DBName:   os.Getenv("DBNAME"),
			User:     os.Getenv("DBUSER"),
			Password: os.Getenv("DBPASS"),
		}
	}

	db, err := sqlx.Connect("postgres", "host="+conf.Host+" user="+conf.User+" password="+conf.Password+" dbname="+conf.DBName+" sslmode=disable")
	if err != nil {
		log.Println("postgresql connection error:", err)
		return nil, err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Println("postgresql migrate init error:", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("postgresql migrate error:", err)
		return nil, err
	}
	return db, nil
}

func MigrateDown(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Println("postgresql migrate init error:", err)
	}

	if err = m.Drop(); err != nil {
		log.Println("postgresql drop error:", err)
		return err
	}

	return nil
}

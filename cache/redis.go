package cache

import (
	"errors"
	"log"
	"strconv"
	"time"

	redisClient "gopkg.in/redis.v5"
)

const KeyNotFound  = "key not found"

// Redis struct to manage redis.
type Redis struct {
	Client *redisClient.Client
	db     Config
}

// NewRedis is responsible for building a redis struct instance
func NewRedis(config Config) (*Redis, error) {

	red := Redis{db: config}
	err := red.Connect()
	if err != nil {
		return nil, err
	}
	return &red, nil
}

// Connect connects on redis database
func (r *Redis) Connect() error {
	db, _ := strconv.Atoi(r.db.GetDatabase())
	r.Client = redisClient.NewClient(&redisClient.Options{
		Addr:     r.db.GetHost() + ":" + strconv.Itoa(r.db.GetPort()),
		Password: r.db.GetPassword(),
		DB:       db,
	})

	_, err := r.Client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// Set set key.
func (r *Redis) Set(key, value string, duration time.Duration) error {
	_, err := r.Client.Set(key, value, duration).Result()
	if err != nil {
		return err
	}
	return nil
}

// Del delete key.
func (r *Redis) Del(key string) error {
	_, err := r.Client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

// Get get key.
func (r *Redis) Get(key string) (string, error) {
	value, err := r.Client.Get(key).Result()
	if err == redisClient.Nil {
		return "", errors.New("key not found")
	} else if err != nil {
		return "", err
	}
	return value, nil
}

// Exist test if key exists.
func (r *Redis) Exist(key string) (bool, error) {
	_, err := r.Client.Get(key).Result()

	if err == redisClient.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Close is responsible for closing redis connection
func (r *Redis) Close() {
	err := r.Client.Close()

	if err != nil {
		log.Println(err)
	}
}
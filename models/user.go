package models

import (
	"database/sql"
	"errors"
	"github.com/go-redis/redis"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IUserStore interface {
	Create(email, password string) (uuid.UUID, error)
	Login(email, password string) (string, error)
	Auth(sesid string) (uuid.UUID, error)
	Logout(sesid string) error
	GetAll() ([]User, error)
}

type User struct {
	ID        uuid.UUID  `db:"id"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	LastLogin *time.Time `db:"last_login"`
}

type UserStore struct {
	db    *sqlx.DB
	redis *redis.Client
}

var ErrLoginIncorrect = errors.New("incorrect email or password")
var ErrAuthIncorrect = errors.New("incorrect or old session")
var ErrAlreadyCreated = errors.New("user already exists")

func (us *UserStore) Create(email, password string) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	bpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	//TODO: pass context to queries
	_, err = us.db.NamedExec("INSERT INTO users (id, email, password, created_at) VALUES (:id,:email,:password,:created_at)", &User{
		ID:        id,
		Email:     email,
		Password:  string(bpw),
		CreatedAt: time.Now(),
	})

	//23505 is postgres' error code that means - item exists
	if pgerr, ok := err.(*pq.Error); ok {
		if pgerr.Code == "23505" {
			return uuid.Nil, ErrAlreadyCreated
		}
	}

	return id, err
}

func (us *UserStore) Auth(sesid string) (uuid.UUID, error) {
	res, err := us.redis.Get("user:session:" + sesid).Result()
	if err != nil {
		if err == redis.Nil {
			return uuid.Nil, ErrAuthIncorrect
		}
		return uuid.Nil, err
	}

	rid := uuid.FromStringOrNil(res)
	if rid == uuid.Nil {
		return uuid.Nil, ErrAuthIncorrect
	}

	return rid, nil
}

func (us *UserStore) Logout(sesid string) error {
	_, err := us.redis.Del("user:session:" + sesid).Result()
	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) Login(email, password string) (string, error) {
	var u User
	err := us.db.Get(&u, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrLoginIncorrect
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", ErrLoginIncorrect
	}

	//TODO: replace to more secure
	sesid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	wset, err := us.redis.SetNX("user:session:"+sesid.String(), u.ID.String(), 3*time.Hour).Result()
	if err != nil || !wset {
		return "", errors.New("session create error")
	}

	_, err = us.db.Exec("UPDATE users SET last_login=$2 WHERE id=$1", u.ID, time.Now())
	if err != nil {
		return "", err
	}
	return sesid.String(), nil
}

func (us *UserStore) GetAll() ([]User, error) {
	var res []User
	err := us.db.Select(&res, "SELECT id,email,password,created_at,last_login FROM users")
	return res, err
}

func NewUserStore(db *sqlx.DB, red *redis.Client) *UserStore {
	return &UserStore{db: db, redis: red}
}

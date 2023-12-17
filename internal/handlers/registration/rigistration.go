package registration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"diplom.com/go-musthave-diploma-tpl/internal/authentication"
	"diplom.com/go-musthave-diploma-tpl/internal/credentials"
	"diplom.com/go-musthave-diploma-tpl/internal/logger"
	"github.com/google/uuid"
)

type RegisterHandler struct {
	setUser     UserSetter
	RedisClient *redis.RedisClient
	log         *logger.Logger
}

//go:generate go run github.com/vektra/mockery/v2@v2.36.1 --all
type UserSetter interface {
	UserSetter(ctx context.Context, user, password, id string) error
}

func NewRegistrationHandler(us UserSetter, redis *redis.RedisClient, log *logger.Logger) *RegisterHandler {
	return &RegisterHandler{
		setUser:     us,
		RedisClient: redis,
		log:         log,
	}
}

func (reg *RegisterHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	var user credentials.User
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = user.SetPassword(user.PasswordHash)
	if err != nil {
		reg.log.LogWarning("error while hashing password", err)
	}

	var token string
	idNew := uuid.New()
	idString := idNew.String()

	err = reg.setUser.UserSetter(ctx, user.Username, user.PasswordHash, idString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	token, err = authentication.BuildJWTString(idString)
	if err != nil {
		reg.log.LogWarning("err when get JWT token while registration:", err)
	}

	err = reg.RedisClient.Set(idString, token)
	if err != nil {
		reg.log.LogWarning("err when set value to redis in auth handler:", err)
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "You have successfully registered and authorized\n")
	//Нужно ли?
	reg.log.LogInfo("user", idString, "successfully registered and authorized")
}

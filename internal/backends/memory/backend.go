package memory

import (
	"sync"

	"github.com/nxshock/signaller/internal"
	"github.com/nxshock/signaller/internal/models"
	mSync "github.com/nxshock/signaller/internal/models/sync"
)

type Backend struct {
	data     map[string]internal.User
	rooms    map[string]internal.Room
	hostname string
	mutex    sync.Mutex // TODO: replace with RW mutex
}

type Token struct {
	Device string
}

func NewBackend(hostname string) *Backend {
	return &Backend{
		hostname: hostname,
		rooms:    make(map[string]internal.Room),
		data:     make(map[string]internal.User)}
}

func (backend *Backend) Register(username, password, device string) (user internal.User, token string, err *models.ApiError) {
	backend.mutex.Lock()
	defer backend.mutex.Unlock()

	if _, ok := backend.data[username]; ok {
		return nil, "", internal.NewError(models.M_USER_IN_USE, "trying to register a user ID which has been taken")
	}

	token = internal.NewToken(internal.DefaultTokenSize)

	user = &User{
		name:     username,
		password: password,
		Tokens: map[string]Token{
			token: {
				Device: device}},
		backend: backend}

	backend.data[username] = user

	return user, token, nil
}

func (backend *Backend) Login(username, password, device string) (token string, err *models.ApiError) {
	backend.mutex.Lock()
	defer backend.mutex.Unlock()

	user, ok := backend.data[username]
	if !ok {
		return "", internal.NewError(models.M_FORBIDDEN, "wrong username")
	}

	if user.Password() != password {
		return "", internal.NewError(models.M_FORBIDDEN, "wrong password")
	}

	token = internal.NewToken(internal.DefaultTokenSize)

	backend.data[username].(*User).Tokens[token] = Token{Device: device}

	return token, nil
}

func (backend *Backend) Logout(token string) *models.ApiError {
	backend.mutex.Lock()
	defer backend.mutex.Unlock()

	for _, user := range backend.data {
		for userToken, _ := range user.(*User).Tokens {
			if userToken == token {
				delete(user.(*User).Tokens, token)
				return nil
			}
		}
	}

	return internal.NewError(models.M_UNKNOWN_TOKEN, "unknown token") // TODO: create error struct
}

func (backend *Backend) Sync(token string, request mSync.SyncRequest) (response *mSync.SyncReply, err *models.ApiError) {
	backend.mutex.Lock()
	defer backend.mutex.Unlock()

	return nil, nil // TODO: implement
}
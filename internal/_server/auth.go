package server

import (
	"bytes"
	"crypto/rand"
	"net/http"
	"strings"

	"github.com/yklcs/panchro/internal/db"
	"golang.org/x/crypto/scrypt"
)

const saltBytes = 128

type Identity struct {
	ID   string
	Key  []byte
	Salt []byte
}

type AuthLayer struct {
	handler http.Handler
	authDB  *db.DB[Identity]
}

func NewAuthLayer(handler http.Handler, dbPath string) *AuthLayer {
	authDB := db.NewDB[Identity](dbPath)
	return &AuthLayer{
		authDB: authDB,
	}
}

func (l *AuthLayer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix(auth, "Bearer")
	if !ok {
		http.Error(w, "invalid bearer token", http.StatusUnauthorized)
		return
	}

	token = strings.TrimSpace(token)
	ok, err := l.VerifyIdentity(token)
	if err != nil || !ok {
		http.Error(w, "invalid bearer token", http.StatusUnauthorized)
	}

	l.handler.ServeHTTP(w, r)
}

func (l *AuthLayer) InitAuth() {
	l.authDB = db.NewDB[Identity]("auth.db.json")
}

func (l *AuthLayer) AddIdentity() (string, error) {
	salt := make([]byte, saltBytes)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	id, err := randBase58String(16)
	if err != nil {
		return "", err
	}
	token, err := randBase58String(16)
	if err != nil {
		return "", err
	}
	bearer := id + token

	key, err := scrypt.Key([]byte(token), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	identity := Identity{
		ID:   id,
		Key:  key,
		Salt: salt,
	}

	err = l.authDB.Put(id, identity)
	if err != nil {
		return "", err
	}

	return bearer, nil
}

func (l *AuthLayer) VerifyIdentity(bearer string) (bool, error) {
	id := bearer[:16]
	token := bearer[16:]
	identity, err := l.authDB.Get(id)
	if err != nil {
		return false, err
	}

	key, err := scrypt.Key([]byte(token), identity.Salt, 32768, 8, 1, 32)
	if err != nil {
		return false, err
	}

	if bytes.Equal(key, identity.Key) {
		return true, nil
	}

	return false, nil
}

func randBase58String(size int) (string, error) {
	const alphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := 0; i < size; i++ {
		b[i] = alphabet[b[i]%byte(len(alphabet))]
	}
	return string(b), nil
}

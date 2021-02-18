package unstable_api

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

//Token is the  base struct of session token
type Token struct {
	Name    string
	Value   string
	Expires time.Time
}
type DataBase struct {
	Token Token
	User  Account
}

//Store to map and mutex
type Store struct {
	M map[string]DataBase
	l sync.Mutex
}

func (store *Store) Init() {
	store.M = make(map[string]DataBase)
}

func (store *Store) Put(user Account, token Token) {
	store.M[token.Value] = DataBase{token, user}
}

func (store *Store) Remove(tokenValue string) {
	delete(store.M, tokenValue)
}

func (store *Store) Destruct(tokenValue string) {
	if store.M[tokenValue].Token.Expires == time.Now() {
		store.Remove(tokenValue)
	}
}
func (store *Store) Get(tokenValue string) (data DataBase, err error) {
	if data, ok := store.M[tokenValue]; ok {
		return data, nil
	}
	return DataBase{}, errors.New("not fount")
}

func GetToken(tokenString string) (data Token, err error) {

	err = json.Unmarshal([]byte(tokenString), &data)
	if err != nil {
		return
	}
	// b, _ := json.Marshal(data)
	return
}

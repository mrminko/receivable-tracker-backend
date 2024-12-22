package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mrminko/receivable-tracker/internal/database"
	"log"
	"net/http"
	"time"
)

func (Query *DBQuery) getAllUsers(w http.ResponseWriter, r *http.Request) {
	type UserJSON struct {
		Id        uuid.UUID      `json:"id"`
		Name      string         `json:"name"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		Phone     sql.NullString `json:"phone"`
	}
	users, err := Query.db.GetAllUsers(r.Context())
	if err != nil {
		log.Println("Error when querying users")
		return
	}
	var userList []UserJSON
	for _, user := range users {
		userList = append(userList, UserJSON{
			Id:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Phone:     user.Phone,
		})
	}
	respondWithJSON(w, 200, userList)
}

func (Query *DBQuery) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := fmt.Sprintf("Error when decoding data: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}
	user, err := Query.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Phone: sql.NullString{
			String: params.Phone,
			Valid:  true,
		},
	})
	if err != nil {
		errMsg := fmt.Sprintf("Error when writing to database: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}
	respondWithJSON(w, 201, user)
}

func (Query *DBQuery) deleteUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Id string `json:"id"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := fmt.Sprintf("Error when decoding data: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}
	userId, err := uuid.Parse(params.Id)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid id given: %v", err)
		respondWithError(w, 500, errMsg)
	}
	user, err := Query.db.DeleteUser(r.Context(), userId)
	if err != nil {
		errMsg := fmt.Sprintf("Error when deleting user: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}
	respondWithJSON(w, 201, user)
}

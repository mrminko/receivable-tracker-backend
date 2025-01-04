package main

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mrminko/receivable-tracker/internal/database"
	"github.com/mrminko/receivable-tracker/utils"
	"log"
	"net/http"
	"time"
)

type ReceivableStatus struct {
	Status string
	Valid  bool
}

func validateReceivableAndAssignStatus(receivable *database.UpdateReceivableParams) (errMsg error) {
	if receivable.AmountTotal == 0 {
		return fmt.Errorf("amount total should not be 0, delete the receivable instead")
	}

	if (receivable.AmountReceived < 0) || (receivable.AmountTotal < 0) || (receivable.AmountLeft < 0) {
		return fmt.Errorf("amounts should not be less than 0")
	}

	receivable.AmountLeft = receivable.AmountTotal - receivable.AmountReceived

	if receivable.AmountLeft <= 0 {
		receivable.Status = "closed"
		return nil
	} else if receivable.AmountLeft == receivable.AmountTotal {
		receivable.Status = "open"
		return nil
	} else {
		receivable.Status = "partial"
		return nil
	}

}

func calculateStatus(receivable *database.UpdateReceivableParams) (status string) {
	if receivable.AmountLeft == 0 {
		status = "closed"
	}
	if (receivable.AmountReceived < receivable.AmountTotal) && (receivable.AmountReceived != 0) {
		status = "partial"
	}
	if receivable.AmountReceived == 0 {
		status = "open"
	}
	return status
}

func validateStatus(status string) ReceivableStatus {
	switch status {
	case "open", "partial", "closed":
		return ReceivableStatus{
			Status: status,
			Valid:  true,
		}
	case "":
		return ReceivableStatus{
			Status: "open",
			Valid:  true,
		}
	default:
		return ReceivableStatus{
			Status: "",
			Valid:  false,
		}
	}

}

func (Query *DBQuery) getAllReceivables(w http.ResponseWriter, r *http.Request) {
	type ReceivableJSON struct {
		Id             uuid.UUID `json:"id"`
		UserId         uuid.UUID `json:"user_id"`
		UserName       string    `json:"user_name"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		AmountTotal    float64   `json:"amount_total"`
		AmountReceived float64   `json:"amount_received"`
		AmountLeft     float64   `json:"amount_left"`
		Status         string    `json:"status"`
	}
	receivables, err := Query.db.GetAllReceivables(r.Context())
	if err != nil {
		log.Println("Error when querying receivables")
		return
	}
	var receivableList []ReceivableJSON
	for _, receivable := range receivables {
		receivableList = append(receivableList, ReceivableJSON{
			Id:             receivable.ID,
			UserId:         receivable.Userid,
			UserName:       receivable.Username,
			CreatedAt:      receivable.CreatedAt,
			UpdatedAt:      receivable.UpdatedAt,
			AmountTotal:    receivable.AmountTotal,
			AmountReceived: receivable.AmountReceived,
			AmountLeft:     receivable.AmountLeft,
			Status:         receivable.Status,
		})
	}
	respondWithJSON(w, 200, receivableList)
}

func (Query *DBQuery) createReceivable(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserId         uuid.UUID `json:"user_id"`
		Date           string    `json:"date"`
		AmountTotal    float64   `json:"amount_total"`
		AmountReceived float64   `json:"amount_received,omitempty"`
		AmountLeft     float64   `json:"amount_left,omitempty"`
		Status         string    `json:"status,omitempty"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := fmt.Sprintf("Error when decoding data: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}

	dateTime, err := utils.ParseTime(params.Date)
	if err != nil {
		errMsg := fmt.Sprintln("Error when parsing date. Please provide in the format \"02-01-2006\"")
		respondWithError(w, 500, errMsg)
		return
	}

	status := validateStatus(params.Status)
	if !status.Valid {
		errMsg := fmt.Sprintf("Status field must be one of \"open\", \"partial\", \"closed\"")
		respondWithError(w, 500, errMsg)
		return
	}

	receivable, err := Query.db.CreateReceivable(r.Context(), database.CreateReceivableParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Userid:         params.UserId,
		Date:           dateTime,
		AmountTotal:    params.AmountTotal,
		AmountReceived: params.AmountReceived,
		AmountLeft:     params.AmountLeft,
		Status:         status.Status,
	})
	if err != nil {
		errMsg := fmt.Sprintf("Error when creating receivable: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}
	respondWithJSON(w, 201, receivable)
}

func (Query *DBQuery) deleteReceivable(w http.ResponseWriter, r *http.Request) {
	receivableId, err := utils.StringToUUID(r, "receivableId")
	if err != nil {
		errMsg := fmt.Sprintf("Invalid receivable id given: %v", err)
		respondWithError(w, 500, errMsg)
	}
	receivable, err := Query.db.DeleteReceivable(r.Context(), receivableId)
	if err != nil {
		errMsg := fmt.Sprintf("Error when deleting user: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}
	respondWithJSON(w, 200, receivable)
}

func (Query *DBQuery) updateReceivable(w http.ResponseWriter, r *http.Request) {
	receivableId, err := utils.StringToUUID(r, "receivableId")
	if err != nil {
		errMsg := fmt.Sprintf("Invalid receivable id given: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}
	existingReceivable, err := Query.db.GetReceivableByID(r.Context(), receivableId)
	if err != nil {
		errMsg := fmt.Sprintf("Entry does not exist: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}

	type parameters struct {
		Date           string   `json:"date,omitempty"`
		AmountTotal    *float64 `json:"amount_total"`
		AmountReceived *float64 `json:"amount_received"`
	}

	params := &parameters{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(params)
	if err != nil {
		errMsg := fmt.Sprintf("Error when decoding data: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}

	UpdatedParams := &database.UpdateReceivableParams{}

	UpdatedParams.ID = existingReceivable.ID

	if params.Date != "" {
		dateTime, err := utils.ParseTime(params.Date)
		if err != nil {
			errMsg := fmt.Sprintln("Error when parsing date. Please provide in the format \"02-01-2006\"")
			respondWithError(w, 500, errMsg)
			return
		}
		UpdatedParams.Date = dateTime
	} else {
		UpdatedParams.Date = existingReceivable.Date
	}

	if params.AmountReceived != nil {
		UpdatedParams.AmountReceived = *params.AmountReceived
	} else {
		UpdatedParams.AmountReceived = existingReceivable.AmountReceived
	}

	if params.AmountTotal != nil {
		UpdatedParams.AmountTotal = *params.AmountTotal
	} else {
		UpdatedParams.AmountTotal = existingReceivable.AmountTotal
	}

	err = validateReceivableAndAssignStatus(UpdatedParams)
	if err != nil {
		errMsg := fmt.Sprintf("Error when validating receivable: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}

	updatedReceivable, err := Query.db.UpdateReceivable(r.Context(), *UpdatedParams)
	if err != nil {
		errMsg := fmt.Sprintf("Error when updating receivable: %v", err)
		respondWithError(w, 500, errMsg)
		return
	}

	respondWithJSON(w, 200, updatedReceivable)
}

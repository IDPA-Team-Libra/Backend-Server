package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Liberatys/libra-back/main/performance"

	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/sec"
	"github.com/Liberatys/libra-back/main/user"
)

//GetUserPerformance loads the performance for a single user
func GetUserPerformance(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		return
	}
	var currentUser User
	err = json.Unmarshal(body, &currentUser)
	if err != nil {
		logger.LogMessage("Anfrage an Performance hatte invalides JSON", logger.INFO)
		w.Write([]byte("Invalid json"))
		return
	}
	validator := sec.NewTokenValidator(currentUser.AccessToken, currentUser.Username)
	if validator.IsValidToken(jwtKey) == false {
		logger.LogMessageWithOrigin(fmt.Sprintf("User: %s was not able to be authenticated", currentUser.Username), logger.WARNING, "performanceRoute.go | GetUserPerformance")
		type Response struct {
			Message string
		}
		obj, _ := json.Marshal(Response{
			Message: "Was not able to validate user",
		})
		w.Write([]byte(obj))
		return
	}
	userID := user.GetUserIDByUsername(currentUser.Username, GetDatabaseInstance())
	performances := performance.LoadUserPerformance(userID, GetDatabaseInstance())
	obj, _ := json.Marshal(performances)
	w.Write(obj)
}

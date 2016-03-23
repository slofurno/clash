package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/slofurno/front/datastore"
	"github.com/slofurno/front/utils"
	"github.com/slofurno/ws"
	"net/http"
)

func createRoom(res http.ResponseWriter, req *http.Request) {

}

func createGame(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	cr := &clashRequest{}

	if decoder.Decode(cr) != nil {
		return
	}

	clash := NewClash(cr.challenge)
	store.Clashes.Insert(clash)

	w.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(w).Encode(clash)

}

func getRooms(w http.ResponseWriter, r *http.Request) {

}

func joinGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["game"]

	_ = r.Header.Get("Authorization")

	//TODO: lookup user from token
}

func postCode(res http.ResponseWriter, req *http.Request) {
	auth := req.Header.Get("Authorization")

	if auth == "" {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	id := utils.Makeid()
	code := &datastore.Code{}
	json.NewDecoder(req.Body).Decode(code)

	fmt.Println(code)

	code.Id = id
	code.User = "esteban"
	code.Time = utils.Epoch_ms()

	store.Codes.Insert(code)
	store.CodeRunner.Push(code.Id)

	res.Write([]byte(id))

	//lookup user via auth token

}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	_ = ws.Upgrade(w, r)

}

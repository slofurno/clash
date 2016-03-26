package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slofurno/front/datastore"
	"github.com/slofurno/front/utils"
	"github.com/slofurno/ws"

	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createLogin(res http.ResponseWriter, req *http.Request) {
	login := &loginRequest{}
	err := json.NewDecoder(req.Body).Decode(login)

	if err != nil {
		//TODO: return an error code or something
		return
	}

	matches := store.Accounts.Get(login.Email)
	var account *datastore.Account

	for _, match := range matches {
		err := bcrypt.CompareHashAndPassword([]byte(match.Password), []byte(login.Password))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		account = match
		break
	}

	if account == nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	newlogin := datastore.NewLogin(account)
	store.Logins.Insert(newlogin)

	res.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(res).Encode(newlogin)
}

func createAccount(res http.ResponseWriter, req *http.Request) {
	a := &datastore.Account{}
	err := json.NewDecoder(req.Body).Decode(a)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	account := datastore.NewAccount(a.Email, a.Password)

	if account == nil {
		return
	}

	store.Accounts.Insert(account)
	login := datastore.NewLogin(account)
	store.Logins.Insert(login)

	res.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(res).Encode(login)
}

func getProblems(res http.ResponseWriter, req *http.Request) {
	problems := store.Problems.Query()

	res.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(res).Encode(problems)
}

func getProblem(res http.ResponseWriter, req *http.Request) {
	problemid := mux.Vars(req)["problem"]
	problem := store.Problems.Get(problemid)
	res.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(res).Encode(problem)
}

func postProblem(res http.ResponseWriter, req *http.Request) {

	problem := &datastore.Problem{}
	err := json.NewDecoder(req.Body).Decode(problem)

	if err != nil {
		return
	}

	problem.Id = utils.Makeid()
	store.Problems.Insert(problem)

	res.Write([]byte(problem.Id))
}

func createRoom(res http.ResponseWriter, req *http.Request) {
	room := &datastore.Room{}
	err := json.NewDecoder(req.Body).Decode(room)

	if err != nil {
		return
	}

	room.Id = utils.Makeid()
	room.Time = utils.Epoch_ms()
	store.Rooms.Insert(room)

	json.NewEncoder(res).Encode(room)
}

func createClash(w http.ResponseWriter, r *http.Request) {
	room := mux.Vars(r)["room"]
	decoder := json.NewDecoder(r.Body)
	cr := &clashRequest{}

	if decoder.Decode(cr) != nil {
		return
	}

	//TODO: maybe make sure you own the room
	_, err := auth(r)

	if err != nil {
		return
	}

	clash := NewClash(cr.Problem)
	fmt.Println(clash)
	store.Clashes.Insert(clash)

	store.Events.Insert(&datastore.Event{
		Id:      utils.Makeid(),
		Time:    utils.Epoch_ms(),
		Subject: room,
		Noun:    clash.Id,
		Verb:    "STARTED_CLASH",
	})

	w.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(w).Encode(clash)
}

func getClash(res http.ResponseWriter, req *http.Request) {
	clashid := mux.Vars(req)["clash"]
	clash := store.Clashes.Get(clashid)
	res.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(res).Encode(clash)
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	subject := mux.Vars(r)["subject"]
	events := store.Events.Query(subject)
	w.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(w).Encode(events)
}

type CodePost struct {
	Clash     string
	User      string
	Code      string
	Signature string
}

func getResults(res http.ResponseWriter, req *http.Request) {
	clash := mux.Vars(req)["clash"]
	results := store.Results.Get(clash)
	res.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(res).Encode(results)
}

func postResult(w http.ResponseWriter, r *http.Request) {
	//TODO: match auth to submitter of code?

	clashid := mux.Vars(r)["clash"]
	codeid := mux.Vars(r)["code"]

	code := store.Codes.Get(codeid)

	if code.Clash != clashid {
		fmt.Println("clash doesn't mtch")
		return
	}

	//clash := store.Clashes.Get(clashid)
	resultid := utils.Makeid()

	//TODO: default code result is success
	store.Results.Insert(&datastore.Result{
		Id:     resultid,
		Clash:  clashid,
		Status: code.Status,
		Time:   code.Time,
		User:   code.User,
		Code:   code.Id,
	})

	w.Write([]byte(resultid))
}

func getCode(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	x := store.Codes.Get(code)
	w.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(w).Encode(x)
}

func getRooms(res http.ResponseWriter, req *http.Request) {
	rooms := store.Rooms.Get()

	if rooms == nil {
		return
	}
	res.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(res).Encode(rooms)
}

func joinRoom(res http.ResponseWriter, req *http.Request) {
	room := mux.Vars(req)["room"]
	login, err := auth(req)

	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	store.Events.Insert(&datastore.Event{
		Id:      utils.Makeid(),
		Time:    utils.Epoch_ms(),
		Subject: room,
		Noun:    login.Account,
		Verb:    "JOINED_LOBBY",
	})

}

func joinGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["game"]

	_ = r.Header.Get("Authorization")

	//TODO: lookup user from token
}

func auth(req *http.Request) (*datastore.Login, error) {
	auth := req.Header.Get("Authorization")

	if auth == "" {
		return nil, errors.New("auth")
	}

	login := store.Logins.Get(auth)

	if login == nil {
		return nil, errors.New("auth")
	}

	return login, nil
}

func postCode(res http.ResponseWriter, req *http.Request) {
	clashid := mux.Vars(req)["clash"]
	clash := store.Clashes.Get(clashid)

	auth := req.Header.Get("Authorization")

	if auth == "" {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	login := store.Logins.Get(auth)
	fmt.Println("authed as:", login)

	id := utils.Makeid()
	code := &datastore.Code{}
	json.NewDecoder(req.Body).Decode(code)

	code.Id = id
	code.User = login.Account
	code.Time = utils.Epoch_ms()
	code.Clash = clashid
	code.Problem = clash.Problem
	code.Status = -1

	store.Codes.Insert(code)
	store.CodeRunner.Push(code.Id)

	res.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(res).Encode(code)
}

type Change struct {
	Type    string `json:"type"`
	Subject string `json:"subject"`
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	sock := ws.Upgrade(w, r)
	defer sock.Close()

	handle := obs.Subscribe()

	go func() {
		obs.AddSubject(handle, "test")
		for m := range handle.c {
			sock.WriteS(m)
		}
	}()

	for {
		m, code, err := sock.Read()
		fmt.Println(m)

		if err != nil || code == ws.Close {
			break
		}

		change := &Change{}
		err = json.Unmarshal([]byte(m), change)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println(change)
		switch change.Type {
		case "SUB":
			obs.AddSubject(handle, change.Subject)
		case "UNSUB":
			obs.RemoveSubject(handle, change.Subject)
		}
	}

	fmt.Println("ws quit")
	obs.UnSubscribe(handle)
}

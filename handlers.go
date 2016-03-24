package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slofurno/front/datastore"
	"github.com/slofurno/front/utils"
	"github.com/slofurno/ws"
)

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

	event := &datastore.Event{
		Id:      utils.Makeid(),
		Subject: "test",
		Noun:    "steve",
		Verb:    "joined",
		Time:    utils.Epoch_ms(),
	}

	store.Events.Insert(event)
}

func createClash(w http.ResponseWriter, r *http.Request) {

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

	//TODO: default code result is success
	store.Results.Insert(&datastore.Result{
		Id:     utils.Makeid(),
		Clash:  clashid,
		Status: code.Status,
		Time:   code.Time,
		User:   code.User,
	})
}

func getCode(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	x := store.Codes.Get(code)
	w.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(w).Encode(x)
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
	clashid := mux.Vars(req)["clash"]
	clash := store.Clashes.Get(clashid)

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
	code.Clash = clashid
	code.Problem = clash.Problem

	store.Codes.Insert(code)
	store.CodeRunner.Push(code.Id)

	res.Write([]byte(id))
	//lookup user via auth token
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

		if err != nil || code == ws.Close {
			break
		}

		change := &Change{}
		err = json.Unmarshal([]byte(m), change)

		if err == nil {
			fmt.Println(change)
		}
	}

	fmt.Println("ws quit")
	obs.UnSubscribe(handle)
}

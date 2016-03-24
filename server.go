package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/slofurno/front/datastore"
	"github.com/slofurno/front/utils"

	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var SubscriptionConfirmation = "SubscriptionConfirmation"
var Notification = "Notification"

var store *datastore.DataStore
var obs *SnsObserver

type Subscription struct {
	subjects map[string]bool
	c        chan string
}

type SnsObserver struct {
	//subs map[string][]chan struct{}
	subs []*Subscription
	m    *sync.Mutex
}

func (s *SnsObserver) Subscribe() *Subscription {
	tevs := &Subscription{
		c:        make(chan string, 16),
		subjects: map[string]bool{},
	}

	s.m.Lock()
	defer s.m.Unlock()

	s.subs = append(s.subs, tevs)
	return tevs
}

func (s *SnsObserver) AddSubject(sub *Subscription, subject string) {
	s.m.Lock()
	defer s.m.Unlock()

	sub.subjects[subject] = true
}

func (s *SnsObserver) Publish(subject, message string) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, sub := range s.subs {
		if !sub.subjects[subject] {
			continue
		}

		select {
		case sub.c <- message:
		default:
			fmt.Println("dropped message")
		}
	}
}

func (s *SnsObserver) RemoveSubject(sub *Subscription, subject string) {
	s.m.Lock()
	defer s.m.Unlock()

	sub.subjects[subject] = false
}

func (s *SnsObserver) UnSubscribe(sub *Subscription) {
	s.m.Lock()
	defer s.m.Unlock()
	subs := s.subs

	for i := 0; i < len(subs); i++ {
		if subs[i] == sub {
			s.subs = append(subs[:i], subs[i+1:]...)
		}
	}

	close(sub.c)
}

type CodeSubmission struct {
	Code   string `json:"code"`
	Runner string `json:"runner"`
}

type clashRequest struct {
	challenge string `json:"challenge"`
}

type Submission struct {
	Id     string `json:"id"`
	User   string `json:"user"`
	Clash  string `json:"clash"`
	Code   string `json:"code"`
	Output string `json:"output"`
	Diff   string `json:"diff"`
	Runner string `json:"runner"`
	Test   string `json:"test"`
}

type SubmissionRequest struct {
	Id        string `json:"id"`
	User      string `json:"user"`
	Clash     string `json:"clash"`
	Code      string `json:"code"`
	Runner    string `json:"runner"`
	Signature string `json:"signature"`
	Challenge string `json:"challenge"`
}

func SignSubmission(sr *SubmissionRequest) string {
	mac := hmac.New(sha256.New, []byte("thisshouldbeasecretkey"))
	mac.Write([]byte(sr.Id))
	mac.Write([]byte(sr.User))
	mac.Write([]byte(sr.Clash))
	mac.Write([]byte(sr.Challenge))
	sum := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

func NewClash(challenge string) *datastore.Clash {
	return &datastore.Clash{
		Id:      utils.Makeid(),
		Time:    utils.Epoch_ms(),
		Problem: challenge,
	}
}

type SnsMessage struct {
	SubscribeURL string
	Type         string
	Subject      string
	Message      string
}

func handleSns(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	m := &SnsMessage{}
	decoder.Decode(m)

	if m.Type == SubscriptionConfirmation {
		http.Get(m.SubscribeURL)
		return
	}

	obs.Publish(m.Subject, m.Message)
}

func main() {

	obs = &SnsObserver{
		m:    &sync.Mutex{},
		subs: []*Subscription{},
	}

	sess := session.New(&aws.Config{Region: aws.String("us-east-1")})
	store = datastore.New()

	/*
		event := &datastore.Event{
			Id:      utils.Makeid(),
			Subject: "othertest",
			Noun:    "steve",
			Verb:    "joined",
		}

		store.Events.Insert(event)

		store.Events.Query("othertest")
	*/

	res, err := http.Get("http://ifconfig.co")

	if err != nil {
		panic(err.Error())
	}

	b, _ := ioutil.ReadAll(res.Body)
	addr := strings.TrimSpace(string(b))
	endpoint := "http://" + addr + ":5555/api/sns"
	fmt.Println(endpoint)

	ps := sns.New(sess)

	_, err = ps.Subscribe(&sns.SubscribeInput{
		Endpoint: aws.String(endpoint),
		TopicArn: aws.String("arn:aws:sns:us-east-1:027082628651:clash_events"),
		Protocol: aws.String("http"),
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	sub := obs.Subscribe()

	go func() {
		obs.AddSubject(sub, "test")

		for m := range sub.c {
			fmt.Println(m)
		}
	}()

	/*
		sub2, err := ps.ConfirmSubscription(&sns.ConfirmSubscriptionInput{
			TopicArn: aws.String("arn:aws:sns:us-east-1:027082628651:clash_events"),
			Token:    aws.String("whatsmytoken"),
		})

		fmt.Println(sub2.GoString())
	*/

	/*
		ps.ConfirmSubscription(&sns.ConfirmSubscriptionInput{

		})
	*/

	/*
		vals := map[string]*dynamodb.AttributeValue{
			":player": &dynamodb.AttributeValue{
				L: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{
						S: aws.String("steve"),
					},
				},
			},
		}

		key := map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{S: aws.String("steve")},
		}

		out, err := ddb.UpdateItem(&dynamodb.UpdateItemInput{
			TableName:                 aws.String("clash"),
			UpdateExpression:          aws.String("SET clashes = list_append(clashes, :player)"),
			ExpressionAttributeValues: vals,
			Key: key,
		})

		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(out.GoString())

	*/

	/*
		bucket := s3.New(config)
		res, err := bucket.PutObject(&s3.PutObjectInput{})


		ddb := dynamodb.New(config)
		ddb.UpdateItem(&dynamodb.UpdateItemInput{
			TableName: "clash",

		})

		if err != nil {
			fmt.Println(err.Error())
		}

	*/
	r := mux.NewRouter()
	r.HandleFunc("/api/clash", createClash).Methods("POST")
	r.HandleFunc("/api/clash/{clash}", postCode).Methods("POST")
	r.HandleFunc("/api/event", createRoom).Methods("POST")
	r.HandleFunc("/api/events/{subject}", getEvents).Methods("GET")
	r.HandleFunc("/api/rooms", getRooms).Methods("GET")

	r.HandleFunc("/api/problems", postProblem).Methods("POST")
	r.HandleFunc("/api/problems", getProblems).Methods("GET")
	r.HandleFunc("/api/problems/{problem}", getProblem).Methods("GET")

	r.HandleFunc("/api/code/{code}", getCode).Methods("GET")
	r.HandleFunc("/api/clash/{clash}/code/{code}", postResult).Methods("POST")

	r.HandleFunc("/api/ws", websocketHandler)
	r.HandleFunc("/api/sns", handleSns).Methods("POST")

	http.Handle("/api/", r)
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":5555", nil)
}

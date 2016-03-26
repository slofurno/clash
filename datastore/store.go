package datastore

import (
	"encoding/json"
	"fmt"
	"github.com/slofurno/front/utils"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"

	"golang.org/x/crypto/bcrypt"
)

const queue string = "https://queue.amazonaws.com/027082628651/myqueue"
const CODE_TABLE = "code"
const CLASH_TABLE = "clashe"
const EVENT_TABLE = "event"
const RESULTS_TABLE = "results"
const PROBLEMS_TABLE = "problems"
const ACCOUNTS_TABLE = "accounts"
const LOGINS_TABLE = "logins"
const ROOMS_TABLE = "rooms"
const CODE_RUNNER = "https://sqs.us-east-1.amazonaws.com/027082628651/coderunner"
const EVENT_TOPIC = "arn:aws:sns:us-east-1:027082628651:clash_events"

type Code struct {
	Id      string `json:"id"`
	Code    string `json:"code"`
	Clash   string `json:"clash"`
	Runner  string `json:"runner"`
	Problem string `json:"problem"`
	User    string `json:"user"`
	Time    int64  `json:"time"`
	Diff    string `json:"diff"`
	Status  int64  `json:"status"`
	Output  string `json:"output"`
}

type Result struct {
	Id     string `json:"id"`
	User   string `json:"user"`
	Clash  string `json:"clash"`
	Time   int64  `json:"time"`
	Status int64  `json:"status"`
	Code   string `json:"code"`
}

type Clash struct {
	Id      string `json:"id"`
	Time    int64  `json:"time"`
	Problem string `json:"problem"`
}

type Problem struct {
	Id     string `json:"id"`
	Text   string `json:"text"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Account struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
}

func NewAccount(email, password string) *Account {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	return &Account{
		Email:    email,
		Password: string(hashed),
		Id:       utils.Makeid(),
	}
}

type Login struct {
	Id      string `json:"id"`
	Account string `json:"account"`
	Token   string `json:"token"`
}

func NewLogin(account *Account) *Login {
	return &Login{
		Id:      utils.Makeid(),
		Token:   utils.Makeid(),
		Account: account.Id,
	}
}

type DataStore struct {
	Clashes    *ClashStore
	Codes      *CodeStore
	Rooms      *RoomStore
	CodeRunner *CodeRunner
	Events     *EventStore
	Results    *ResultStore
	Problems   *ProblemStore
	Accounts   *AccountStore
	Logins     *LoginStore
}

type EventPusher struct {
	pusher *sns.SNS
}

func (s *EventPusher) Publish(event *Event) {

	b, err := json.Marshal(event)

	if err != nil {
		return
	}

	s.pusher.Publish(&sns.PublishInput{
		TopicArn: aws.String(EVENT_TOPIC),
		Subject:  aws.String(event.Subject),
		Message:  aws.String(string(b)),
	})
}

func New() *DataStore {

	sess := session.New(&aws.Config{Region: aws.String("us-east-1")})

	ddb := dynamodb.New(sess)
	mysqs := sqs.New(sess)
	pusher := sns.New(sess)

	clashes := &ClashStore{db: ddb}
	codes := &CodeStore{db: ddb}
	coderunner := &CodeRunner{queue: mysqs}
	events := &EventStore{db: ddb, pub: pusher}
	rooms := &RoomStore{db: ddb}
	results := &ResultStore{db: ddb}
	problems := &ProblemStore{db: ddb}
	accounts := &AccountStore{db: ddb}
	logins := &LoginStore{db: ddb}

	return &DataStore{
		Clashes:    clashes,
		Codes:      codes,
		CodeRunner: coderunner,
		Events:     events,
		Rooms:      rooms,
		Results:    results,
		Problems:   problems,
		Accounts:   accounts,
		Logins:     logins,
	}
}

type Room struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
	Id   string `json:"id"`
}

type ClashStore struct {
	db *dynamodb.DynamoDB
}

type CodeStore struct {
	db *dynamodb.DynamoDB
}

type EventStore struct {
	db  *dynamodb.DynamoDB
	pub *sns.SNS
}

type RoomStore struct {
	db *dynamodb.DynamoDB
}

type ResultStore struct {
	db *dynamodb.DynamoDB
}

type CodeRunner struct {
	queue *sqs.SQS
}

type ProblemStore struct {
	db *dynamodb.DynamoDB
}

type AccountStore struct {
	db *dynamodb.DynamoDB
}

type LoginStore struct {
	db *dynamodb.DynamoDB
}

type Event struct {
	Id      string `json:"id"`
	Subject string `json:"subject"`
	Noun    string `json:"noun"`
	Verb    string `json:"verb"`
	Time    int64  `json:"time"`
}

type Event2 struct {
	Topic string      `json:"topic"`
	Load  interface{} `json:"load"`
	Meta  string      `json:"meta"`
}

func (s *AccountStore) Insert(account *Account) {

	b, _ := json.Marshal(account.Meta)

	item := map[string]*dynamodb.AttributeValue{
		"id":       S(account.Id),
		"email":    S(account.Email),
		"password": S(account.Password),
		"meta":     S(string(b)),
	}

	s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(ACCOUNTS_TABLE),
		Item:      item,
	})
}

func (s *AccountStore) PutMeta(id, meta string) {
	update := map[string]*dynamodb.AttributeValue{
		":meta": S(meta),
	}
	key := map[string]*dynamodb.AttributeValue{
		"id": S(id),
	}
	out, err := s.db.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:                 aws.String(ACCOUNTS_TABLE),
		Key:                       key,
		UpdateExpression:          aws.String("SET meta = :meta"),
		ExpressionAttributeValues: update,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(out.String())
}

func (s *AccountStore) GetMeta(id string) string {
	key := map[string]*dynamodb.AttributeValue{
		"id": S(id),
	}

	out, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ACCOUNTS_TABLE),
		Key:       key,
	})

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	meta := out.Item["meta"]

	if meta == nil {
		return ""
	}
	return *meta.S
}

func (s *AccountStore) Query(loginEmail string) []*Account {

	item := map[string]*dynamodb.AttributeValue{
		":email": S(loginEmail),
	}

	out, err := s.db.Scan(&dynamodb.ScanInput{
		TableName:                 aws.String(ACCOUNTS_TABLE),
		FilterExpression:          aws.String("email = :email"),
		ExpressionAttributeValues: item,
	})

	matches := []*Account{}
	if err != nil {
		fmt.Println(err.Error())
		return matches
	}

	for _, item := range out.Items {
		email := item["email"]
		password := item["password"]
		id := item["id"]

		if email == nil || password == nil || id == nil {
			continue
		}

		matches = append(matches, &Account{
			Id:       *id.S,
			Email:    *email.S,
			Password: *password.S,
		})
	}

	return matches
}

func (s *LoginStore) Insert(login *Login) {
	item := map[string]*dynamodb.AttributeValue{
		"id":      S(login.Id),
		"account": S(login.Account),
		"token":   S(login.Token),
	}
	s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(LOGINS_TABLE),
		Item:      item,
	})

}

func (s *LoginStore) Get(token string) *Login {
	key := map[string]*dynamodb.AttributeValue{
		"token": S(token),
	}
	res, _ := s.db.GetItem(&dynamodb.GetItemInput{
		ConsistentRead: aws.Bool(true),
		TableName:      aws.String(LOGINS_TABLE),
		Key:            key,
	})

	item := res.Item
	id := item["id"]
	account := item["account"]
	expectedToken := item["token"]

	if expectedToken == nil || id == nil || account == nil {
		return nil
	}

	return &Login{
		Id:      *id.S,
		Account: *account.S,
		Token:   token,
	}
}

func (s *ProblemStore) Insert(problem *Problem) {
	d, _ := json.Marshal(problem)

	item := map[string]*dynamodb.AttributeValue{
		"id":   S(problem.Id),
		"json": S(string(d)),
	}

	s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(PROBLEMS_TABLE),
		Item:      item,
	})

}

func (s *ProblemStore) Get(id string) *Problem {
	key := map[string]*dynamodb.AttributeValue{
		"id": S(id),
	}

	res, _ := s.db.GetItem(&dynamodb.GetItemInput{
		ConsistentRead: aws.Bool(true),
		TableName:      aws.String(PROBLEMS_TABLE),
		Key:            key,
	})

	item := res.Item["json"]

	if item == nil {
		return nil
	}
	d := *item.S
	problem := &Problem{}
	json.Unmarshal([]byte(d), problem)
	return problem
}

func (s *ProblemStore) Query() []*Problem {
	results := []*Problem{}

	res, err := s.db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(PROBLEMS_TABLE),
	})

	if err != nil {
		fmt.Println(err.Error())
		return results
	}

	for _, item := range res.Items {
		d := *item["json"].S
		problem := &Problem{}
		err := json.Unmarshal([]byte(d), problem)

		if err != nil {
			continue
		}

		results = append(results, problem)
	}

	return results
}

func (s *ResultStore) Insert(result *Result) {

	d, err := json.Marshal(result)

	item := map[string]*dynamodb.AttributeValue{
		"id":    S(result.Id),
		"clash": S(result.Clash),
		"json":  S(string(d)),
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(RESULTS_TABLE),
		Item:      item,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (s *ResultStore) Get(clash string) []*Result {

	item := map[string]*dynamodb.AttributeValue{
		":clash": S(clash),
	}

	out, _ := s.db.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(RESULTS_TABLE),
		KeyConditionExpression:    aws.String("clash = :clash"),
		ExpressionAttributeValues: item,
	})

	results := []*Result{}

	for _, item := range out.Items {
		d := item["json"]
		if d == nil {
			continue
		}

		result := &Result{}

		err := json.Unmarshal([]byte(*d.S), result)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		results = append(results, result)
	}

	return results

}

func (s *EventStore) Insert(event *Event) {
	b, err := json.Marshal(event)
	if err != nil {
		return
	}

	item := map[string]*dynamodb.AttributeValue{
		"id":      S(event.Id),
		"subject": S(event.Subject),
		"noun":    S(event.Noun),
		"verb":    S(event.Verb),
		"time":    N(event.Time),
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(EVENT_TABLE),
		Item:      item,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s.pub.Publish(&sns.PublishInput{
		TopicArn: aws.String(EVENT_TOPIC),
		Subject:  aws.String(event.Subject),
		Message:  aws.String(string(b)),
	})
}

func (s *EventStore) Query(subject string) []*Event {

	item := map[string]*dynamodb.AttributeValue{
		":sub": S(subject),
	}

	out, err := s.db.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(EVENT_TABLE),
		KeyConditionExpression:    aws.String("subject = :sub"),
		ExpressionAttributeValues: item,
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	events := []*Event{}

	for _, item := range out.Items {

		verb := item["verb"]
		subject := item["subject"]
		noun := item["noun"]
		id := item["id"]
		time := item["time"]

		if verb == nil || subject == nil || noun == nil || id == nil || time == nil {
			continue
		}

		i, _ := strconv.ParseInt(*time.N, 10, 64)

		event := &Event{
			Id:      *id.S,
			Noun:    *noun.S,
			Subject: *subject.S,
			Verb:    *verb.S,
			Time:    i,
		}

		events = append(events, event)
	}

	return events
}

func (s *CodeStore) Insert(code *Code) {

	item := map[string]*dynamodb.AttributeValue{
		"id":      S(code.Id),
		"user":    S(code.User),
		"time":    N(code.Time),
		"code":    S(code.Code),
		"problem": S(code.Problem),
		"runner":  S(code.Runner),
		"diff":    S(code.Diff),
		"status":  N(code.Status),
		"clash":   S(code.Clash),
		"output":  S(code.Output),
	}

	_, err := s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(CODE_TABLE),
		Item:      item,
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}

func (s *CodeStore) Get(id string) *Code {
	key := map[string]*dynamodb.AttributeValue{
		"id": S(id),
	}

	res, err := s.db.GetItem(&dynamodb.GetItemInput{
		ConsistentRead: aws.Bool(true),
		TableName:      aws.String(CODE_TABLE),
		Key:            key,
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	X := res.Item
	ret := &Code{}

	code := X["code"]
	problem := X["problem"]
	runner := X["runner"]
	user := X["user"]
	time := X["time"]
	clash := X["clash"]

	if code == nil || problem == nil || runner == nil || user == nil || time == nil || clash == nil {
		return nil
	}

	ret.Code = *code.S
	ret.Problem = *problem.S
	ret.Runner = *runner.S
	ret.User = *user.S
	ret.Clash = *clash.S
	ret.Id = id

	i, err := strconv.ParseInt(*time.N, 10, 64)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		ret.Time = i
	}

	diff := X["diff"]
	status := X["status"]
	output := X["output"]

	if diff != nil && status != nil {
		i, err := strconv.ParseInt(*status.N, 10, 64)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			ret.Status = i
		}

		ret.Diff = *diff.S
	}

	if output != nil {
		ret.Output = *output.S
	}
	return ret
}

func (s *RoomStore) Insert(room *Room) {
	item := map[string]*dynamodb.AttributeValue{
		"id":   S(room.Id),
		"time": N(room.Time),
		"name": S(room.Name),
	}

	_, err := s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(ROOMS_TABLE),
		Item:      item,
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}

func (s *RoomStore) Get() []*Room {
	out, err := s.db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(ROOMS_TABLE),
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	rooms := []*Room{}
	for _, item := range out.Items {
		id := item["id"]
		time := item["time"]
		name := item["name"]

		if id == nil || time == nil || name == nil {
			continue
		}

		i, _ := strconv.ParseInt(*time.N, 10, 64)
		rooms = append(rooms, &Room{
			Id:   *id.S,
			Time: i,
			Name: *name.S,
		})
	}

	return rooms
}

func (s *ClashStore) Get(clashid string) *Clash {

	key := map[string]*dynamodb.AttributeValue{
		"id": S(clashid),
	}

	res, err := s.db.GetItem(&dynamodb.GetItemInput{
		ConsistentRead: aws.Bool(true),
		TableName:      aws.String(CLASH_TABLE),
		Key:            key,
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	d := res.Item["json"]
	b := []byte(*d.S)

	clash := &Clash{}
	json.Unmarshal(b, clash)

	return clash
}

func (s *ClashStore) Insert(clash *Clash) (*dynamodb.PutItemOutput, error) {

	b, _ := json.Marshal(clash)

	items := map[string]*dynamodb.AttributeValue{
		"id":   S(clash.Id),
		"json": S(string(b)),
		"time": N(clash.Time),
	}

	return s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(CLASH_TABLE),
		Item:      items,
	})
}

func (s *CodeRunner) Receive() (string, *string) {
	res, err := s.queue.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(CODE_RUNNER),
		WaitTimeSeconds:     aws.Int64(20),
		MaxNumberOfMessages: aws.Int64(1),
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	messages := res.Messages

	if len(messages) != 1 {
		return "", nil
	}

	message := messages[0]
	handle := message.ReceiptHandle
	body := message.Body

	return *body, handle
}

func (s *CodeRunner) Delete(handle *string) {
	s.queue.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(CODE_RUNNER),
		ReceiptHandle: handle,
	})
}

func (s *CodeRunner) Push(code string) {

	s.queue.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(CODE_RUNNER),
		MessageBody: aws.String(code),
	})
}

func S(val string) *dynamodb.AttributeValue {
	//TODO: why doesn't this sdk deal in json?
	if val == "" {
		val = " "
	}
	return &dynamodb.AttributeValue{S: aws.String(val)}
}

func N(val int64) *dynamodb.AttributeValue {
	s := fmt.Sprint(val)
	return &dynamodb.AttributeValue{N: aws.String(s)}
}

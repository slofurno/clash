package datastore

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const queue string = "https://queue.amazonaws.com/027082628651/myqueue"
const CODE_TABLE = "code"
const CLASH_TABLE = "clashes"
const EVENT_TABLE = "event"
const CODE_RUNNER = "https://sqs.us-east-1.amazonaws.com/027082628651/coderunner"
const EVENT_TOPIC = "arn:aws:sns:us-east-1:027082628651:clash_events"

type Code struct {
	Id      string `json:"id"`
	Code    string `json:"code"`
	Runner  string `json:"runner"`
	Problem string `json:"problem"`
	User    string `json:"user"`
	Time    int64  `json:"time"`
	Diff    string `json:"diff"`
	Status  int64  `json:"status"`
}
type Clash struct {
	Id        string `json:"id"`
	Time      int64  `json:"time"`
	Challenge string `json:"challenge"`
}

type DataStore struct {
	Clashes    *ClashStore
	Codes      *CodeStore
	Rooms      *RoomStore
	CodeRunner *CodeRunner
	Events     *EventStore
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

	return &DataStore{
		Clashes:    clashes,
		Codes:      codes,
		CodeRunner: coderunner,
		Events:     events,
		Rooms:      rooms,
	}
}

type Room struct {
	Name string
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

type CodeRunner struct {
	queue *sqs.SQS
}

type Event struct {
	Id      string `json:"id"`
	Subject string `json:"subject"`
	Noun    string `json:"noun"`
	Verb    string `json:"verb"`
	Time    int64  `json:"time"`
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

	fmt.Println(code)

	item := map[string]*dynamodb.AttributeValue{
		"id":      S(code.Id),
		"user":    S(code.User),
		"time":    N(code.Time),
		"code":    S(code.Code),
		"problem": S(code.Problem),
		"runner":  S(code.Runner),
		"diff":    S(code.Diff),
		"status":  N(code.Status),
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
	fmt.Println(X)
	fmt.Println(res.String())

	ret := &Code{}

	code := X["code"]
	problem := X["problem"]
	runner := X["runner"]
	user := X["user"]
	time := X["time"]

	if code == nil || problem == nil || runner == nil || user == nil || time == nil {
		return nil
	}

	ret.Code = *code.S
	ret.Problem = *problem.S
	ret.Runner = *runner.S
	ret.User = *user.S
	ret.Id = id

	i, err := strconv.ParseInt(*time.N, 10, 64)

	if err != nil {
		ret.Time = i
	}

	diff := X["diff"]
	status := X["status"]

	if diff != nil && status != nil {
		i, err := strconv.ParseInt(*status.N, 10, 64)
		if err != nil {
			ret.Status = i
		}

		ret.Diff = *status.S
	}

	return ret

}

func (s *RoomStore) Insert(room *Room) {

	_ = map[string]*dynamodb.AttributeValue{
		"id": S(room.Name),
	}
}

func (s *ClashStore) Insert(clash *Clash) (*dynamodb.PutItemOutput, error) {

	items := map[string]*dynamodb.AttributeValue{
		"id":        S(clash.Id),
		"challenge": S(clash.Challenge),
		"time":      N(clash.Time),
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

	return *message.Body, handle
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

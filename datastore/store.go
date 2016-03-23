package datastore

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const queue string = "https://queue.amazonaws.com/027082628651/myqueue"
const CODE_TABLE = "code"
const CLASH_TABLE = "clashes"
const CODE_RUNNER = "https://sqs.us-east-1.amazonaws.com/027082628651/coderunner"

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
	CodeRunner *CodeRunner
}

func New() *DataStore {

	sess := session.New(&aws.Config{Region: aws.String("us-east-1")})

	ddb := dynamodb.New(sess)
	mysqs := sqs.New(sess)

	clashes := &ClashStore{db: ddb}
	codes := &CodeStore{db: ddb}
	coderunner := &CodeRunner{queue: mysqs}

	return &DataStore{
		Clashes:    clashes,
		Codes:      codes,
		CodeRunner: coderunner,
	}
}

type ClashStore struct {
	db *dynamodb.DynamoDB
}

type CodeStore struct {
	db *dynamodb.DynamoDB
}

type CodeRunner struct {
	queue *sqs.SQS
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

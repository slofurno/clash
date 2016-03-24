package main

import (
	"fmt"
	"github.com/slofurno/front/datastore"
	"github.com/slofurno/front/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {

	store := datastore.New()

	for {

		m, handle := store.CodeRunner.Receive()

		if handle == nil {
			fmt.Println("no messages yet...")
			continue
		}

		code := store.Codes.Get(m)

		if code == nil {
			fmt.Println("missing code")
			store.CodeRunner.Delete(handle)
			continue
		}

		fmt.Println("is this your code?", code.Code)
		fmt.Println("looking for problem:", code.Problem)
		problem := store.Problems.Get(code.Problem)

		if problem == nil {
			fmt.Println("missing problem")
			store.CodeRunner.Delete(handle)
			continue
		}
		fmt.Println("expected output:", problem.Output)
		res, status := diff(code.Code, problem.Output)
		fmt.Println(res)

		code.Diff = res
		code.Status = status

		store.Codes.Insert(code)

		store.Events.Insert(&datastore.Event{
			Id:      utils.Makeid(),
			Subject: code.Id,
			Verb:    "ran",
		})

		store.CodeRunner.Delete(handle)
	}
}

const exs = `
for (var i = 0; i < 10; i++) {
	console.log(i);
}
`

const ex2 = `
for (var i = 0; i < 10; i++) {
	console.log(6);
}
`

func diff(code, expected string) (string, int64) {

	fmt.Println("diff?")
	t1, err := ioutil.TempFile("", "d")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(t1.Name())

	t2, err := ioutil.TempFile("", "d")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(t2.Name())

	if _, err := t1.WriteString(code); err != nil {
		log.Fatal(err)
	}

	if _, err := t2.WriteString(expected); err != nil {
		log.Fatal(err)
	}

	diff := exec.Command("./noderunner.sh", t1.Name(), t2.Name())
	//diff := exec.Command("diff", "<(node "+t1.Name()+")", "<(node "+t2.Name()+")")

	var cout []byte
	var equal int64

	if cout, err = diff.CombinedOutput(); err != nil {
		fmt.Println(err.Error())
		equal = 1
	}

	if err := t1.Close(); err != nil {
		log.Fatal(err)
	}

	if err := t2.Close(); err != nil {
		log.Fatal(err)
	}

	return string(cout), equal
}

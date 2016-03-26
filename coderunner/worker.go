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

		problem := store.Problems.Get(code.Problem)

		if problem == nil {
			fmt.Println("missing problem")
			store.CodeRunner.Delete(handle)
			continue
		}
		fmt.Println("expected output:", problem.Output)

		//FIXME
		run(code)

		fmt.Println("???")
		fmt.Println(code.Output)

		res, status := diff(code.Output, problem.Output)

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

var runners = map[string]string{
	"js": "node",
}

func run(code *datastore.Code) {
	t, err := ioutil.TempFile("", "d")

	defer os.Remove(t.Name())
	t.WriteString(code.Code)

	var cmd string
	var ok bool

	if cmd, ok = runners[code.Runner]; !ok {
		fmt.Println("invalid runner")
		return
	}

	p := exec.Command(cmd, t.Name())
	cout, err := p.CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
	}

	code.Output = string(cout)
	t.Close()
}

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

	diff := exec.Command("diff", "-w", t1.Name(), t2.Name())
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

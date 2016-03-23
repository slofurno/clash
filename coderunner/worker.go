package main

import (
	"fmt"
	"github.com/slofurno/front/datastore"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {

	store := datastore.New()

	for {

		m, handle := store.CodeRunner.Receive()

		if handle != nil {
			code := store.Codes.Get(m)
			fmt.Println(code.Code, code.Problem, code.Runner)
			res, status := diff()

			code.Diff = res
			code.Status = status

			store.Codes.Insert(code)
			store.CodeRunner.Delete(handle)
		} else {
			fmt.Println("no messages yet...")
		}
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

func diff() (string, int64) {

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

	if _, err := t1.WriteString(exs); err != nil {
		log.Fatal(err)
	}

	if _, err := t2.WriteString(exs); err != nil {
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

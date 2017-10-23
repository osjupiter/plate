package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Plate struct {
	Header   string
	Template string
	Record   []map[string]string
}

func main() {
	b, e := ioutil.ReadAll(os.Stdin)
	if e != nil {
		panic(e)
	}
	t := Plate{}
	err := yaml.Unmarshal(b, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Print(t.Header)
	for _, v := range t.Record {
		text := t.Template
		for key, value := range v {
			text = strings.Replace(text, "${"+key+"}", value, -1)
		}
		fmt.Printf("%s", text)
	}

}

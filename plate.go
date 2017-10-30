package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

func strings2Intarfaces(strs []string) []interface{} {
	ret := make([]interface{}, 0)
	for _, s := range strs {
		ret = append(ret, s)
	}
	return ret

}

type GenerateFunc struct {
	Fmt  string
	Args []string
}

func (gf *GenerateFunc) run(record map[string]interface{}) string {
	vars := gf.convert(gf.Args, record)
	return fmt.Sprintf(gf.Fmt, (vars)...)

}
func (gf *GenerateFunc) convert(args []string, record map[string]interface{}) []interface{} {
	ret := make([]interface{}, 0)
	for _, key := range args {
		if value, ok := record[key]; ok {
			ret = append(ret, value)
		} else {
			panic(fmt.Sprint("value of key:%s was not found!"))
		}
	}
	return ret
}

type Plate struct {
	Header   string
	Template string
	Generate map[string]GenerateFunc
	Record   []map[string]interface{}
}

func (p *Plate) generateInputs() []map[string]interface{} {
	ret := make([]map[string]interface{}, 0)
	for _, v := range p.Record {
		rmap := make(map[string]interface{})
		//set default
		rmap["unixtime"] = int32(time.Now().Unix())
		rmap["datetime"] = time.Now().Format("060102150405")

		//copy record
		for key, value := range v {
			rmap[key] = value
		}
		//geranate
		for key, gf := range p.Generate {
			generated := gf.run(rmap)
			rmap[key] = generated
		}
		ret = append(ret, rmap)
	}

	return ret

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

	inputs := t.generateInputs()
	fmt.Print(t.Header)
	for _, v := range inputs {
		text := t.Template
		for key, value := range v {
			repStr := interface2String(value)

			text = strings.Replace(text, "${"+key+"}", repStr, -1)
		}
		fmt.Printf("%s", text)
	}

}

func interface2String(data interface{}) string {

	switch data.(type) {
	case string:
		return data.(string)
	case int:
		return fmt.Sprintf("%d", data.(int))
	case int32:
		return fmt.Sprintf("%d", data.(int32))
	case nil:
		return ""
	default:
		fmt.Fprintf(os.Stderr, "cant toString "+fmt.Sprintf("%#v", data))
	}
	return ""
}

package jsonpick

/**
https://www.codeleading.com/article/17814079501/
*/
import (
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
	"io/ioutil"
	"testing"
)

var jsonPick = gojsonq.New()
var json = ""

func init() {
	s, _ := ioutil.ReadFile("./json.json")
	//fmt.Println(string(s))
	json = string(s)
}

func Test_a(t *testing.T) {
	name := jsonPick.FromString(json).Find("name.first")
	println(name.(string))
}

func Test_b(t *testing.T) {
	avg := jsonPick.FromString(json).From("temperatures").Avg()
	fmt.Printf("Average temperature: %.2f\n", avg)
}

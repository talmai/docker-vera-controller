package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

func main() {

	var myClient = &http.Client{Timeout: 10 * time.Second}

	response, err := myClient.Get("http://192.168.2.7:3480/data_request?id=user_data&output_format=json")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	str2 := string(bytes)

	//ref:  http://michaelheap.com/golang-encodedecode-arbitrary-json/  18/8/15

	// Using some hand crafted JSON. This could come from a file, web service, anything
	// str2 := "{\"foo\":{\"baz\": [1,2,3]}, \"flag\":true, \"list\":[\"one\", 2, true, \"4\", {\"key\":\"value\"}, [1, true]]}"

	var y map[string]interface{}
	json.Unmarshal([]byte(str2), &y)

	// fmt.Printf("%+v\n", y)
	//# => map[foo:map[baz:[1 2 3]] flag:true list:[one 2 true 4 map[key:value] [1 true]]]

	/*
	   As we’re un-marshalling into an interface, we need to inform go what data type
	   each key is before we can perform operations on it. Go provides a the "reflect"
	   package which we can use to process arbitrarily complex data structures:
	*/

	the_list := y["devices"].([]interface{})
	for n, v := range the_list {
		fmt.Printf("\tindex:%d  value:%v  kind:%s  type:%s\n", n, v, reflect.TypeOf(v).Kind(), reflect.TypeOf(v))
	}

	//# =>
	//index:0  value:one  kind:string  type:string
	//index:1  value:2  kind:float64  type:float64
	//index:2  value:true  kind:bool  type:bool
	//index:3  value:4  kind:string  type:string
	//index:4  value:map[key:value]  kind:map  type:map[string]interface {}
	//index:5  value:[1 true]  kind:slice  type:[]interface {}
}

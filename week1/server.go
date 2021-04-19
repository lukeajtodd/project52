package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Country struct {
	country_id string
}

type Result struct {
	name    string
	age     float64
	gender  string
	country Country
}

func main() {
	http.HandleFunc("/predict", predict)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func endpointHandler(endpoint string, chn chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	fmt.Println(endpoint)

	response, err := http.Get(fmt.Sprint("https://api.", endpoint, ".io?name=michael"))
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	chn <- string(data)
}

func predict(w http.ResponseWriter, r *http.Request) {
	chn := make(chan string)
	var wg sync.WaitGroup
	var result Result

	endpoints := []string{
		"nationalize",
		"genderize",
		"agify",
	}

	for _, endpoint := range endpoints {
		go endpointHandler(endpoint, chn, &wg)
	}

	wg.Wait()

	//
	// This below has lots of issues, not sure why/how to pull
	// the values out of the unmarshalled JSON data.
	//
	for v := range chn {
		// fmt.Println(v)
		var holding map[string]interface{}
		json.Unmarshal([]byte(v), &holding)

		if holding["age"] != nil {
			fmt.Println(holding["age"])
			result.age = holding["age"].(float64)
		}

		if holding["gender"] != nil {
			result.gender = holding["gender"].(string)
		}

		// if holding["country"] != nil {
		// 	result.country = holding["country"].(Country)
		// }
	}

	fmt.Println(result)

	close(chn)
}

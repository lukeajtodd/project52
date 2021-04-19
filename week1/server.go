package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("Error during connection upgrade: ", err)
		return
	}

	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error during message read: ", err)
			break
		}

		log.Printf("Received: %s", message)

		err = conn.WriteMessage(messageType, message)

		if err != nil {
			log.Println("Error during message write: ", err)
			break
		}
	}
}

func main() {
	// Set routing rules
	http.HandleFunc("/", home)
	http.HandleFunc("/predict", predict)

	//Use the default DefaultServeMux.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func predictAge() {
	response, err := http.Get("https://api.agify.io?name=michael")
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	data, _ := ioutil.ReadAll(response.Body)
	
	fmt.Println(string(data)))
}

// func predictGender() {

// }

// func predictNationality() {

// }

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Index")
}

func predict(w http.ResponseWriter, r *http.Request) {
	socketHandler(w, r)
	// go predictAge()
	// go predictGender()
	// go predictNationality()
}

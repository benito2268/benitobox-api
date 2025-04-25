package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

var MESSAGES_FILE = "./messages.json"
var USERS_FILE = "./users.json"

type Message struct {
	Author   string `json:"author"`
	Datetime int    `json:"datetime"`
	Content  string `json:"content"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func messages(w http.ResponseWriter, req *http.Request) {
	// ensure this is a GET request
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	file, err := os.Open(MESSAGES_FILE)
	if err != nil {
		http.Error(w, "error opening file: "+MESSAGES_FILE, http.StatusInternalServerError)
		return
	}

	defer file.Close()

	data, dataErr := io.ReadAll(file)
	if dataErr != nil {
		http.Error(w, "error reading file: "+MESSAGES_FILE, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func register(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//read the new user from the request body
	var newUser User
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&newUser); err != nil {
		http.Error(w, "recieved ill formatted JSON string", http.StatusBadRequest)
	}

	//read in the existing users
	var users []User
	oldUsers, err := os.ReadFile(USERS_FILE)
	if err != nil {
		http.Error(w, "error reading users", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(oldUsers, &users)
	if err != nil {
		http.Error(w, "error reading users", http.StatusInternalServerError)
		return
	}

	// check if the user already exists
	for _, user := range users {
		if user.Username == newUser.Username {
			http.Error(w, "user already exists", http.StatusBadRequest)
		}
	}

	//hash the password

}

func post(w http.ResponseWriter, req *http.Request) {
	// ensure this is a POST request
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msgs []Message

	//read the messages and deserialize into Message structs
	filedata, readerr := os.ReadFile(MESSAGES_FILE)
	if readerr != nil {
		http.Error(w, "error reading messages", http.StatusInternalServerError)
		return
	}

	jsonerr := json.Unmarshal(filedata, &msgs)
	if jsonerr != nil {
		http.Error(w, "error reading messages", http.StatusInternalServerError)
		return
	}

	//read the new message out of the request body
	var newMsg Message
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&newMsg); err != nil {
		http.Error(w, "recieved ill formatted JSON string", http.StatusBadRequest)
		return
	}

	msgs = append(msgs, newMsg)

	//re-write the file
	toWrite, err := json.MarshalIndent(msgs, "", "\t")
	if err != nil {
		http.Error(w, "error writing message", http.StatusInternalServerError)
	}

	os.Truncate(MESSAGES_FILE, 0)
	err = os.WriteFile(MESSAGES_FILE, toWrite, 0644)
}

func main() {
	http.HandleFunc("/messages", messages)
	http.HandleFunc("/post", post)
	http.HandleFunc("/register", register)
	http.ListenAndServe(":80", nil)
}

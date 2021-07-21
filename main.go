package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	InvalidNumberOfArguments = errors.New("invalid number of arguments")
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	if _, err := http.Get("http://localhost:8080/loadItems"); err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Print("kvstorage-cli> ")
		inp, _ := reader.ReadString('\n')
		inp = strings.TrimSpace(strings.Replace(inp, "\n", "", -1))

		command := strings.Split(inp, " ")

		switch len(command) {
		case 1, 2, 3:
			switch command[0] {
			case "set":
				url := fmt.Sprintf("http://localhost:8080/items/%s/%s", command[1], command[2])
				resp, err := Request(url, "PUT")
				if err != nil {
					//not sure about this kind of error handling
					fmt.Println(err.Error())
				}
				fmt.Println(resp)
			case "get":
				url := fmt.Sprintf("http://localhost:8080/items/%s", command[1])
				resp, err := Request(url, "GET")
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(resp)
			case "delete":
				url := fmt.Sprintf("http://localhost:8080/items/%s", command[1])
				resp, err := Request(url, "DELETE")
				if err != nil {
					//not sure about this kind of error handling
					fmt.Println(err.Error())
				}
				fmt.Println(resp)
			case "exit":
				//TODO: handle request properly
				if _, err := http.Get("http://localhost:8080/saveItems"); err != nil {
					fmt.Println(err)
				}
				os.Exit(0)
			default:
				fmt.Println("Invalid command")
			}
		default:
			fmt.Println(InvalidNumberOfArguments)
		}
	}
}

func Request(url, method string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK && (method == "PUT" || method == "DELETE") {
		return "OK", nil
	}

	body, err := io.ReadAll(resp.Body)
	m := make(map[string]string)
	err = json.Unmarshal(body, &m)

	//if error is present in json, there will be no value
	_, found := m["error"]
	if found {
		return m["error"], nil
	}

	resp.Body.Close()
	return m["value"], nil
}

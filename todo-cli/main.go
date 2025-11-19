package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Task struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("cmd: ")
		// ReadString blocks until the '\n' character (Enter key) is encountered
		input, err := reader.ReadString('\n')
		CheckError(err)

		parts := strings.Split(input, " ")
		cmd := strings.ReplaceAll(parts[0], "\n", "")

		switch cmd {
		case "list":
			if len(parts) == 1 {
				listall()
			} else {
				id := strings.ReplaceAll(parts[1], "\n", "")
				listbyid(id)
			}
		case "create":
			createtask()
		case "update":
			id := strings.ReplaceAll(parts[1], "\n", "")
			updatetask(id)
		case "delete":
			id := strings.ReplaceAll(parts[1], "\n", "")
			deletetask(id)
		case "exit":
			return
		default:
			fmt.Println("Invalid Command")
		}
	}
}

func deletetask(id string) {
	url := "http://localhost:8080/tasks/" + id
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	CheckError(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	CheckError(err)

	if resp.StatusCode == 200 {
		fmt.Println("Deleted Record:", id)
	} else {
		fmt.Println("No record found with id:", id)
	}
	resp.Body.Close()
}

func updatetask(id string) {
	var task Task
	var priority int
	fmt.Print("Usage: <name>, <description>, <priority>:")
	reader := bufio.NewReader(os.Stdin)

	// Read a line of input until a newline character is encountered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	parts := strings.SplitN(input, ",", 3)
	if len(parts) != 3 {
		fmt.Println("Usage: <name>, <description>, <priority>")
		return
	} else {
		response, err := http.Get("http://localhost:8080/tasks/" + id)
		CheckError(err)
		responseData, err := io.ReadAll(response.Body)
		CheckError(err)

		err = json.Unmarshal(responseData, &task)
		CheckError(err)
		response.Body.Close()

		if len(strings.TrimSpace(parts[0])) > 0 {
			task.Name = strings.TrimSpace(parts[0])
		}
		if len(strings.TrimSpace(parts[1])) > 0 {
			task.Description = strings.TrimSpace(parts[1])
		}
		if len(strings.TrimSpace(parts[2])) > 0 {
			priorityStr := strings.TrimSpace(parts[2])
			_, err = fmt.Sscanf(priorityStr, "%d", &priority)
			if err != nil {
				fmt.Println("Invalid priority. Please enter a valid integer.")
				return
			}
			task.Priority = priority
		}

		jsonData, err := json.Marshal(task)
		CheckError(err)
		url := "http://localhost:8080/tasks/" + fmt.Sprintf("%s", id)
		req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(jsonData)))
		CheckError(err)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		CheckError(err)
		fmt.Println("Updated Record:", id)
		fmt.Println("Response Status:", resp.Status)
		resp.Body.Close()
	}
}

func listall() {
	var response *http.Response
	var err error
	var tasks []Task

	response, err = http.Get("http://localhost:8080/tasks")
	CheckError(err)
	responseData, err := io.ReadAll(response.Body)
	CheckError(err)
	err = json.Unmarshal(responseData, &tasks)
	CheckError(err)
	for _, t := range tasks {
		fmt.Printf("ID: %d, Name: %s, Description: %s, Priority %d\n", t.Id, t.Name, t.Description, t.Priority)
	}
	response.Body.Close()
}

func listbyid(id string) {
	var response *http.Response
	var err error
	var task Task

	response, err = http.Get("http://localhost:8080/tasks/" + id)
	CheckError(err)
	responseData, err := io.ReadAll(response.Body)
	CheckError(err)
	err = json.Unmarshal(responseData, &task)
	CheckError(err)
	fmt.Printf("ID: %d, Name: %s, Description: %s, Priority %d\n", task.Id, task.Name, task.Description, task.Priority)
	response.Body.Close()
}

func createtask() {
	var name, description string
	var priority int
	fmt.Print("Usage: <name>, <description>, <priority>:")
	reader := bufio.NewReader(os.Stdin)

	// Read a line of input until a newline character is encountered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	parts := strings.SplitN(input, ",", 3)
	if len(parts) < 3 {
		fmt.Println("Usage: <name>, <description>, <priority>")
		return
	}
	name = strings.TrimSpace(parts[0])
	description = strings.TrimSpace(parts[1])
	priorityStr := strings.TrimSpace(parts[2])
	_, err = fmt.Sscanf(priorityStr, "%d", &priority)
	if err != nil {
		fmt.Println("Invalid priority. Please enter a valid integer.")
		return
	}

	newTask := Task{
		Name:        name,
		Description: description,
		Priority:    priority,
	}

	fmt.Println(newTask)

	jsonData, err := json.Marshal(newTask)
	CheckError(err)

	response, err := http.Post("http://localhost:8080/tasks", "application/json", strings.NewReader(string(jsonData)))
	CheckError(err)

	responseData, err := io.ReadAll(response.Body)
	CheckError(err)

	fmt.Println(string(responseData))
	response.Body.Close()
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

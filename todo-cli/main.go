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
	//var input string
	var response *http.Response
	//var err error
	var tasks []Task
	var task Task

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
				response, err = http.Get("http://localhost:8080/tasks")
				CheckError(err)
				responseData, err := io.ReadAll(response.Body)
				CheckError(err)
				err = json.Unmarshal(responseData, &tasks)
				CheckError(err)
				for _, t := range tasks {
					fmt.Printf("ID: %d, Name: %s, Description: %s, Priority %d\n", t.Id, t.Name, t.Description, t.Priority)
				}
			} else {
				id := strings.ReplaceAll(parts[1], "\n", "")
				response, err = http.Get("http://localhost:8080/tasks/" + id)
				CheckError(err)
				responseData, err := io.ReadAll(response.Body)
				CheckError(err)
				err = json.Unmarshal(responseData, &task)
				CheckError(err)
				fmt.Printf("ID: %d, Name: %s, Description: %s, Priority %d\n", task.Id, task.Name, task.Description, task.Priority)
			}
			response.Body.Close()
		case "create":
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
				continue
			}
			name = strings.TrimSpace(parts[0])
			description = strings.TrimSpace(parts[1])
			priorityStr := strings.TrimSpace(parts[2])
			_, err = fmt.Sscanf(priorityStr, "%d", &priority)
			if err != nil {
				fmt.Println("Invalid priority. Please enter a valid integer.")
				continue
			}

			newTask := Task{
				Name:        name,
				Description: description,
				Priority:    priority,
			}

			fmt.Println(newTask)

			jsonData, err := json.Marshal(newTask)
			CheckError(err)

			response, err = http.Post("http://localhost:8080/tasks", "application/json", strings.NewReader(string(jsonData)))
			CheckError(err)

			responseData, err := io.ReadAll(response.Body)
			CheckError(err)

			fmt.Println(string(responseData))
			response.Body.Close()
		case "update":
			var priority int
			id := strings.ReplaceAll(parts[1], "\n", "")
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
				continue
			} else {
				response, err := http.Get("http://localhost:8080/tasks/" + id)
				CheckError(err)
				responseData, err := io.ReadAll(response.Body)
				CheckError(err)
				fmt.Println(string(responseData))

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
						continue
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
		case "delete":
			id := strings.ReplaceAll(parts[1], "\n", "")
			url := "http://localhost:8080/task/" + id
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
		case "exit":
			return
		default:
			fmt.Println("Invalid Command")
		}
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

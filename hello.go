package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const delay = 1
const httpOk = 200

func main() {

	introduction()

	for {
		showIntroMenu()

		switch getUserInput() {
		case 1:
			checkStatus()
		case 2:
			retreiveLogs()
		case 3:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Unknown command")
			os.Exit(-1)
		}
	}
}

func introduction() {
	name := "J Freitas"
	version := 1.01

	fmt.Println("Hello,", name)
	fmt.Println("Version: ", version)
}

func getUserInput() int {
	var userInput int

	fmt.Scan(&userInput)
	fmt.Println("Your input is", userInput)

	return userInput
}

func showIntroMenu() {
	fmt.Println("1 - Check status")
	fmt.Println("2 - Show logs")
	fmt.Println("3 - Exit")
}

func retreiveLogs() {
	fmt.Println("Retrieving logs...")

	file, err := ioutil.ReadFile("logs.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(file))
	}
}

func checkStatus() {
	fmt.Println("Checking status...")

	sites := showSites()

	for _, site := range sites {
		testSite(site)
		time.Sleep(delay * time.Second)
	}
}

func showSites() []string {
	var sites []string

	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Error to open the file: ", err)
	}

	reader := bufio.NewReader(file)

	for {
		row, fileErr := reader.ReadString('\n')
		row = strings.TrimSpace(row)

		sites = append(sites, row)

		if fileErr == io.EOF {
			break
		}

	}

	file.Close()

	return sites
}

func testSite(site string) {
	httpResp, err := http.Get(site)

	fmt.Println("Testing ", site)

	if err != nil {
		fmt.Println(err)
	} else {
		if &httpResp.StatusCode != nil {
			if httpResp.StatusCode == httpOk {
				fmt.Println("Status: ", httpResp.StatusCode, " - Successfuly loaded")
			} else {
				fmt.Println("Status: ", httpResp.StatusCode, " We are unable to load this website")
			}
			registerLog(site, httpResp.StatusCode, httpResp.Status)
		} else {
			fmt.Println("Something went wrong")
		}
	}
}

func registerLog(site string, statusCode int, status string) {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	} else {
		file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - ")
		file.WriteString(site + " - status - " + strconv.FormatInt(int64(statusCode), 10) + " - " + status + "\n")
	}

	file.Close()
}

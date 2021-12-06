package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

const regex = `<.*?>`

type Question struct {
	QuestionNo int    `json:"questionNo"`
	Question   string `json:"question"`
	Hint       string `json:"hint"`
}

type Data struct {
	Questions []Question `json:"questions"`
}

type Room struct {
	TotalData []Data `json:"data"`
}

func strip_html_regex(s string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllLiteralString(s, "")
}

func write_to_file(room Room, room_name string) {

	file, err := os.OpenFile(room_name+".md", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	list := "-"
	title := "# " + room_name + "\n\n"
	if _, err := file.WriteString(title); err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < len(room.TotalData[0].Questions); i++ {
		temp := room.TotalData[0].Questions[i].Question
		question := strip_html_regex(temp)

		content := "## " + question + "\n\n" + list + "\n\n"
		if _, err := file.WriteString(content); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	room_name := os.Args[1]
	req := "https://tryhackme.com/api/tasks/" + room_name

	resp, err := http.Get(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	room := Room{}

	if err := json.Unmarshal(body, &room); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	if len(room.TotalData) == 0 {
		fmt.Println("This room does not exist")
		os.Exit(1)
	}

	write_to_file(room, room_name)

}

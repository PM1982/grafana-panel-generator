package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	SwitchIP    = []string{"10.32.0.210", "10.34.0.210", "10.36.0.210"}
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Port(id int) (string, string) {
	if id == 1 {
		return "1", "IP1"
	}
	if id == 2 {
		return "128", "IP128"
	}
	if id == 67 {
		return "193", "Management1"
	}
	if id > 2 && id < 67 {
		return fmt.Sprintf("%d", 126+id), fmt.Sprintf("Ethernet%d", id-2)
	}
	return "-1", "FUUUUUUUUUUU"
}

func main() {
	beg, err := template.ParseFiles("begin.tmpl")
	if err != nil {
		fmt.Print(err)
		return
	}

	end, err := template.ParseFiles("end.tmpl")
	if err != nil {
		fmt.Print(err)
		return
	}

	t, err := template.ParseFiles("dashboard.tmpl")
	if err != nil {
		log.Print(err)
		return
	}
	for k, switchip := range SwitchIP {
		var doc bytes.Buffer
		f, err := os.Create(switchip + ".json")
		f.Close()
		f, err = os.OpenFile(switchip+".json", os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Println("Open file: ", err)
			return
		}
		defer f.Close()
		configbeg := map[string]string{
			"SwitchIP": switchip,
		}
		configend := map[string]string{
			"SwitchIP": switchip,
			"UID":      RandStringRunes(8),
		}
		doc.Reset()
		err = beg.Execute(&doc, configbeg)
		if err != nil {
			log.Print("Can't execute ", err)
			return
		}
		f.WriteString(doc.String())

		for i := 1; i <= 67; i++ {
			portidx, portdesc := Port(i)
			X := (i - 1) % 6 * 4
			Y := (i - 1) / 6 * 7
			config := map[string]string{
				"SwitchIP": switchip,
				"PortIdx":  portidx,
				"ifIndex":  "{{ifIndex}}",
				"PortDesc": portdesc,
				"ID":       fmt.Sprint(i),
				"H":        "7",
				"W":        "4",
				"X":        fmt.Sprint(X),
				"Y":        fmt.Sprint(Y),
			}
			doc.Reset()
			err = t.Execute(&doc, config)
			if err != nil {
				log.Print("Can't execute ", err)
				return
			}
			f.WriteString(doc.String())
			if k < len(SwitchIP) && i < 67 {
				_, err := f.WriteString("\n,\n")
				if err != nil {
					log.Print("Can't execute ", err)
					return
				}
			}
		}
		doc.Reset()
		err = end.Execute(&doc, configend)
		if err != nil {
			log.Print("Can't execute ", err)
			return
		}
		f.WriteString(doc.String())
	}
}

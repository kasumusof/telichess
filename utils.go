package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	chess "github.com/kasumusof/chessgo"
	"github.com/kasumusof/telichess/models"
)

func makeURL(bot *Bot, telemethod string) (string, error) {
	method := bot.endpoint + telemethod

	u, err := url.ParseRequestURI(bot.APIURL)
	if err != nil {
		return "", err
	}
	u.Path = method
	urlStr := u.String()
	return urlStr, nil
}
func makeRequest(bot *Bot, method, url string, data url.Values) (*http.Response, error) {
	c := strings.NewReader(data.Encode())
	// fmt.Println("From request:", c)
	r, _ := http.NewRequest(method, url, c)                                    // URL-encoded payload
	r.Header.Add("Authorization", fmt.Sprintf("auth_token=\"%s\"", bot.Token)) // htis line not really needed
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")          // remove this line to delete webhook
	// r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := bot.Client.Do(r)
	if err != nil {
		log.Println()
		return nil, err
	}
	return resp, nil
}

func requestJSON(bot *Bot, method, url string, data string) (*http.Response, error) {
	r, _ := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))        // URL-encoded payload
	r.Header.Add("Authorization", fmt.Sprintf("auth_token=\"%s\"", bot.Token)) // htis line not really needed
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")          // remove this line to delete webhook
	// r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	var v []byte
	_, _ = r.Body.Read(v)
	resp, err := bot.Client.Do(r)
	if err != nil {
		log.Println()
		return nil, err
	}

	return resp, nil
}

func structToMap(inter interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(inter).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {

		f := iVal.Field(i)

		var v string
		var err error
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
			// log.Println("in structToMap String", typ.Field(i).Tag.Get("json"), v)
		case bool:
			v = strconv.FormatBool(f.Bool())
		default:
			c := f.Interface()
			v, err = structToJSON(c)
		}
		if err != nil || v == "null" || v == "" || v == "false" {
			continue
		}
		values.Set(typ.Field(i).Tag.Get("json"), v)
		// fmt.Println("Data to be sent", values)

	}
	return
}

func structToJSON(inter interface{}) (string, error) {
	vv, err := json.Marshal(inter)
	if err != nil {
		log.Println("from structToJson Error", err)
		return "", err
	}
	return string(vv), nil
}

func addUser(chatID int) {
	userID := fmt.Sprintf("%d", chatID)
	dum := &models.User{}
	dum, err = models.GetUser(userID)
	if err != nil {
		log.Println(err)
		dum = &models.User{
			ChatID:  userID,
			FEN:     chess.StartPos,
			History: "",
		}
		_, err = models.AddUser(dum)
		if err != nil {
			log.Panic("error adding user", err)
		}

	}

}

func sliceToString(slice []string) string {
	return strings.Join(slice, " ")
}
func stringToSlice(str string) []string {
	return strings.Split(str, " ")
}

var trimSpaces = func(text string) string {
	text = strings.TrimSpace(text)
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")
	return text
}

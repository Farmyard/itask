package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-zoo/bone"
	"github.com/renstrom/shortuuid"
)

type BaseResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func responseString(code int, message string) string {
	t, _ := json.Marshal(BaseResponse{Code: code, Message: message})
	return string(t)
}

func responseData(code int, data map[string]interface{}, message string) string {
	t, _ := json.Marshal(BaseResponse{Code: code, Data: data, Message: message})
	return string(t)
}

type Controller struct {
	//
}

func NewController() *Controller {
	return new(Controller)
}

func (c *Controller) Ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Fprint(w, responseData(200, map[string]interface{}{"version": "1.0.0"}, "pong"))
}

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()
	key := shortuuid.New()

	deviceToken := r.FormValue("devicetoken")
	if len(deviceToken) <= 0 {
		fmt.Fprint(w, responseString(400, "deviceToken 不能为空"))
		return
	}

	m := NewModel()
	oldKey := r.FormValue("key")
	if err := m.CheckKeyExists(oldKey); err == nil {
		key = oldKey
	} else {
		m.AddDeviceInfo(deviceToken, key)
	}

	log.Println("注册设备成功")
	log.Println("key: ", key)
	log.Println("deviceToken: ", deviceToken)
	fmt.Fprint(w, responseData(200, map[string]interface{}{"key": key}, "注册成功"))
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	key := bone.GetValue(r, "key")
	category := bone.GetValue(r, "category")
	title := bone.GetValue(r, "title")
	body := bone.GetValue(r, "body")

	defer r.Body.Close()

	m := NewModel()
	device, _ := m.GetDeviceInfoByKey(key)
	deviceToken := device.Token

	r.ParseForm()

	if len(title) <= 0 && len(body) <= 0 {
		//url中不包含 title body，则从Form里取
		for key, value := range r.Form {
			if strings.ToLower(key) == "title" {
				title = value[0]
			} else if strings.ToLower(key) == "body" {
				body = value[0]
			}
		}

	}

	if len(body) <= 0 {
		body = "无推送文字内容"
	}

	params := make(map[string]interface{})
	for key, value := range r.Form {
		params[strings.ToLower(key)] = value[0]
	}

	log.Println(" ========================== ")
	log.Println("key: ", key)
	log.Println("category: ", category)
	log.Println("title: ", title)
	log.Println("body: ", body)
	log.Println("params: ", params)
	log.Println(" ========================== ")

	p := NewPush()
	err := p.PostPush(category, title, body, deviceToken, params)
	if err != nil {
		fmt.Fprint(w, responseString(400, err.Error()))
	} else {
		fmt.Fprint(w, responseString(200, ""))
	}
}

var Stop bool = false

func (c *Controller) Stop(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	Stop = true
	fmt.Fprint(w, responseData(200, map[string]interface{}{"version": "1.0.0"}, "stop"))
}

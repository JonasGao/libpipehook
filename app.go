package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type libHookHandler struct {
	roboUrl string
	key     string
}

type userModel struct {
	Name string
}

type attributesModel struct {
	Id             int    `json:"id"`
	Ref            string `json:"ref"`
	Status         string `json:"status"`
	DetailedStatus string `json:"detailed_status"`
	CreatedAt      string `json:"created_at"`
	FinishedAt     string `json:"finished_at"`
	Duration       int    `json:"duration"`
}

type projectModel struct {
	Name   string `json:"name"`
	WebUrl string `json:"web_url"`
}

type authorModel struct {
	Name  string
	Email string
}

type commitModel struct {
	Message string
	Author  authorModel
}

type buildsModel struct {
	Stage  string
	Status string
}

type hookModel struct {
	ObjectKind       string          `json:"object_kind"`
	ObjectAttributes attributesModel `json:"object_attributes"`
	User             userModel       `json:"user"`
	Project          projectModel    `json:"project"`
	Commit           commitModel     `json:"commit"`
	Builds           []buildsModel   `json:"builds"`
}

func (h *libHookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r := recover()
		var msg string
		if r != nil {
			err := r.(error)
			msg = err.Error()
			http.Error(w, msg, 500)
		}
	}()
	if r.Body == nil {
		w.WriteHeader(400)
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", bytes)
	var hook hookModel
	err = json.Unmarshal(bytes, &hook)
	if err != nil {
		panic(err)
	}
	h.convertAndSend(hook)
	w.WriteHeader(200)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (h *libHookHandler) getSign() (int64, string) {
	//Long timestamp = System.currentTimeMillis();
	stamp := makeTimestamp()
	//String stringToSign = timestamp + "\n" + secret;
	stringToSign := fmt.Sprintf("%d\n%s", stamp, h.key)
	//Mac mac = Mac.getInstance("HmacSHA256");
	//mac.init(new SecretKeySpec(secret.getBytes("UTF-8"), "HmacSHA256"));
	//byte[] signData = mac.doFinal(stringToSign.getBytes("UTF-8"));
	//return URLEncoder.encode(new String(Base64.encodeBase64(signData)),"UTF-8");
	m := hmac.New(sha256.New, []byte(h.key))
	_, err := io.WriteString(m, stringToSign)
	if err != nil {
		panic(err)
	}
	str := base64.StdEncoding.EncodeToString(m.Sum(nil))
	return stamp, url.QueryEscape(str)
}

func (h *libHookHandler) sendMsg(msg string) {
	stamp, sign := h.getSign()
	u := fmt.Sprintf("%s&timestamp=%d&sign=%s", h.roboUrl, stamp, sign)
	fmt.Println("will send", msg)
	resp, err := http.Post(u, "application/json", strings.NewReader(msg))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("send response [%s]%s\n", resp.Status, body)
}

func (h *libHookHandler) convertAndSend(hook hookModel) {
	h.sendMsg(getMsgBody(hook))
}

func main() {
	roboUrl := os.Args[1]
	key := os.Args[2]
	if roboUrl == "" {
		fmt.Println(`Please provide dingtalk robot "send" url`)
		return
	}
	fmt.Println("Startup!")
	err := http.ListenAndServe(":12345", &libHookHandler{
		roboUrl: roboUrl,
		key:     key,
	})
	if err != nil {
		fmt.Println("Startup fail!", err)
	}
}

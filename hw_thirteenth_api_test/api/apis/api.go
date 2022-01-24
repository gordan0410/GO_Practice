package apis

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
)

// request
type req_api struct {
	Id string
}

type req_v1 struct {
	Number string
}

type req_v2 struct {
	S string
}

// respond api
type res_api[T res_internal_data] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// respond v1
type res_v1[T res_internal_data] struct {
	Status string `json:"status"`
	Text   string `json:"text"`
	Result T      `json:"result"`
}

// respong v2
type res_v2[T res_internal_data] struct {
	Err_no  int    `json:"err_no"`
	Err_txt string `json:"err_txt"`
	Result  T      `json:"result"`
}

type res_msg[T res_internal_data] interface {
	res_api[T] | res_v1[T] | res_v2[T]
}

type req_body interface {
	req_api | req_v1 | req_v2 | map[string]string
}

// 小包
type res_api_1_data struct {
	Name string
}

type res_api_2_data struct {
	Book string
	Num  int
}

type res_v1_person_data struct {
	Name string
	Age  int
	Job  string
}

type res_v1_book_data struct {
	Name     string
	Position string
}

type res_v2_teacher_data struct {
	Class string
}

type res_v2_student_data struct {
	Id   string
	Name string
}

type res_internal_data interface {
	res_api_1_data | res_api_2_data | res_v1_person_data | res_v1_book_data | res_v2_student_data | res_v2_teacher_data
}

func Api_get(id string) {
	query := url.Values{}
	query.Set(conf.API.Path1.ID, id)
	msg, err := api[res_api[res_api_1_data], res_api_1_data, map[string]string](http.MethodGet, conf.API.Path1.Path, nil, query, conf.API.Key, conf.API.Value, conf.API.Port)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", msg)
}

func Api_post(id string) {
	req := req_api{
		Id: id,
	}
	msg, err := api[res_api[res_api_2_data], res_api_2_data](http.MethodPost, conf.API.Path2.Path, req, nil, conf.API.Key, conf.API.Value, conf.API.Port)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", msg)
}

func V1_get(id string) {
	query := url.Values{}
	query.Set(conf.V1.Path1.ID, id)
	msg, err := api[res_v1[res_v1_person_data], res_v1_person_data, map[string]string](http.MethodGet, conf.V1.Path1.Path, nil, query, conf.V1.Key, conf.V1.Value, conf.V1.Port)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", msg)
}

func V1_post(id string) {
	req := req_v1{
		Number: id,
	}
	msg, err := api[res_v1[res_v1_book_data], res_v1_book_data](http.MethodPost, conf.V1.Path2.Path, req, nil, conf.V1.Key, conf.V1.Value, conf.V1.Port)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", msg)
}

func V2_get(id string) {
	query := url.Values{}
	query.Set(conf.V2.Path1.ID, id)
	msg, err := api[res_v2[res_v2_teacher_data], res_v2_teacher_data, map[string]string](http.MethodGet, conf.V2.Path1.Path, nil, query, conf.V2.Key, conf.V2.Value, conf.V2.Port)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", msg)
}

func V2_post(id string) {
	req := req_v2{
		S: id,
	}
	msg, err := api[res_v2[res_v2_student_data], res_v2_student_data](http.MethodPost, conf.V2.Path2.Path, req, nil, conf.V2.Key, conf.V2.Value, conf.V2.Port)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", msg)
}

func api[T2 res_msg[T3], T3 res_internal_data, T1 req_body](method, path string, body T1, query url.Values, key, value, port string) (T2, error) {
	// build basic link
	link := url.URL{
		Scheme:   conf.Scheme,
		Host:     net.JoinHostPort(conf.Host, port),
		Path:     path,
		RawQuery: query.Encode(),
	}

	// build new request
	var req *http.Request
	var err error
	var res_data T2
	if method == http.MethodGet {
		// build request
		req, err = http.NewRequest(method, link.String(), nil)
		if err != nil {
			return res_data, err
		}
	} else if method == http.MethodPost {
		// body to json
		json_data, err := json.Marshal(body)
		if err != nil {
			return res_data, err
		}

		// make json as a reader(payload)
		payload := bytes.NewBuffer(json_data)

		// build request
		req, err = http.NewRequest(method, link.String(), payload)
		if err != nil {
			return res_data, err
		}
	}

	req.Header.Add(key, value)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return res_data, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&res_data)
	if err != nil {
		return res_data, err
	}
	return res_data, err
}

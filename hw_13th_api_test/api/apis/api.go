package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

// requestx3, api, v1, v2
type req_api struct {
	Id string
}

type req_v1 struct {
	Number string
}

type req_v2 struct {
	S string
}

// 也可以都改成指標不使用empty_body, 但若碰到不能使用指標形式的項目時, 這個方法會比較實用
type empty_body struct {
}

// request generic
type req_body interface {
	req_api | req_v1 | req_v2 | *empty_body
}

// respond msg
// respond api
type res_api[T Res_data] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// respond v1
type res_v1[T Res_data] struct {
	Status string `json:"status"`
	Text   string `json:"text"`
	Result T      `json:"result"`
}

// respond v2
type res_v2[T Res_data] struct {
	Err_no  int    `json:"err_no"`
	Err_txt string `json:"err_txt"`
	Result  T      `json:"result"`
}

// respond generic
type res_msg[T Res_data] interface {
	*res_api[T] | *res_v1[T] | *res_v2[T]
}

// data
type Res_api_1_data struct {
	Name string
}

type Res_api_2_data struct {
	Book string
	Num  int
}

type Res_v1_person_data struct {
	Name string
	Age  int
	Job  string
}

type Res_v1_book_data struct {
	Name     string
	Position string
}

type Res_v2_teacher_data struct {
	Class string
}

type Res_v2_student_data struct {
	Id   string
	Name string
}

type Res_data interface {
	Res_api_1_data | Res_api_2_data | Res_v1_person_data | Res_v1_book_data | Res_v2_student_data | Res_v2_teacher_data
}

// 資料使用
func Me(ID string) {
	api_get, err := Api_get(ID)
	if err != nil {
		fmt.Println(err)
	}
	api_post, err := Api_post(ID)
	if err != nil {
		fmt.Println(err)
	}
	v1_get, err := V1_get(ID)
	if err != nil {
		fmt.Println(err)
	}
	v1_post, err := V1_post(ID)
	if err != nil {
		fmt.Println(err)
	}
	v2_get, err := V2_get(ID)
	if err != nil {
		fmt.Println(err)
	}
	v2_post, err := V2_post(ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n", *api_get, *api_post, *v1_get, *v1_post, *v2_get, *v2_post)

}

// 因為不同method回傳檔案格式都不一樣，所以不同method使用不同func，回傳的data struct才能唯一；
// 若集中單個func處理則需要再使用func前預先指定回傳的struct格式，邏輯上本末倒至
// api 最終是為了取得data or result
func Api_get(id string) (*Res_api_1_data, error) {
	// 設定url parameter
	query := url.Values{}
	query.Set(conf.API.Path1.ID, id)
	// 給中間層
	res_msg, err := ident_api[Res_api_1_data, *empty_body](http.MethodGet, "/api/1", nil, query)
	return res_msg, err
}

func Api_post(id string) (*Res_api_2_data, error) {
	// 設定request body
	req := req_api{
		Id: id,
	}
	// 給中間層
	res_msg, err := ident_api[Res_api_2_data](http.MethodPost, "/api/2", req, nil)
	return res_msg, err
}

// V1
func V1_get(id string) (*Res_v1_person_data, error) {
	query := url.Values{}
	query.Set(conf.V1.Path1.ID, id)
	res_msg, err := ident_v1[Res_v1_person_data, *empty_body](http.MethodGet, "/v1/person", nil, query)
	return res_msg, err
}

func V1_post(id string) (*Res_v1_book_data, error) {
	req := req_v1{
		Number: id,
	}
	res_msg, err := ident_v1[Res_v1_book_data](http.MethodPost, "/v1/book", req, nil)
	return res_msg, err
}

// V2 request型態寫死
func V2_get(id string) (*Res_v2_teacher_data, error) {
	// 設定url parameter
	query := url.Values{}
	query.Set(conf.V2.Path1.ID, id)
	// 給中間層
	res_msg, err := ident_v2[Res_v2_teacher_data](http.MethodGet, "/v2/teacher", &req_v2{}, query)
	return res_msg, err
}

func V2_post(id string) (*Res_v2_student_data, error) {
	req := req_v2{
		S: id,
	}
	res_msg, err := ident_v2[Res_v2_student_data](http.MethodPost, "/v2/student", &req, nil)
	return res_msg, err
}

// 中間層，處理key等config參數，並解析及判斷回傳資訊，正確則給予result or data，失敗則給予error msg
func ident_api[T3 Res_data, T1 req_body](method, path string, body T1, query url.Values) (*T3, error) {
	var key string
	var value string
	var port string

	// 設定key與port
	key = "key"
	value = conf.API.Value
	port = conf.API.Port

	// 接收底層資料，並根據msg回傳資料
	msg, err := api[*res_api[T3], T3](method, path, body, query, key, value, port)
	if err != nil {
		return nil, err
	}
	if msg.Code == 1 {
		return &msg.Data, nil
	} else {
		return nil, errors.New(msg.Msg)
	}
}

func ident_v1[T3 Res_data, T1 req_body](method, path string, body T1, query url.Values) (*T3, error) {
	var key string
	var value string
	var port string

	key = "api-key"
	value = conf.V1.Value
	port = conf.V1.Port

	msg, err := api[*res_v1[T3], T3](method, path, body, query, key, value, port)
	if err != nil {
		return nil, err
	}

	if msg.Status == "Y" {
		return &msg.Result, nil
	} else {
		return nil, errors.New(msg.Text)
	}

}

// req_body型態寫死，使用api2，返回*http.respond，改為於中間層解json，可簡化中間層及底層泛型複雜度
func ident_v2[T3 Res_data](method, path string, body *req_v2, query url.Values) (*T3, error) {
	var key string
	var value string
	var port string

	// 設定key與port
	key = "token"
	value = conf.V2.Value
	port = conf.V2.Port

	// 接收底層資料
	res, err := api2(method, path, body, query, key, value, port)
	if err != nil {
		return nil, err
	}

	// 解*http.respond's json to struct
	var msg res_v2[T3]
	err = json.NewDecoder(res.Body).Decode(&msg)
	if err != nil {
		return nil, err
	}

	// 根據msg回傳資料
	if msg.Err_no == 0 {
		return &msg.Result, nil
	} else {
		return nil, errors.New(msg.Err_txt)
	}
}

// 底層，發送請求接受回覆(return json)
func api[T2 res_msg[T3], T3 Res_data, T1 req_body](method, path string, body T1, query url.Values, key, value, port string) (T2, error) {

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
			return nil, err
		}
	} else if method == http.MethodPost {
		// body to json
		json_data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		// make json as a reader(payload)
		payload := bytes.NewBuffer(json_data)

		// build request
		req, err = http.NewRequest(method, link.String(), payload)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add(key, value)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&res_data)
	if err != nil {
		return nil, err
	}

	return res_data, nil
}

// 底層，發送請求接受回覆(return *http.respond)
func api2[T1 req_body](method, path string, body *T1, query url.Values, key, value, port string) (*http.Response, error) {

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
	if method == http.MethodGet {
		// build request
		req, err = http.NewRequest(method, link.String(), nil)
		if err != nil {
			return nil, err
		}
	} else if method == http.MethodPost {
		// body to json
		json_data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		// make json as a reader(payload)
		payload := bytes.NewBuffer(json_data)

		// build request
		req, err = http.NewRequest(method, link.String(), payload)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add(key, value)

	client := &http.Client{}
	return client.Do(req)
}

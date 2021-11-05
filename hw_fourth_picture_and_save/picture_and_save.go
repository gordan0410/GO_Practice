package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"gocv.io/x/gocv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

//received Json from client which data is like {Msg: data} 
type msg struct {
	Msg string
}

//decode_img to base64
func encode_img(file_path string) (img_base64_str string ){
	ff, _ := ioutil.ReadFile(file_path)                     //讀入
	img_base64_str = base64.StdEncoding.EncodeToString(ff)  //文件轉base64
	return img_base64_str	
}

func decode_img(msg string){
	ff, _ := base64.StdEncoding.DecodeString(msg)
	img, _, _ := image.Decode(bytes.NewReader(ff))
	out, err := os.Create("./save_selfie.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = jpeg.Encode(out, img, nil)

	if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}

}

//picture
func camera(device_id int) {
	// 確認catch到攝影機
	webcam, err := gocv.VideoCaptureDevice(device_id)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", device_id)
		return
	}
	defer webcam.Close()

	// 新增Window and Mat
	img := gocv.NewMat()
	defer img.Close()
	for i := 0; i < 1; i++ {
		webcam.Read(&img)
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Erro reading the device : %v\n", device_id)
		}
	}
	gocv.IMWrite("selfie.jpg", img)
}

// take picture and decode, then send back to client
func ws_takeandshow_pic(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	device_id := 0
	camera(device_id)
	f := encode_img("./selfie.jpg")
    

	time.Sleep(time.Second)
	if err = c.WriteJSON(f); err != nil {
		log.Println(err)
	}

	defer c.Close()
}


// take picture and save pic in database
func ws_save_pic(w http.ResponseWriter, r *http.Request){
	c, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	defer c.Close()
	req := msg{}

	err = c.ReadJSON(&req)
	if err != nil {
		fmt.Println("Error reading json.", err)
	}

	r_msg := req.Msg
	pic_base64 := strings.Replace(r_msg, "\"", "", -1)
	decode_img(pic_base64)

}

func start_websocket(){
	http.HandleFunc("/ws_takeandshow_pic", ws_takeandshow_pic)
	http.HandleFunc("/ws_save_pic", ws_save_pic)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServ: ", err)
	}
}



func main() {
	start_websocket()
}

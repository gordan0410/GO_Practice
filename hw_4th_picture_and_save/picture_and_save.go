package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./template")))
	http.HandleFunc("/ws", ws_connect)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("ListenAndServ failed: %v", err)
	}
}

func ws_connect(w http.ResponseWriter, r *http.Request) {
	//判斷請求是否為websocket升級請求。
	if websocket.IsWebSocketUpgrade(r) {
		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			fmt.Println("websocket upgrader.Upgrade failed: ", err)
			return
		}
		
		//設定設備id、執行拍照、取得img.jpg的[]byte
		device_id := 0
		img_byte, err := camera(device_id)
		if err != nil {
			fmt.Println("func camera failed: ", err)
			return
		}

		//轉jpg -> base64並傳送至前端
		img_base64 := encode_img(img_byte)
		err = conn.WriteMessage(websocket.TextMessage, []byte(img_base64))
		if err != nil {
			fmt.Println("conn.WriteMessage failed: ", err)
			return
		}

		// 持續接收與回覆訊息
		for{
		_, c, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("websocket conn.ReadMessage() failed: ", err)
			return
		}
		img_base64_return := string(c)
		err = decode_img_and_save(img_base64_return)
		if err == nil {
			conn.WriteMessage(websocket.TextMessage, []byte("儲存成功"))
		}
		}
	} else {
		fmt.Println("not connected")
		return
	}
}

//picture
func camera(device_id int) ([]byte, error) {
	// 確認catch到攝影機
	webcam, err := gocv.VideoCaptureDevice(device_id)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", err)
		return nil, err
	}

	// 等待鏡頭開啟（暖機）(約300ms抓500ms)
	time.Sleep(500 * time.Millisecond)

	// 新增Mat
	img := gocv.NewMat()

	// 讀取畫面進Mat
	if ok := webcam.Read(&img); !ok {
		fmt.Println("Erro happended in webcam.Read.")
		return nil, err
	}

	// 照片轉為jpeg碼並轉為[]byte形式
	img_jpg, err := gocv.IMEncode(gocv.JPEGFileExt, img)
	if err != nil {
		fmt.Println("gocv.IMEncode failed")
		return nil, err
	}
	img_byte := img_jpg.GetBytes()

	// 關機
	err = img.Close()
	if err != nil {
		fmt.Println("img.Close() failed")
		return nil, err
	}
	err = webcam.Close()
	if err != nil {
		fmt.Println("webcam.Close failed")
		return nil, err
	}
	return img_byte, nil
}

func encode_img(img []byte) (img_base64_str string) { //讀入
	img_base64_str = base64.StdEncoding.EncodeToString(img) //文件轉base64
	return img_base64_str
}

func decode_img_and_save(img_base64 string) error {
	ff, err := base64.StdEncoding.DecodeString(img_base64)
	if err != nil {
		fmt.Println("os.WriteFile faild: ", err)
		return err
	}
	err = os.WriteFile("selfie_save.jpg", ff, 0740)
	if err != nil {
		fmt.Println("os.WriteFile faild: ", err)
		return err
	}
	return nil
}

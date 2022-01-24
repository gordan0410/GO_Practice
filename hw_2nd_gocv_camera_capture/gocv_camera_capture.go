package main

import (
	"log"
	"time"

	"gocv.io/x/gocv"
)

func main() {

	// 本機攝影機代號
	device_id := 0

	// 確認catch到攝影機
	webcam, err := gocv.VideoCaptureDevice(device_id)
	if err != nil {
		log.Fatalf("Error opening video capture device: %v\n", err)
	}
	// 等待鏡頭開啟（暖機）(約300ms抓500ms)
	time.Sleep(500 * time.Millisecond)

	// 新增Mat
	img := gocv.NewMat()

	// 讀取畫面進Mat
	webcam.Read(&img)
	if ok := webcam.Read(&img); !ok {
		log.Fatal("Erro happended in webcam.Read.\n")
	}

	// 照片儲存
	gocv.IMWrite("selfie.jpg", img)

	// 關機
	// classifier.Close()
	img.Close()
	webcam.Close()
}

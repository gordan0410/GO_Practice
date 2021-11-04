package main

import (
	"fmt"
	"time"

	"gocv.io/x/gocv"
)

func close_time(sec int) (window_switch int) {
	time.Sleep(time.Duration(sec) * time.Second)
	window_switch = 0
	return window_switch
}

func main() {
	device_id := 0

	// 確認catch到攝影機
	webcam, err := gocv.VideoCaptureDevice(device_id)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", device_id)
		return
	}

	// 新增Window and Mat
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	//設置開關
	window_switch := 1

	if window_switch == 1 {
		for i := 0; i < 100; i++ {
			webcam.Read(&img)
			if ok := webcam.Read(&img); !ok {
				fmt.Printf("Erro reading the device : %v\n", device_id)
			}
			window.IMShow(img)
			window.WaitKey(1)
		}

		close_time(1)
		gocv.IMWrite("selfie.jpg", img)

	} else if window_switch == 0 {
		webcam.Close()
		window.Close()
		img.Close()
	}
}

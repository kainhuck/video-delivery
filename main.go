package main

import (
	"fmt"
	"log"
	"os"
	audio_maker "video-delivery/audio-maker"
	"video-delivery/clamber"
	"video-delivery/deliverer"
	video_maker "video-delivery/video-maker"
)

func main() {
	base := "data"

	// 1. 抓取图片
	fmt.Println("开始内容爬取")
	clamber := clamber.NewIfeng(base)
	title, articleFile, imageFile, err := clamber.Crawl("https://ishare.ifeng.com/c/s/v002YRDPuR6stZtW3GijNKVqFhJw4el98RshG1A1PZRUfrc__")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("内容爬取完成，标题：%v，文章: %v, 图片: %v\n", title, articleFile, imageFile)

	// 2. 文字转语音
	fmt.Println("开始文字转语音")
	audioMaker := audio_maker.NewXunfei(base, os.Getenv("XUNFEI_APPID"), os.Getenv("XUNFEI_APISECRET"), os.Getenv("XUNFEI_APIKEY"))
	audioFile, err := audioMaker.CovertTextToAudio(articleFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("文字转语音完成，语音: %v\n", audioFile)

	// 3. 合成视频
	fmt.Println("开始合成视频")
	videoMaker := video_maker.NewVideoMaker(base)
	videoFile, err := videoMaker.MergeImageAudio(imageFile, audioFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("合成视频完成，视频: %v\n", videoFile)

	// 4. 投放视频
	fmt.Println("开始视频投放")
	deliverer := deliverer.NewBilibili(base)
	if err := deliverer.Delivery(videoFile, imageFile, title, "-"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("视频投放成功")
}

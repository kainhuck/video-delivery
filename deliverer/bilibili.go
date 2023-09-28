package deliverer

import (
	bilibili_go "github.com/kainhuck/bilibili-go"
	"log"
)

type Bilibili struct {
	client *bilibili_go.Client
	path   string
}

func NewBilibili(path string) Deliverer {
	client := bilibili_go.NewClient(bilibili_go.WithAuthStorage(bilibili_go.NewFileAuthStorage("bilibili.hyk.json")))
	client.LoginWithQrCode()

	return &Bilibili{
		client: client,
		path:   path,
	}
}

func (b *Bilibili) Delivery(videoFile string, cover string, title string, desc string, custom ...interface{}) error {
	if err := b.client.RefreshAuthInfo(); err != nil {
		return err
	}

	// 1. 上传视频
	video, err := b.client.UploadVideoFromDisk(videoFile)
	if err != nil {
		return err
	}
	log.Println("视频上传成功")

	// 2. 上传封面
	cover_, err := b.client.UploadCoverFromDisk(cover)
	if err != nil {
		return err
	}
	log.Println("封面上传成功")

	// 3. 投稿
	result, err := b.client.SubmitVideo(&bilibili_go.SubmitRequest{
		Cover:     cover_.Url,
		Title:     title,
		Copyright: 1,
		TID:       37,
		Tag:       "下饭视频,摸鱼音频,蹲坑视频",
		Desc:      desc,
		Recreate:  -1,
		Videos: []*bilibili_go.SubmitVideo{
			video,
		},
		NoReprint: 1,
		WebOS:     2,
	})
	if err != nil {
		return err
	}
	log.Printf("投稿成功🏅️AV号: %v, BV号: %v\n", result.Aid, result.Bvid)

	return nil
}

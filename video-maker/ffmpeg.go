package video_maker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"video-delivery/utils"
)

type Ffmpeg struct {
	path string
}

func (f *Ffmpeg) MergeImageAudio(imageFile, audioFile string) (videoFile string, err error) {
	videoFile = filepath.Join(f.path, "video", utils.TrimFilename(imageFile)+".mp4")

	if utils.ExistFile(videoFile) {
		fmt.Printf("视频文件：%v，已经存在\n", videoFile)
		return videoFile, nil
	}

	cmd := exec.Command("ffmpeg", "-i", imageFile, "-i", audioFile, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental", videoFile)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(videoFile)
		return "", err
	}

	return videoFile, nil
}

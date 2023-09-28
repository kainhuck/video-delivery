package video_maker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	duration, err := getMediaDuration(audioFile)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", imageFile, "-i", audioFile, "-c:v", "libx264", "-t", fmt.Sprintf("%.2f", duration), "-pix_fmt", "yuv420p", videoFile)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(videoFile)
		return "", err
	}

	return videoFile, nil
}

// 获取媒体文件的时长（以秒为单位）
func getMediaDuration(filename string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filename)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(string(output))
	duration := 0.0
	fmt.Sscanf(durationStr, "%f", &duration)

	return duration, nil
}

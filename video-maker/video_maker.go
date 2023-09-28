package video_maker

/*
	ffmpeg
	语音图片合成视频
*/

type VideoMaker interface {
	// MergeImageAudio 将图片，音频合成视频
	MergeImageAudio(imageFile, audioFile string) (videoFile string, err error)
}

func NewVideoMaker(path string) VideoMaker {
	return &Ffmpeg{
		path: path,
	}
}

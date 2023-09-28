package audio_maker

/*
	tts
	文本转语音
*/

type AudioMaker interface {
	// CovertTextToAudio 文字转语音
	CovertTextToAudio(textFile string) (audioFile string, err error)
}

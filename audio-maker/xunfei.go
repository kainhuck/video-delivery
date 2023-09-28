package audio_maker

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
	"video-delivery/utils"
)

type Xunfei struct {
	path      string
	appId     string
	apiSecret string
	apiKey    string
}

func NewXunfei(path string, appId string, apiSecret string, apiKey string) AudioMaker {
	return &Xunfei{
		path:      path,
		appId:     appId,
		apiSecret: apiSecret,
		apiKey:    apiKey,
	}
}

func (x *Xunfei) CovertTextToAudio(textFile string) (audioFile string, err error) {
	audioFile = filepath.Join(x.path, "audio", utils.TrimFilename(textFile)+".mp3")

	if utils.ExistFile(audioFile) {
		fmt.Printf("音频: %v 已经存在\n", audioFile)
		return audioFile, nil
	}

	content, err := os.ReadFile(textFile)
	if err != nil {
		return "", err
	}

	taskId, err := x.createJobZh(X4LingxiaoqiAssist, string(content))
	if err != nil {
		return "", err
	}

	uri, err := x.getJobResultUrl(taskId)
	if err != nil {
		return "", err
	}

	if err := saveFileFromUrl(uri, audioFile); err != nil {
		return "", err
	}

	return audioFile, nil
}

func (x *Xunfei) signUri(uri string) string {
	u, _ := url.Parse(uri)
	gmt := time.FixedZone("GMT", 0)
	timeFormat := time.Now().In(gmt).Format(time.RFC1123)

	signatureOrigin := "host: api-dx.xf-yun.com\n"
	signatureOrigin += "date: " + timeFormat + "\n"
	signatureOrigin += "POST " + u.Path + " HTTP/1.1"
	signature := hmacSha256Base64(signatureOrigin, x.apiSecret)
	authorizationOrigin := fmt.Sprintf(`api_key="%s", algorithm="hmac-sha256", headers="host date request-line", signature="%s"`, x.apiKey, signature)
	authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))

	return fmt.Sprintf("%s?host=api-dx.xf-yun.com&date=%s&authorization=%s", uri, url.QueryEscape(timeFormat), authorization)
}

func hmacSha256Base64(signatureOrigin, apiSecret string) string {
	apiSecretBytes := []byte(apiSecret)
	signatureOriginBytes := []byte(signatureOrigin)

	// Create an HMAC-SHA256 hasher
	hmacSha256 := hmac.New(sha256.New, apiSecretBytes)

	// Write the signature origin bytes to the hasher
	hmacSha256.Write(signatureOriginBytes)

	// Calculate the HMAC-SHA256 digest
	signatureSha := hmacSha256.Sum(nil)

	// Encode the digest as base64 and convert it to a string
	encodedSignature := base64.StdEncoding.EncodeToString(signatureSha)

	return encodedSignature
}

type VCN string

const (
	X4Pengfei          VCN = "x4_pengfei"           // 男声 较年轻
	X4Yeting           VCN = "x4_yeting"            // 女声 较年轻
	X4Qianxue          VCN = "x4_qianxue"           // 女声 较成熟
	X4Guanshan         VCN = "x4_guanshan"          // 男声 较成熟
	X4LingxiaoqiAssist VCN = "x4_lingxiaoqi_assist" // 女声 较年轻
)

type Header struct {
	AppID  string `json:"app_id"`
	TaskID string `json:"task_id,omitempty"`
}

type Audio struct {
	Encoding   string `json:"encoding"`
	SampleRate int    `json:"sample_rate"`
}

type Pybuf struct {
	Encoding string `json:"encoding"`
	Compress string `json:"compress"`
	Format   string `json:"format"`
}

type Dts struct {
	Vcn      VCN    `json:"vcn"`
	Language string `json:"language"`
	Speed    int    `json:"speed"`
	Volume   int    `json:"volume"`
	Pitch    int    `json:"pitch"`
	Rhy      int    `json:"rhy"`
	Audio    *Audio `json:"audio"`
	Pybuf    *Pybuf `json:"pybuf"`
}

type Parameter struct {
	Dts *Dts `json:"dts"`
}

type Text struct {
	Encoding string `json:"encoding"`
	Compress string `json:"compress"`
	Format   string `json:"format"`
	Text     string `json:"text"`
}

type Payload struct {
	Text *Text `json:"text"`
}

type CreateJobReq struct {
	Header    *Header    `json:"header"`
	Parameter *Parameter `json:"parameter"`
	Payload   *Payload   `json:"payload"`
}

type QueryJobReq struct {
	Header *Header `json:"header"`
}

type Response struct {
	Header struct {
		Code       int    `json:"code"`
		Message    string `json:"message"`
		Sid        string `json:"sid"`
		TaskId     string `json:"task_id"`
		TaskStatus string `json:"task_status"` // 1-任务创建成功 2-任务派发失败 4-结果处理中 5-结果处理完成（包含成功/失败）
	} `json:"header"`
	Payload struct {
		Audio struct {
			Audio      string `json:"audio"`
			BitDepth   string `json:"bit_depth"`
			Channels   string `json:"channels"`
			Encoding   string `json:"encoding"`
			SampleRate string `json:"sample_rate"`
		} `json:"audio"`
		Pybuf struct {
			Encoding string `json:"encoding"`
			Text     string `json:"text"`
		} `json:"pybuf"`
	} `json:"payload"`
}

const (
	CreateJobUrl = "https://api-dx.xf-yun.com/v1/private/dts_create"
	QueryJobUrl  = "https://api-dx.xf-yun.com/v1/private/dts_query"
)

func (x *Xunfei) createJob(req *CreateJobReq) (*Response, error) {
	bts, _ := json.Marshal(req)

	request, err := http.NewRequest(http.MethodPost, x.signUri(CreateJobUrl), bytes.NewBuffer(bts))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bts, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	err = json.Unmarshal(bts, &result)

	return &result, err
}

func (x *Xunfei) createJobZh(vcn VCN, text string) (string, error) {
	result, err := x.createJob(&CreateJobReq{
		Header: &Header{AppID: x.appId},
		Parameter: &Parameter{Dts: &Dts{
			Vcn:      vcn,
			Language: "zh",
			Speed:    50,
			Volume:   50,
			Pitch:    50,
			Rhy:      1,
			Audio: &Audio{
				Encoding:   "lame",
				SampleRate: 16000,
			},
			Pybuf: &Pybuf{
				Encoding: "utf8",
				Compress: "raw",
				Format:   "plain",
			},
		}},
		Payload: &Payload{Text: &Text{
			Encoding: "utf8",
			Compress: "raw",
			Format:   "plain",
			Text:     base64.StdEncoding.EncodeToString([]byte(text)),
		}},
	})

	if err != nil {
		return "", err
	}
	if result.Header.Code != 0 {
		return "", fmt.Errorf(result.Header.Message)
	}

	return result.Header.TaskId, nil
}

func (x *Xunfei) queryJob(req *QueryJobReq) (*Response, error) {
	bts, _ := json.Marshal(req)

	request, err := http.NewRequest(http.MethodPost, x.signUri(QueryJobUrl), bytes.NewBuffer(bts))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bts, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	err = json.Unmarshal(bts, &result)

	return &result, err
}

func (x *Xunfei) getJobResultUrl(taskId string) (string, error) {
	for {
		result, err := x.queryJob(&QueryJobReq{Header: &Header{
			AppID:  x.appId,
			TaskID: taskId,
		}})
		if err != nil {
			return "", err
		}

		if result.Header.Code != 0 {
			return "", fmt.Errorf(result.Header.Message)
		}

		if result.Header.TaskStatus == "2" {
			return "", fmt.Errorf("任务派发失败")
		}

		if result.Header.TaskStatus == "5" {
			uri, err := base64.StdEncoding.DecodeString(result.Payload.Audio.Audio)
			if err != nil {
				return "", err
			}
			return string(uri), nil
		}

		time.Sleep(1 * time.Second)
	}
}

func saveFileFromUrl(uri, filename string) error {
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	return err
}

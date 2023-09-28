# video-delivery

视频营销号工具

支持从指定链接爬取文章，文章转音频，自动合成视频，自动上传B站

## 依赖

1. 文字转语音依赖讯飞接口，需提前购买讯飞长文本语音合成服务，并配置以下环境变量(https://console.xfyun.cn/services/long_text)
   XUNFEI_APPID
   XUNFEI_APISECRET
   XUNFEI_APIKEY
2. 视频合成依赖ffmpeg，需提前安装好
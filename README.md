# video-delivery

视频营销号工具

支持从指定链接爬取文章，文章转音频，自动合成视频，自动上传B站

## 依赖

1. 文字转语音依赖讯飞接口，需提前购买讯飞长文本语音合成服务，并配置以下环境变量(https://console.xfyun.cn/services/long_text)

   XUNFEI_APPID

   XUNFEI_APISECRET
   
   XUNFEI_APIKEY
2. 视频合成依赖ffmpeg，需提前安装好
   
   Mac
   ```bash
   brew install ffmpeg
   ```
   
   Ubuntu
   ```bash
   sudo apt-get install ffmpeg
   ```
   
   Arch
   ```bash
   sudo pacman -S ffmpeg
   ```
   
   Windows
   1. 访问FFmpeg的官方网站(https://ffmpeg.org/)，并下载最新版本的FFmpeg。

   2. 解压下载的文件。

   3. 将解压后的文件夹移动到C盘根目录下。

   4. 将FFmpeg的bin目录添加到系统环境变量中。在Windows搜索栏中输入“环境变量”，并打开“编辑系统环境变量”。在“系统属性”窗口中点击“环境变量”按钮，在“系统变量”中找到“Path”变量，点击“编辑”按钮，在“变量值”末尾添加FFmpeg的bin目录路径，多个路径之间用分号隔开。

   5. 点击“确定”按钮保存更改。

   6. 验证FFmpeg是否已成功安装。在命令提示符或PowerShell中输入以下命令：

      ```
      ffmpeg -version
      ```

      如果看到类似于以下内容的输出，则说明FFmpeg已成功安装：

      ```
      ffmpeg version 4.4-static https://johnvansickle.com/ffmpeg/ Copyright (c) 2000-2021 the FFmpeg developers
      built with gcc 8 (Debian 8.4.0-3)
      configuration: --enable-gpl --enable-version3 --enable-static --disable-debug --disable-ffplay --disable-indev=sndio --disable-outdev=sndio --cc=gcc
      libavutil      56. 70.100 / 56. 70.100
      libavcodec     58.134.100 / 58.134.100
      libavformat    58. 76.100 / 58. 76.100
      libavdevice    58. 13.100 / 58. 13.100
      libavfilter     7.110.100 /  7.110.100
      libswscale      5.  9.100 /  5.  9.100
      libswresample   3.  9.100 /  3.  9.100
      ```
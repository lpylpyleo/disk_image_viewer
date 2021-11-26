# Disk Image Viewer
一个简单显示磁盘上的图片和视频服务，支持缓存（永久），视频206 Partial Content。

A simple web service that can view images on the disk.

类似`http-server`和`http.server`，但是把图片直接显示出来而不是给一个链接，方便看远程机器上的图

It's similar to `http-server` and `http.server`, but it displays images in place instead of providing a link. So it's convinient to view images on remote vm. 

##### 从nodejs改成Go重写了一下，方便使用

## Usage


`go run main.go <image_dir_absolute_path>`


## Screenshot

![](doc/screenshot0.png)

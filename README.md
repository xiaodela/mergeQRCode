# mergeQRCode
使用Go语言实现合并微信和支付宝的二维码
# 使用方法
1. 首先你需要安装[Zbar](http://zbar.sourceforge.net/download.html)这个软件
2. 将微信二维码放在**\qrcode\wechat**中
3. 将支付宝二维码放在**\qrcode\alipay**中
4. 你可以选择直接下载[mergeQRCode.exe]()直接运行，或者编译源代码运行
5. 生成index.html后，将此文件部署到**HTTP服务器**中（可以使用GitHub/GitLab的静态服务器
6. 将服务器地址生成为二维码（可以用用网上的生成二维码网站，我这里推荐一个[QR Code Generator](https://www.qr-code-generator.com/)），用支付宝或微信扫描即可
# 注意事项
- 需要把Zbar的根目录添加到Path中
- Zbar的二进制文件路径设置在pay.go的第14行，如有需要请更改
- 生成的index.html需要部署在**HTTP服务器**中,否则无效
# 实现原理
通过Zbar扫描二维码，然后用JavaScript实现跳转达到扫描二维码的效果
跳转使用检测UA的方式实现检测平台

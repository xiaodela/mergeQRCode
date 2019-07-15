package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const qrcodeAlipayPath = "./qrcode/alipay"
const qrcodeWechatPath = "./qrcode/wechat"

const zbar = "zbarimg.exe"

var paths string
var pathsAlipayQRcode string
var pathsWechatQRcode string

var linkAlipay string
var linkWechat string

var javascriptCodeTemplate = `
<!DOCTYPE html>
<html lang="cn">
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <meta http-equiv="X-UA-Compatible" content="ie=edge" />
        <title>Pay</title>
        <style>
        #code-wechat {
            text-align: center;
            margin: 20% 7%;
        }
        img {
            max-width: 70%;
            max-height: 70%;
        }
        .prompt {
            margin: 10%;
            font-size: 120%;
        }
        </style>
    </head>
    <body>
        <div id="code-wechat" style="display:none">
                <div class="prompt">长按识别以下二维码付款</div>
                <img id="wechat-url" />
        </div>
        <script>
            var setting = {
                wechatUrl: "{{Wechat}}",
                aliUrl: "{{Alipay}}",
                qrcodeApi: "https://www.zhihu.com/qrcode?url="
            };
            if (
                !/^http(s*):\/\//.test(location.href) ||
                /^http(s*):\/\/localhost/.test(location.href)
            ) {
                alert("在本地打开无效，请上传到GitHub/GitLab或静态网站！");
            }
            if (navigator.userAgent.match(/Alipay/i)) {
                window.location.href = setting.aliUrl;
            } else if (navigator.userAgent.match(/MicroMessenger\//i)) {
                document.getElementById("code-wechat").style.display = "block";
                document.getElementById("wechat-url").src =
                    setting.qrcodeApi + urlEncode(setting.wechatUrl);
                document.getElementById("code-wechat").style.display = "block";
            } else {
                alert("请使用支付宝或微信扫描二维码！")
            }
            function urlEncode(String) {
                return encodeURIComponent(String)
                    .replace(/'/g, "%27")
                    .replace(/"/g, "%22");
            }
        </script>
    </body>
</html>
`

//扫描二维码(使用Zbar)
func scanQRcode(path string) (string, error) {
	cmd := exec.Command(zbar, "-q", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "Command execution failed", err
	}
	tmp := strings.TrimSpace(string(out))
	str := strings.Split(tmp, "Code:")

	return str[1], nil

}

//获取指定目录下的所有文件
func walkfunc(path string, info os.FileInfo, err error) error {
	paths = path
	return err
}

//得到支付宝和微信二维码文件文件地址
func getPath() {
	filepath.Walk(qrcodeAlipayPath, walkfunc)
	pathsAlipayQRcode = paths
	filepath.Walk(qrcodeWechatPath, walkfunc)
	pathsWechatQRcode = paths
}

//获得支付宝和微信的二维码链接
func getLink() {
	linkAlipay, _ = scanQRcode(pathsAlipayQRcode)
	linkWechat, _ = scanQRcode(pathsWechatQRcode)
}

func createIndexFile() (string, error) {

	tmpA := strings.Replace(javascriptCodeTemplate, "{{Alipay}}", linkAlipay, 1)
	tmpW := strings.Replace(tmpA, "{{Wechat}}", linkWechat, 1)

	f, err := os.Create("index.html")
	defer f.Close()
	if err != nil {
		return "file creation failed", err
	}

	_, err = f.Write([]byte(tmpW))
	if err != nil {
		return "file write failed", err
	}

	return tmpW, nil
}
func run() (string, error) {
	getPath()
	getLink()
	str, err := createIndexFile()
	return str, err
}

func main() {
	data, err := run()
	if err != nil {
		fmt.Println("生成文件错误，错误为：", err)
		fmt.Println(data)
	}
	fmt.Println("文件index.html创建成功！生成在本程序根目录下！")
	fmt.Println("按回车键退出...")
	fmt.Scanln()
	return
}

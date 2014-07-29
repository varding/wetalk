package maintenance

import (
	"fmt"
	"net/http"
)

func handler(rw http.ResponseWriter, req *http.Request) {
	var respData = "<html><body><h1 style='text-align:center;margin-top:20px;'>Go友团网站升级中。。。明天上午恢复正常。。。谢谢！ :)</h1></body></html>"
	rw.Header().Add("Content-Type", "text/html;charset=utf-8")
	rw.Write([]byte(respData))
}

func main() {
	addr := "golanghome:80"
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
	}
}

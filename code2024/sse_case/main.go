package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type SSE struct {
}

var msg = make(chan string, 1)

func (sse *SSE) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	id := req.URL.Query().Get("id")
	log.Println("id=", id)
	channelsMap[id] = make(chan string)

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// // 指定sse
	//go func() {
	//	for {
	//		select {
	//		case <-req.Context().Done():
	//			fmt.Println("req done...")
	//			delete(channelsMap, id)
	//			return
	//			//case <-time.After(1 * time.Second):
	//			//	log.Println("send...")
	//			//	// 必须格式 event: {}\ndata: {}\n\n
	//			//	fmt.Fprintf(rw, "event: ping\ndata: %d\n\n", time.Now().Unix(), time.Now().Unix())
	//			//	flusher.Flush()
	//		}
	//	}
	//}()
	//for msg := range channelsMap[id] {
	//	log.Println("send ->", msg)
	//	fmt.Fprintf(rw, "event: ping\ndata: %s\n\n", msg, time.Now().Unix())
	//	flusher.Flush()
	//}

	for {
		select {
		case <-req.Context().Done():
			fmt.Println("req done...")
			delete(channelsMap, id)
			return

		case <-time.After(1 * time.Second):
			//log.Println("send...")
			// 必须格式 event: {}\ndata: {}\n\n
			fmt.Fprintf(rw, "event: ping\ndata: %d\n\n", time.Now().Unix(), time.Now().Unix())
			flusher.Flush()

		case rse := <-msg: //:= <-channelsMap[id]:
			log.Println("send ->", rse)
			fmt.Fprintf(rw, "event: ping\ndata: %s\n\n", rse, time.Now().Unix())
			flusher.Flush()

		}
	}
}

func main() {
	http.Handle("/sse", &SSE{})
	http.HandleFunc("/user", handler)
	http.HandleFunc("/send", send)
	http.ListenAndServe(":8080", nil)
}

var channelsMap = map[string]chan string{}

func handler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	htmlContent := fmt.Sprintf(`<html>
<body>
<h1>test</h1>
</body>
<script>
    if (window.EventSource) {
        const source = new EventSource("/sse?id=%s");
        source.onopen = (e) => {
            console.log("链接成功");
            console.log(e);
        };
        source.onmessage = function (e) {
            console.log(e.data);
        };
        source.onerror = function (err) {
            console.log(err);
        };
    }
</script>
</html>`, id)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))

	//file, err := os.Open("client.html") // 替换为你想输出的文件路径
	//if err != nil {
	//	http.Error(w, "File not found", http.StatusNotFound)
	//	return
	//}
	//defer file.Close()
	//io.Copy(w, file)
}

func send(w http.ResponseWriter, r *http.Request) {
	//id := r.URL.Query().Get("id")
	//log.Println("id=", id)
	//if c, ok := channelsMap[id]; ok {
	//	c <- "send"
	//}
	msg <- "end"
}

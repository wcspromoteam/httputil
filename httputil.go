package httputil

import (
	"http"
	"fmt"
	"io"
	"os"
	"time"
	"mime"
	"path"
)

// Returns a map of cookies, with the name mapped to the value of the cookie. Additional metadata (expiration date, etc) is thrown away.
func GetCookies(r *http.Request) map[string]string {
	cookiearray := r.Cookie
	fmt.Println(r.Cookie)
	cookies := make(map[string]string, len(cookiearray))
	// Construct a map of cookie names and values.
	for _, cookie := range cookiearray {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}

// Serves a file, with the name given by path, to the connection conn (although it could be any io.Writer). The info string contains extra info about how the serving went, for use in the log. E.g. "file not modified". Differs from http.ServeFile in that it will *never* list directories.
func ServeFile(conn io.Writer, filePath string, req *http.Request) (info string, err os.Error) {
	file, e := os.Open(filePath)
	if e != nil {
		fmt.Println("Error opening file:", e)
		fmt.Fprintln(conn, "HTTP/1.1 404", e)
		fmt.Fprintln(conn, "Content-Length:", len(e.String()))
		fmt.Fprintln(conn, "")
		fmt.Fprintln(conn, e)
		return
	}
	finfo, e := file.Stat()
	if e != nil {
		fmt.Println("Error stating file:", e)
		fmt.Fprintln(conn, "HTTP/1.1 500", e)
		fmt.Fprintln(conn, "Content-Length: 0")
		fmt.Fprintln(conn, "\n")
		return
	}
	// Stolen from HTTP library
	if t, _ := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); t != nil && finfo.Mtime_ns/1e9 <= t.Seconds() {
		fmt.Fprintln(conn, "HTTP/1.1 304 Not Modified")
		fmt.Fprintln(conn, "")
		//fmt.Println("[File not modified]")
		info = "file not modified"
		return
	}
	ext := path.Ext(filePath)
	ctype := mime.TypeByExtension(ext)
	// No extension, serve as a generic binary
	// TODO: Add check for plain text
	if ctype == "" {
		ctype = "application/octlet-stream"
	}
	fmt.Println(ext, ctype)
	fmt.Println("Serving file", filePath)
	fmt.Fprintln(conn, "HTTP/1.1 200 OK")
	fmt.Fprintln(conn, "Last-Modified:", time.SecondsToUTC(finfo.Mtime_ns/1e9).Format(http.TimeFormat))
	fmt.Fprintln(conn, "Content-Type:", ctype)
	fmt.Fprintln(conn, "Content-Length:", finfo.Size)
	fmt.Fprintln(conn, "")
	io.Copy(conn, file)
	//http.ServeFile(c, r, r.URL.Path[1:])
	return
}

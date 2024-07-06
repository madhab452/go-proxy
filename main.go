package main

import (
	"context"
	_ "embed"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed go-proxy.yaml
var conf string

type ProxyServer struct {
	ServerPort string `yaml:"serverPort"`
	Proxies    []struct {
		Target string `yaml:"target"`
		Path   string `yaml:"path"`
	}
}

func main() {
	ps := ProxyServer{}
	err := yaml.Unmarshal([]byte(conf), &ps)
	if err != nil {
		slog.Error("yaml.Unmarshall() failed", "error", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, conf := range ps.Proxies {
			if strings.HasPrefix(r.URL.Path, conf.Path) {
				slog.Info("request url", "req-url", r.RequestURI)
				dest := conf.Target + strings.Replace(r.RequestURI, conf.Path, "", 1)
				req, err := http.NewRequestWithContext(context.Background(), r.Method, dest, r.Body)
				if err != nil {
					slog.Error("http.NewRequestWithContext():", "error", err)
					return
				}
				slog.Info("sending request to:", "url", dest)
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					slog.Error("clnt.Do()", "error", err)
					return
				}
				if res.Body != nil {
					defer res.Body.Close()
				}
				bytes, err := io.ReadAll(res.Body)
				if err != nil {
					if len(bytes) == 0 {
						return
					}
				}
				w.WriteHeader(res.StatusCode)
				w.Write(bytes)
			}
		}
	})

	slog.Info("server running on port: " + ps.ServerPort)
	http.ListenAndServe(ps.ServerPort, mux)
}

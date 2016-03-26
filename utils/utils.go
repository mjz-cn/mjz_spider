package utils

import (
    "net/http"
    "net/url"
    "os/exec"
    "time"

    "golang.org/x/net/proxy"
)

func ExecuteCmd(cmd string, args []string) (res string, err error){
    path, err := exec.LookPath(cmd)
    if err != nil {
        panic(err) // 这里是核心，不能出错
    }

    cmdObject := exec.Command(path, args...)
    var resBytes []byte
    if resBytes, err = cmdObject.Output(); err != nil {
        return "", err
    }
    return string(resBytes), nil
}

func NewClient(proxyUrl string) (*http.Client, bool) {
    client := &http.Client{
        Timeout: time.Second * 100,
    }
    if len(proxyUrl) > 0 {
        proxyParsedUrl, err := url.Parse(proxyUrl)
        if err != nil {
            return client, false
        }
        var transport *http.Transport
        switch proxyParsedUrl.Scheme {
        case "http":
            transport = &http.Transport{Proxy: http.ProxyURL(proxyParsedUrl)}
        case "socks5":
            dialer, err := proxy.FromURL(proxyParsedUrl, proxy.Direct)
            if err != nil {
                return client, false
            }
            transport = &http.Transport{Dial: dialer.Dial}
        default:
            return client, false
        }
        client.Transport = transport
    }
    return client, true
}
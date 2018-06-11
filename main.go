package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
)

func parseCreds(credsPath string) (string, string, error) {
	file, err := os.Open(credsPath)
	if err != nil {
		return "", "", err
	}
	reader := bufio.NewReader(file)
	key, _ := reader.ReadString('\n')
	secret, _ := reader.ReadString('\n')
	return strings.TrimSpace(key), strings.TrimSpace(secret), err
}

// ./renew-aliyun-cdn-cert -creds aliyun.creds -cert my.nat.com.pem -private nat.com.key -domain my.nat.com
func main() {
	region := flag.String("region", "cn-shenzhen", "Region")
	credsPath := flag.String("creds", "", "File path of Aliyun credentials")
	domain := flag.String("domain", "", "Domain name of CDN")
	certPath := flag.String("cert", "", "File path of certificate")
	pkPath := flag.String("private", "", "File path of private key of certificate")

	flag.Parse()

	key, secret, err := parseCreds(*credsPath)
	if err != nil {
		panic(err)
	}
	cdnCli, err := cdn.NewClientWithAccessKey(*region, key, secret)
	if err != nil {
		panic(err)
	}
	cert, err := ioutil.ReadFile(*certPath)
	if err != nil {
		panic(err)
	}
	pk, err := ioutil.ReadFile(*pkPath)
	if err != nil {
		panic(err)
	}

	t := time.Now()
	certName := "cert-" + t.Format("2006-01-02T15:04:05")

	req := cdn.CreateSetDomainServerCertificateRequest()
	req.DomainName = *domain
	req.ServerCertificateStatus = "on"
	req.CertName = certName
	req.ServerCertificate = string(cert)
	req.PrivateKey = string(pk)

	res, err := cdnCli.SetDomainServerCertificate(req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Signed with name: %v\n", certName)
	fmt.Println(res)
}

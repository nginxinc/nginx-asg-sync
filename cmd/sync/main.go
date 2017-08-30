package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var configFile = flag.String("config_path", "/etc/nginx/aws.yaml", "Path to the config file")
var logFile = flag.String("log_path", "", "Path to the log file. If the file doesn't exist, it will be created")
var version = "0.1-2"

const connTimeoutInSecs = 10

func main() {
	flag.Parse()

	if *logFile != "" {
		logF, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Printf("Couldn't open the log file: %v", err)
			os.Exit(10)
		}
		log.SetOutput(io.MultiWriter(logF, os.Stderr))
	}

	log.Printf("nginx-asg-sync version %s", version)

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Printf("Couldn't open the config file: %v", err)
		os.Exit(10)
	}

	cfg, err := parseConfig(data)
	if err != nil {
		log.Printf("Couldn't parse the config file %v: %v", *configFile, err)
		os.Exit(10)
	}

	nginx, err := NewNginxClient(cfg.UpstreamConfEndpont, cfg.StatusEndpoint, connTimeoutInSecs*time.Second)
	if err != nil {
		log.Printf("Couldn't create NGINX client: %v", err)
		os.Exit(10)
	}

	awsClient := createAWSClient(cfg.Region)

	for _, ups := range cfg.Upstreams {
		if ups.Kind == "http" {
			err = nginx.CheckIfHTTPUpstreamExists(ups.Name)
		} else {
			err = nginx.CheckIfStreamUpstreamExists(ups.Name)
		}
		if err != nil {
			log.Printf("Problem with the NGINX configuration: %v", err)
			os.Exit(10)
		}
		exists, err := awsClient.CheckIfAutoscalingGroupExists(ups.AutoscalingGroup)
		if err != nil {
			log.Printf("Couldn't check if an Auto Scaling group exists: %v", err)
			os.Exit(10)
		} else if !exists {
			log.Printf("Warning: Auto Scaling group %v doesn't exists", ups.AutoscalingGroup)
		}
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)

	for {
		for _, ups := range cfg.Upstreams {
			ips, err := awsClient.GetPrivateIPsOfInstancesOfAutoscalingGroup(ups.AutoscalingGroup)
			if err != nil {
				log.Printf("Couldn't get the IP addresses of instances of the Auto Scaling group %v: %v", ups.AutoscalingGroup, err)
				continue
			}

			var backends []string
			for _, ip := range ips {
				backend := fmt.Sprintf("%v:%v", ip, ups.Port)
				backends = append(backends, backend)
			}

			var added, removed []string

			if ups.Kind == "http" {
				added, removed, err = nginx.UpdateHTTPServers(ups.Name, backends)
			} else {
				added, removed, err = nginx.UpdateStreamServers(ups.Name, backends)
			}
			if err != nil {
				log.Printf("Couldn't update servers in NGINX: %v", err)
				continue
			}
			if len(removed) > 0 || len(added) > 0 {
				log.Printf("Upstream: %v has been updated; Added: %v; Removed: %v\n", ups.Name, added, removed)
			}
		}

		select {
		case <-time.After(cfg.SyncIntervalInSeconds * time.Second):
		case <-sigterm:
			log.Println("Terminating...")
			return
		}
	}
}

func createAWSClient(region string) *AWSClient {
	httpClient := &http.Client{Timeout: connTimeoutInSecs * time.Second}
	cfg := &aws.Config{Region: aws.String(region), HTTPClient: httpClient}
	session := session.New(cfg)
	svcAutoscaling := autoscaling.New(session)
	svcEC2 := ec2.New(session)
	return NewAWSClient(svcEC2, svcAutoscaling)
}

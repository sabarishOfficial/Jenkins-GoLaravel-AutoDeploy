package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func main() {

	// permission check variable values get into jenkins environment values
	folders := os.Getenv("folders")
	dir := strings.Split(folders, ", ")
	CheckPermission(dir...)

	// migrations // ask the developer what are migrations run like school superAdmin placement Example.
	workDir := os.Getenv("workDir")
	migrate := os.Getenv("migrate")
	// in case not set or empty in migration the migrate function not working
	if len(migrate) >= 1 {
		cmd := strings.Split(migrate, ", ")
		Migrations(workDir, cmd...)
	}

	zoneID := os.Getenv("zoneID")
	api := os.Getenv("api")
	CacheClear(zoneID, api)
	ID := os.Getenv("CloudFrontID")
	CloudFrontCacheInvalidate(ID)

}

func Migrations(workdir string, cmd ...string) string {
	var output []byte
	for _, migrate := range cmd {
		migrateCmd := exec.Command("bash", "-c", migrate)
		migrateCmd.Dir = workdir
		out, err := migrateCmd.CombinedOutput()

		if migrateCmd.ProcessState.ExitCode() != 0 {
			panic(err)
		} else {
			fmt.Printf("command %s\noutput %s\n", migrate, string(out))

		}
		output = append(output, out...)
	}
	return string(output)

}
func CheckPermission(dir ...string) string {
	var out []byte
	var err error

	for _, dirs := range dir {
		CheckPer := exec.Command("sudo", "stat", "-c", "%U", dirs)
		out, err = CheckPer.CombinedOutput()
		if err != nil {
			panic(err)
		}
		if owner := strings.TrimSpace(string(out)); owner != "apache" {
			fmt.Println("you need change the permission")
			ChangeOwnership(dirs)
		}
	}
	return string(out)

}

func ChangeOwnership(dir ...string) string {
	var out []byte
	var err error
	for _, dirs := range dir {
		changeOwner := exec.Command("sudo", "chown", "-R", "apache:apache", dirs)

		_, err = changeOwner.CombinedOutput()
		if err != nil {
			panic(err)
		}

	}
	return string(out)
}
func CacheClear(zoneId, token string) {
	cloudFlare := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", zoneId)
	fmt.Println(cloudFlare)
	data := []byte(`{"purge_everything": true}`)

	req, reqErr := http.NewRequest("POST", cloudFlare, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("content-Type", "application/json")
	if reqErr != nil {
		panic(reqErr)
	}
	client := http.Client{}

	resp, respErr := client.Do(req)

	if respErr != nil {
		panic(respErr)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))

}
func CloudFrontCacheInvalidate(Id string) {
	cloudFrontConnect := cloudfront.New(session.New())
	params := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(Id),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(
				fmt.Sprintf("goinvali%s", time.Now().Format("2006/01/02,15:04:05"))),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(1),
				Items: []*string{
					aws.String("/*"),
				},
			},
		},
	}
	req, err := cloudFrontConnect.CreateInvalidation(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(req)

	invalidationID := *req.Invalidation.Id
	checkInvalidationStatus(cloudFrontConnect, Id, invalidationID)
}
func checkInvalidationStatus(cloudFrontConnect *cloudfront.CloudFront, ID string, InvalidateID string) {
	Param := &cloudfront.GetInvalidationInput{
		DistributionId: aws.String(ID),
		Id:             aws.String(InvalidateID),
	}
	for {
		req, err := cloudFrontConnect.GetInvalidation(Param)

		if err != nil {
			panic(err)
		}
		status := *req.Invalidation.Status

		if status == "Completed" || status == "Failed" {
			fmt.Println("CloudFront Cache Clear")
			break
		} else {
			fmt.Println("InProgress")
		}
		time.Sleep(5 * time.Second)
	}

}

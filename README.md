
# Jenkins Laravel Deployment Pipeline

# Overview
This Jenkins pipeline script automates the deployment process for Laravel projects. It encompasses Laravel migration, storage and bootstrap permission checks, Jenkins code migration, domain cache clearing, and CDN (AWS CloudFront and Cloudflare) invalidation.

# Installation
Make sure to install the required Go packages by running the following commands:
```
go get github.com/aws/aws-sdk-go/aws
go get github.com/aws/aws-sdk-go/aws/awsutil@v1.50.9
```
Download all the dependencies needed for the project.

# Prerequisites
Before using this tool, ensure you have the following prerequisites:
- Cloudflare API key and zone ID
- AWS access key and secret key for CloudFront invalidation
- AWS CloudFront ID

# Variables Example
Here's an example of the variables that need to be configured before running the tool:

```
workDir = "/usr/share/nginx/html"
folders = "/usr/share/nginx/html/resources/pending-files, /usr/share/nginx/html/vendor, /usr/share/nginx/html/storage, /usr/share/nginx/html/bootstrap"
migrateCommand = "sudo php artisan migrate"
zoneID = "jyb1f3in"
apiKey = "123676"
cloudFrontID = "12390"
```
# Usage
To execute the deployment tool, run the following command:
```
go run main.go
```
This command will trigger the Laravel migration, check permissions for storage and bootstrap folders, migrate Jenkins code, clear domain cache, and invalidate CDN on both AWS CloudFront and Cloudflare.
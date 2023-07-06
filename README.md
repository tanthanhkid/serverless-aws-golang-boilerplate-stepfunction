<!--
title: .'HTTP GET and POST'
description: 'Boilerplate code for Golang with GET and POST example'
framework: v1
platform: AWS
language: Go
priority: 10
authorLink: 'https://github.com/pramonow'
authorName: 'Pramono Winata'
authorAvatar: 'https://avatars0.githubusercontent.com/u/28787057?v=4&s=140'
-->

# Serverless-golang facecollection with AWS Rekognition 
This example is using AWS Request and Response Proxy Model, provided by AWS itself.
If you want to test any changes don't forget to run `make` inside the service directory.
 

Run this to create rekognition collection, if you change the name of the collection, make sure to change in serverless.yml file at line 26
```
aws rekognition create-collection --collection-id "viblo"
```

Run this to create S3 bucket to store yours images, if you change the name of the bucket, make sure to change in serverless.yml file at line 27
```
aws s3api create-bucket \
    --bucket viblo-facecollection \ 
    --region ap-southeast-1
```

# API Specs
- Viết 2 api 
  - index face: add face to collection
  - search face: search by image

### API Index
 
- input:
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "username": {{string}},
        "image": {{string}}, //base64 encoded image
    }
}
```
- output:
```
{
    "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}}
}
```

### API Search
 
- input:
```
{
    "requestId": {{uuid}},
    "requestTime": {{timeRPC3339}},
    "data": {
        "image": {{string}}, //base64 encoded image
    }
}
```

- output:
```
{
    "responseId": {{requestId}},
    "responseTime": {{timeRPC3339}},
    "responseCode": {{string}},
    "responseMessage": {{string}},
    "data": [] //faces matches array
}
```
 
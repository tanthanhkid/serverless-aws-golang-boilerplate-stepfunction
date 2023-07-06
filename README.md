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

# Serverless-golang Stepfunction detect feedback sentiment from customer

This example is using AWS Request and Response Proxy Model, provided by AWS itself.
If you want to test any changes don't forget to run `make` inside the service directory.

# API Specs

- input:

```
{
  "Comment": "dịch vụ ở đây không tốt lắm"
}
```

- output:

```
{
    "executionArn": "arn:aws:states:ap-southeast-1:xxxx:execution:CustomerFeedbackSentiment:xxxxx",
    "startDate": 1.688637505746E9
}
```
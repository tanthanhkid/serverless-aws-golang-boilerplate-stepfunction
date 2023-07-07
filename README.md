

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
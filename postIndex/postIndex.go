package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"

	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var logger *log.Logger

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	RequestId   string      `json:"requestId"`
	RequestTime string      `json:"requestTime"`
	Data        DataRequest `json:"data"`
}

type DataRequest struct {
	UserName string `json:"userName"`
	Image    string `json:"image"`
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseId      string `json:"responseId"`
	ResponseTime    string `json:"responseTime"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

type App struct {
	S3 *s3.S3
}

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	datetime := time.Now().UTC()
	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{
		RequestId: "",
	}

	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 401}, nil
	}

	//verify uuid not null
	if bodyRequest.RequestId == "" {
		return events.APIGatewayProxyResponse{Body: "requestId can not be null", StatusCode: 401}, nil

	}

	logger.SetPrefix("[Request ID:" + bodyRequest.RequestId + "] - ")

	//verify datetime format RFC3339
	parsedTime, err := time.Parse(time.RFC3339, bodyRequest.RequestTime)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error() + "parsedTime: " + parsedTime.GoString(), StatusCode: 401}, nil
	}

	if bodyRequest.Data.Image == "" || bodyRequest.Data.UserName == "" {
		return events.APIGatewayProxyResponse{Body: "User object can not be null", StatusCode: 401}, nil
	}

	// call AWS Rekognition to index face
	output, err := indexFace(bodyRequest.Data.UserName, bodyRequest.Data.Image)

	responseCode := "06"
	if len(output.FaceRecords) > 0 && err == nil {
		responseCode = "00"
	}

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		ResponseId:      uuid.New().String(),
		ResponseTime:    datetime.Format(time.RFC3339),
		ResponseCode:    responseCode,
		ResponseMessage: "face indexed: " + strconv.Itoa(len(output.FaceRecords)),
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	// print response json
	responseJson, err := json.Marshal(bodyResponse)
	if err != nil {
		logger.Fatalln("cannot parse response to json")
	}
	logger.Println("RESPONSE: " + string(responseJson))

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	logger = log.New(os.Stderr, "", log.LstdFlags)
	lambda.Start(Handler)
}

func indexFace(userName string, image string) (*rekognition.IndexFacesOutput, error) {
	collectionId := os.Getenv("COLLECTION_ID")
	facesBucket := os.Getenv("FACES_BUCKET")

	// Load the SDK's configuration from environment and shared config, and
	// create the client with this.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Fatalf("failed to load SDK configuration, %v", err)
	}

	client := rekognition.NewFromConfig(cfg)

	//parse image from base 64 and upload to S3 bucket
	decodedSignature, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		log.Fatalf("decode base64 failed, %v", err)
	}
	r := bytes.NewReader(decodedSignature)

	//create s3 input
	s3ObjectName := userName + ".jpg"

	s3Input := &s3.PutObjectInput{
		Body:   r,
		Bucket: &facesBucket,
		Key:    &s3ObjectName,
	}

	//create new session
	sess, err := createSession()
	if err != nil {
		log.Fatalf("failed to create AWS session, %v", err)
	}
	s3 := s3.New(sess)

	//upload image file to s3
	s3output, err := s3.PutObject(s3Input)

	logger.Printf("S3 output message: %v", s3output)

	if err != nil {
		log.Fatalf("failed to put object, %v", err)
	}

	//get image from S3 bucket and index with rekognition
	input := &rekognition.IndexFacesInput{
		Image: &types.Image{
			S3Object: &types.S3Object{
				Bucket: &facesBucket,
				Name:   &s3ObjectName,
			},
		},
		CollectionId:    &collectionId,
		ExternalImageId: &userName,
	}

	output, err := client.IndexFaces(context.TODO(), input)

	if err != nil {
		logger.Fatalf("err when index image, %v", err)
	}

	return output, err
}

func createSession() (*session.Session, error) {
	// create sesssion step
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	return sess, err
}

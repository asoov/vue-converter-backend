package cloudwatch

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func LogToCloudWatch(logMessage string) {
	sess := session.Must(session.NewSession())

	// Create a CloudWatch Logs client
	cw := cloudwatchlogs.New(sess, &aws.Config{Region: aws.String("eu-central-1")})

	// Define log group and log stream
	logGroup := "VueConverter"
	logStream := "VueConverterBackend"

	// Get the current time in milliseconds
	timestamp := time.Now().UnixMilli()

	// Create log event
	logEvent := cloudwatchlogs.InputLogEvent{
		Message:   aws.String(logMessage),
		Timestamp: aws.Int64(timestamp),
	}

	// Put log event
	_, err := cw.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &logGroup,
		LogStreamName: &logStream,
		LogEvents:     []*cloudwatchlogs.InputLogEvent{&logEvent},
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Successfully logged to CloudWatch")
}

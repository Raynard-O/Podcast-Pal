package library

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"os"
	"sync"
)

var once = &sync.Once{}

var (
	BUCKET  = "podcastpal"
	AWS_KEY = os.Getenv("AWS_KEY")
	AWS_ID  = os.Getenv("AWS_ID")
	region  = "us-east-2"
)

var policy = fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action":["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"],"Sid": ""}]}`, BUCKET)
var Client *minio.Client

type Response struct {
	Url string `json:"url"`
}

func InitClient() {
	var err error
	Client, err = minio.New("s3.amazonaws.com", &minio.Options{
		Creds:  credentials.NewStaticV4(AWS_ID, AWS_KEY, ""),
		Secure: true,
	})
	if err != nil {
		log.Panicln(err)
	}

	if exists, err := Client.BucketExists(context.Background(), BUCKET); !exists || err != nil {

		if err = Client.MakeBucket(context.Background(), BUCKET, minio.MakeBucketOptions{Region: region}); err != nil {
			log.Panicln(err)
		}
		log.Print(AWS_ID)
		Client.SetBucketPolicy(context.Background(), BUCKET, policy)
	}
}

func AwsUpload(b io.Reader, size int64, file_key string) error {
	address := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", BUCKET, region, file_key)
	log.Println(address) // <= file url
	res, err := Client.PutObject(context.Background(), BUCKET, file_key, b, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}

func Trans(url string) transcribeservice.StartTranscriptionJobInput {
	jobIn := transcribeservice.StartTranscriptionJobInput{

		JobExecutionSettings: &transcribeservice.JobExecutionSettings{
			AllowDeferredExecution: aws.Bool(true),
			DataAccessRoleArn:      aws.String("my-arn"),
		},
		LanguageCode: aws.String("en-US"),
		MediaFormat:  aws.String("mp3"),
		Media:        &transcribeservice.Media{MediaFileUri: aws.String(url)},
		Settings: &transcribeservice.Settings{
			ChannelIdentification: nil,
			MaxAlternatives:       aws.Int64(2),
			MaxSpeakerLabels:      aws.Int64(2),
			ShowAlternatives:      aws.Bool(true),
			ShowSpeakerLabels:     aws.Bool(true),
		},
		TranscriptionJobName: aws.String("jobName"),
	}
	return jobIn
}

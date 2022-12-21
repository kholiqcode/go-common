package s3

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	common_utils "github.com/kholiqcode/go-common/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Init(baseConfig *common_utils.Config) (*s3.Client, error) {
	creeds := credentials.NewStaticCredentialsProvider(baseConfig.S3.AccessKey, baseConfig.S3.SecretKey, "")

	customEndpointResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               baseConfig.S3.Endpoint,
			SigningRegion:     baseConfig.S3.Region,
			HostnameImmutable: true,
		}, nil
	})

	logMode := aws.ClientLogMode(0)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	no_ssl_verify := &http.Client{Transport: tr}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creeds),
		config.WithRegion(baseConfig.S3.Region),
		config.WithEndpointResolver(customEndpointResolver),
		config.WithClientLogMode(logMode),
		config.WithHTTPClient(no_ssl_verify),
		config.WithRetryer(func() aws.Retryer {
			return aws.NopRetryer{}
		}),
	)

	return s3.NewFromConfig(cfg), err
}

func PreSignClient(client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(client)
}

type S3ClientImpl struct {
	client  S3Client
	preSign S3PreSign
	config  *common_utils.Config
}

func NewS3Client(
	client S3Client,
	preSign S3PreSign,
	config *common_utils.Config,
) S3File {
	return &S3ClientImpl{
		client:  client,
		preSign: preSign,
		config:  config,
	}
}

func (s *S3ClientImpl) UploadPrivateFile(ctx context.Context, file multipart.File, path string) (string, error) {

	_, err := s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.config.S3.PrivateBucketName),
		Key:    aws.String(path),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	log.Println("success uploaded file to S3: ", path)

	preSignUrl, err := s.GetPreSignUrl(ctx, path)

	if err != nil {
		return "", err
	}

	return preSignUrl, nil
}

func (s *S3ClientImpl) UploadPublicFile(ctx context.Context, file multipart.File, path string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.config.S3.PublicBucketName),
		Key:         aws.String(path),
		Body:        file,
		ContentType: aws.String(common_utils.GetExt(path)),
	})

	if err != nil {
		return "", err
	}

	return s.BuildPublicUrl(path), nil
}

func (s *S3ClientImpl) DeleteFile(ctx context.Context, bucketName string, path string) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
	})

	return err
}

func (s *S3ClientImpl) GetPreSignUrl(ctx context.Context, path string) (string, error) {
	params := &s3.GetObjectInput{
		Bucket:              aws.String(s.config.S3.PrivateBucketName),
		Key:                 aws.String(path),
		ResponseContentType: aws.String(common_utils.GetExt(path)),
	}

	resp, err := s.preSign.PresignGetObject(ctx, params, func(po *s3.PresignOptions) {
		po.Expires = s.config.S3.PreSignUrlDuration * time.Second
	})

	if err != nil {
		return "", err
	}

	return resp.URL, nil
}

func (s *S3ClientImpl) BuildPublicUrl(path string) string {

	if s.config.S3.PublicUrl == "" {
		return s.config.S3.PublicUrl + path
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		s.config.S3.PublicBucketName,
		s.config.S3.Region,
		path,
	)
}

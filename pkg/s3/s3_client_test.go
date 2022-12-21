package s3

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"testing"

	"github.com/kholiqcode/go-common/helper"
	mock_s3 "github.com/kholiqcode/go-common/pkg/s3/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/stretchr/testify/assert"
)

const (
	PATH                      = "users/userId/test-image.png"
	IMAGE_NAME                = "test-image.png"
	IMAGE_PRE_SIGNED_URL      = "https://amazon.s3.com/users/userId/test-image.png?pre-signed"
	IMAGE_URL                 = "https://amazon.s3.com/users/userId/test-image.png"
	FAILED_PUT_OBJECT         = "failed when put object"
	FAILED_GET_PRE_SIGNED_URL = "failed when get pre sign url"
	FAILED_DELETE_IMAGE       = "failed when deleting image"
)

func loadConfig() *common_utils.Config {
	cfg, err := common_utils.LoadConfig("")
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func initS3Client(ctrl *gomock.Controller, config *common_utils.Config) (S3File, *mock_s3.MockS3Client, *mock_s3.MockS3PreSign) {
	client := mock_s3.NewMockS3Client(ctrl)
	preSignClient := mock_s3.NewMockS3PreSign(ctrl)

	return NewS3Client(client, preSignClient, config), client, preSignClient
}

func createTestFile() multipart.File {
	filesHeader := helper.CreateFilesHeader(1, IMAGE_NAME)
	file, _ := filesHeader[0].Open()
	return file
}

func TestInitS3(t *testing.T) {
	config := loadConfig()
	s3Client, err := Init(config)

	assert.NotNil(t, s3Client)
	assert.Nil(t, err)
}

func TestUploadPrivateFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	config := loadConfig()
	file := createTestFile()
	defer file.Close()

	s3Client, client, preSignClient := initS3Client(ctrl, config)

	t.Run("Success upload private file and get pre sign url", func(t *testing.T) {
		client.EXPECT().PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(config.S3.PrivateBucketName),
			Key:    aws.String(PATH),
			Body:   file,
		}).Return(&s3.PutObjectOutput{}, nil).Times(1)

		preSignClient.EXPECT().PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket:              aws.String(config.S3.PrivateBucketName),
			Key:                 aws.String(PATH),
			ResponseContentType: aws.String(common_utils.GetExt(PATH)),
		}, gomock.Not(nil)).DoAndReturn(func(_ interface{}, _ interface{}, callback func(po *s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
			opt := s3.PresignOptions{}
			callback(&opt)

			return &v4.PresignedHTTPRequest{
				URL: IMAGE_PRE_SIGNED_URL}, nil
		}).Times(1)

		preSignedUrl, err := s3Client.UploadPrivateFile(ctx, file, PATH)

		assert.NoError(t, err)
		assert.Equal(t, IMAGE_PRE_SIGNED_URL, preSignedUrl)
	})

	t.Run("Failed put object", func(t *testing.T) {
		client.EXPECT().PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(config.S3.PrivateBucketName),
			Key:    aws.String(PATH),
			Body:   file,
		}).Return(&s3.PutObjectOutput{}, errors.New(FAILED_PUT_OBJECT)).Times(1)

		preSignClient.EXPECT().PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket:              aws.String(config.S3.PrivateBucketName),
			Key:                 aws.String(PATH),
			ResponseContentType: aws.String(common_utils.GetExt(PATH)),
		}, gomock.Not(nil)).Return(&v4.PresignedHTTPRequest{
			URL: IMAGE_PRE_SIGNED_URL,
		}, nil).Times(0)

		preSignedUrl, err := s3Client.UploadPrivateFile(ctx, file, PATH)

		assert.Equal(t, FAILED_PUT_OBJECT, err.Error())
		assert.Equal(t, "", preSignedUrl)
	})

	t.Run("Failed get pre sign url", func(t *testing.T) {
		client.EXPECT().PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(config.S3.PrivateBucketName),
			Key:    aws.String(PATH),
			Body:   file,
		}).Return(&s3.PutObjectOutput{}, nil).Times(1)

		preSignClient.EXPECT().PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket:              aws.String(config.S3.PrivateBucketName),
			Key:                 aws.String(PATH),
			ResponseContentType: aws.String(common_utils.GetExt(PATH)),
		}, gomock.Not(nil)).Return(&v4.PresignedHTTPRequest{}, errors.New(FAILED_GET_PRE_SIGNED_URL)).
			Times(1)

		preSignedUrl, err := s3Client.UploadPrivateFile(ctx, file, PATH)

		assert.Equal(t, FAILED_GET_PRE_SIGNED_URL, err.Error())
		assert.Equal(t, "", preSignedUrl)
	})
}

func TestUploadPublicFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	config := loadConfig()
	file := createTestFile()
	defer file.Close()

	s3Client, client, _ := initS3Client(ctrl, config)

	t.Run("Success upload public file", func(t *testing.T) {
		client.EXPECT().PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(config.S3.PublicBucketName),
			Key:         aws.String(PATH),
			Body:        file,
			ContentType: aws.String(common_utils.GetExt(PATH)),
		}).Return(&s3.PutObjectOutput{}, nil).Times(1)

		url, err := s3Client.UploadPublicFile(ctx, file, PATH)

		assert.NoError(t, err)
		assert.NotEqual(t, IMAGE_URL, url)
	})
	t.Run("Failed put object", func(t *testing.T) {
		client.EXPECT().PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(config.S3.PublicBucketName),
			Key:         aws.String(PATH),
			Body:        file,
			ContentType: aws.String(common_utils.GetExt(PATH)),
		}).Return(&s3.PutObjectOutput{}, errors.New(FAILED_PUT_OBJECT)).
			Times(1)

		url, err := s3Client.UploadPublicFile(ctx, file, PATH)
		assert.Equal(t, FAILED_PUT_OBJECT, err.Error())
		assert.Equal(t, "", url)
	})
}

func TestDeleteImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	config := loadConfig()
	file := createTestFile()
	defer file.Close()

	s3Client, client, _ := initS3Client(ctrl, config)

	t.Run("Success delete image", func(t *testing.T) {
		client.EXPECT().DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(config.S3.PrivateBucketName),
			Key:    aws.String(PATH),
		}).
			Return(&s3.DeleteObjectOutput{}, nil).Times(1)

		err := s3Client.DeleteFile(ctx, config.S3.PrivateBucketName, PATH)
		assert.NoError(t, err)
	})

	t.Run("Failed delete image", func(t *testing.T) {
		client.EXPECT().DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(config.S3.PrivateBucketName),
			Key:    aws.String(PATH),
		}).
			Return(&s3.DeleteObjectOutput{}, errors.New(FAILED_DELETE_IMAGE)).Times(1)

		err := s3Client.DeleteFile(ctx, config.S3.PrivateBucketName, PATH)
		assert.Equal(t, FAILED_DELETE_IMAGE, err.Error())
	})
}

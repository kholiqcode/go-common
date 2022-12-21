package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kholiqcode/go-common/pkg/jwt"
	"github.com/kholiqcode/go-common/pkg/log"
	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Inside "helper" folder: JUST FOR TESTING PURPOSE
type TestCaseHandler struct {
	Name          string
	SetHeaders    func(req *http.Request)
	Payload       interface{}
	Method        string
	ReqUrl        string
	BuildStub     func(input interface{}, stubs ...interface{})
	CheckResponse func(recorder *httptest.ResponseRecorder, expected interface{})
}

func SetHeaderApplicationJson(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

func SetHeaderMultiPartForm(req *http.Request, contentType string) {
	req.Header.Add("Content-Type", contentType)
}

func SetAuthorizationHeader(req *http.Request, symmetricKey string, userId uuid.UUID) {
	logger := log.NewLogger("console", "debug")
	tokenMaker := jwt.NewManager(logger, &common_utils.Config{
		JWT: common_utils.JWT{
			Secret:  "secret",
			Expires: time.Hour * 2,
		},
	})

	token, _ := tokenMaker.Generate(userId)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
}

func ParseResponseBody(body *bytes.Buffer, output interface{}) {
	json.NewDecoder(body).Decode(output)
}

func ParseInterfaceToMap(data interface{}, field string) map[string]interface{} {
	return data.(map[string]interface{})[field].(map[string]interface{})
}

func ParseInterfaceToSlice(data interface{}, field string) []interface{} {
	return data.(map[string]interface{})[field].([]interface{})
}

func ParseInterfaceToString(data interface{}, field string) string {
	return data.(map[string]interface{})[field].(string)
}

func ParseErrorMessage(data interface{}) string {
	errors := data.([]interface{})

	errorsResp := ""
	for i, err := range errors {
		errMap := err.(map[string]interface{})
		if i > 0 {
			errorsResp = fmt.Sprintf("%s, %s", errorsResp, errMap["message"])

		} else {
			errorsResp = errMap["message"].(string)
		}

	}
	return errorsResp
}

func CreateFormFile(n int, filename string) (*bytes.Buffer, string) {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	for i := 0; i < n; i++ {
		newFileName := fmt.Sprintf("%s-%d", filename, i+1)
		formFile, _ := writer.CreateFormFile("files", newFileName)
		formFile.Write([]byte(fmt.Sprintf("JUST FOR TESTING %s", newFileName)))
	}
	writer.Close()

	return form, writer.FormDataContentType()
}

func CreateFilesHeader(n int, filename string) []*multipart.FileHeader {
	formFile, header := CreateFormFile(n, filename)

	req := httptest.NewRequest(http.MethodPost, "/upload", formFile)
	defer req.Body.Close()
	req.Header.Add("Content-Type", header)

	req.ParseMultipartForm(5242880)
	filesHeader := req.MultipartForm.File["files"]

	return filesHeader
}

func CreateFakeGooglePubSub(t *testing.T, project string) (gPubSub *pubsub.Client, close func()) {
	ctx := context.Background()
	// Start a fake server running locally.
	srv := pstest.NewServer()

	// Connect to the server without using TLS.
	conn, err := grpc.Dial(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	// Use the connection when creating a pubsub client.
	client, err := pubsub.NewClient(ctx, project, option.WithGRPCConn(conn))
	assert.NoError(t, err)

	_ = client
	return client, func() {
		srv.Close()
		conn.Close()
		client.Close()
	}
}

func CheckTokenPayloadCtx(ctx context.Context) gomock.Matcher {

	return gomock.AssignableToTypeOf(context.WithValue(ctx, gomock.Nil(), ""))
}

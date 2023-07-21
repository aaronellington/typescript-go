package typescript_test

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aaronellington/typescript-go/typescript"
)

type UserID uint64

type Group struct {
	Name      string `json:"groupName"`
	UpdatedAt time.Time
	DeletedAt *time.Time
	Data      any
	MoreData  interface{}
}

type User struct {
	UserID         UserID   `json:"userID"`
	PrimaryGroup   Group    `json:"primaryGroup"`
	SecondaryGroup *Group   `json:"secondaryGroup,omitempty"`
	Tags           []string `json:"tags"`
	Private        any      `json:"-"`
	unexported     any
}

func TestPrimary(t *testing.T) {
	_ = User{}.unexported

	service := typescript.New(map[string]any{
		"foobar":     UserID(0),
		"group":      Group{},
		"SystemUser": User{},
	})

	testThePackage(t, service)
}

func testThePackage(t *testing.T, service *typescript.Service) {
	actualFileName := "test_files/" + t.Name() + "_actual.ts"
	actualFile, err := os.Create(actualFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer actualFile.Close()

	actualFileBuffer := bytes.NewBuffer([]byte{})

	writer := io.MultiWriter(actualFile, actualFileBuffer)

	if err := service.Generate(writer); err != nil {
		t.Fatal(err)
	}

	expectedContents, err := os.ReadFile("test_files/" + t.Name() + "_expected.ts")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actualFileBuffer.Bytes(), expectedContents) {
		wd, _ := os.Getwd()
		t.Fatal("contents don't match: " + wd + "/" + actualFileName)
	}
}

package test

import (
	"bufio"
	"fmt"
	"github.com/DimaKuptsov/task-man/handlers"
	httpErrors "github.com/DimaKuptsov/task-man/handlers/error"
	"github.com/DimaKuptsov/task-man/handlers/projects"
	"github.com/DimaKuptsov/task-man/logger"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProjectsDeleteWithInvalidId(t *testing.T) {
	err := logger.Init()
	if err != nil {
		t.Errorf("Projects.DeleteRoute: failed to init logger. Error: %s", err.Error())
	}
	router := handlers.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()
	var testsWithInvalidIDs = []struct {
		invalidId string
	}{
		{"0"},
		{"123"},
		{"QIEw-12312qw-sad123"},
	}
	for _, test := range testsWithInvalidIDs {
		deleteUrl := fmt.Sprintf("%s/projects/delete/%s", srv.URL, test.invalidId)
		request, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
		if err != nil {
			t.Errorf("Projects.DeleteRoute: expected create request without errors. Got %s", err.Error())
		}
		client := &http.Client{}
		res, err := client.Do(request)
		if err != nil {
			t.Errorf("Projects.DeleteRoute: expected send request without errors. Got %s", err.Error())
		}
		reader := bufio.NewReader(res.Body)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			t.Errorf("Projects.DeleteRoute: expected response reading without errors. Got %s", err.Error())
		}
		if !strings.Contains(line, httpErrors.UnprocessableEntityMessage) {
			t.Errorf("Projects.DeleteRoute: expected unprocessable entity request. Got %s", line)
		}
		if !strings.Contains(line, httpErrors.GetBadParameterErrorMessage(projects.ProjectIDField)) {
			t.Errorf("Projects.DeleteRoute: expected bad parameter %s. Got %s", projects.ProjectIDField, line)
		}
		res.Body.Close()
	}
}

func TestProjectsDeleteExpectedInternalServerError(t *testing.T) {
	err := logger.Init()
	if err != nil {
		t.Errorf("Projects.DeleteRoute: failed to init logger. Error: %s", err.Error())
	}
	router := handlers.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()
	deleteUrl := fmt.Sprintf("%s/projects/delete/%s", srv.URL, uuid.New().String())
	request, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		t.Errorf("Projects.DeleteRoute: expected create request without errors. Got %s", err.Error())
	}
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		t.Errorf("Projects.DeleteRoute: expected send request without errors. Got %s", err.Error())
	}
	reader := bufio.NewReader(res.Body)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		t.Errorf("Projects.DeleteRoute: expected response reading without errors. Got %s", err.Error())
	}
	if !strings.Contains(line, httpErrors.InternalServerErrorMessage) {
		t.Errorf("Projects.DeleteRoute: expected internal server error response. Got %s", line)
	}
	if !strings.Contains(line, httpErrors.InternalServerErrorDescription) {
		t.Errorf("Projects.DeleteRoute: expected internal server error description. Got %s", line)
	}
	res.Body.Close()
}

package test

import (
	"bufio"
	"fmt"
	"github.com/DimaKuptsov/task-man/handlers"
	"github.com/DimaKuptsov/task-man/handlers/columns"
	httpErrors "github.com/DimaKuptsov/task-man/handlers/error"
	"github.com/DimaKuptsov/task-man/logger"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestColumnsGetForProjectWithInvalidID(t *testing.T) {
	err := logger.Init()
	if err != nil {
		t.Errorf("Columns.GetForProjectRoute: failed to init logger. Error: %s", err.Error())
	}
	router := handlers.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()
	var testsWithInvalidIds = []struct {
		invalidId string
	}{
		{"0"},
		{"123"},
		{"QIEw-12312qw-sad123"},
	}
	for _, test := range testsWithInvalidIds {
		getUrl := fmt.Sprintf("%s/columns/project/%s", srv.URL, test.invalidId)
		res, err := http.Get(getUrl)
		if err != nil {
			t.Errorf("Columns.GetForProjectRoute: expected send request without errors. Got %s", err.Error())
		}
		reader := bufio.NewReader(res.Body)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			t.Errorf("Columns.GetForProjectRoute: expected response reading without errors. Got %s", err.Error())
		}
		if !strings.Contains(line, httpErrors.BadRequestMessage) {
			t.Errorf("Columns.GetForProjectRoute: expected bad request. Got %s", line)
		}
		if !strings.Contains(line, httpErrors.GetBadParameterErrorMessage(columns.ProjectIDField)) {
			t.Errorf("Columns.GetForProjectRoute: expected bad parameter %s. Got %s", columns.ProjectIDField, line)
		}
		res.Body.Close()
	}
}

func TestColumnsGetForProjectExpectedInternalServerError(t *testing.T) {
	err := logger.Init()
	if err != nil {
		t.Errorf("Columns.GetForProjectRoute: failed to init logger. Error: %s", err.Error())
	}
	router := handlers.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()
	getUrl := fmt.Sprintf("%s/columns/project/%s", srv.URL, uuid.New().String())
	res, err := http.Get(getUrl)
	if err != nil {
		t.Errorf("Columns.GetForProjectRoute: expected send request without errors. Got %s", err.Error())
	}
	reader := bufio.NewReader(res.Body)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		t.Errorf("Columns.GetForProjectRoute: expected response reading without errors. Got %s", err.Error())
	}
	if !strings.Contains(line, httpErrors.InternalServerErrorMessage) {
		t.Errorf("Columns.GetForProjectRoute: expected internal server error response. Got %s", line)
	}
	if !strings.Contains(line, httpErrors.InternalServerErrorDescription) {
		t.Errorf("Columns.GetForProjectRoute: expected internal server error description. Got %s", line)
	}
	res.Body.Close()
}

package service

import (
	"errors"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"
	"github.com/stretchr/testify/assert"
)

var (
	src     = rand.New(rand.NewSource(time.Now().UnixNano()))
	formats = [...]string{
		"pretty",
		"yaml",
		"json",
	}
)

func genCfg() *model.Config {
	return model.NewConfig(
		helpers.RandString(30),
		formats[rand.Intn(len(formats))],
		int(src.Int63()))
}

func genWl() *model.Work {
	tags := make([]string, src.Intn(128))
	for index := range tags {
		tags[index] = helpers.RandString(src.Intn(30))
	}

	return model.NewWork(
		helpers.RandString(30),
		helpers.RandString(30),
		helpers.RandString(30),
		int(src.Int63()),
		tags,
		time.Now(),
	)
}

func TestConfigure(t *testing.T) {
	var tests = []struct {
		name string
		cfg  *model.Config
		err  error
	}{
		{
			name: "Success",
			cfg:  genCfg(),
			err:  nil,
		},
		{
			name: "Errored",
			cfg:  genCfg(),
			err:  errors.New("Errored"),
		},
	}

	for _, testItem := range tests {
		mockRepo := new(repository.MockRepo)
		mockRepo.On("Configure", testItem.cfg).
			Return(testItem.err)
		svc := NewWorklogService(mockRepo)

		t.Run(testItem.name, func(t *testing.T) {
			returnedErr := svc.Congfigure(testItem.cfg)

			if returnedErr != nil {
				assert.EqualError(t, testItem.err, returnedErr.Error())
			} else {
				assert.Nil(t, returnedErr)
			}
			mockRepo.AssertExpectations(t)
			mockRepo.AssertCalled(t, "Configure", testItem.cfg)
		})
	}
}

func TestCreateWorklog(t *testing.T) {
	var tests = []struct {
		name   string
		wl     *model.Work
		exCode int
		err    error
	}{
		{
			name:   "Success",
			wl:     genWl(),
			exCode: http.StatusCreated,
			err:    nil,
		},
		{
			name:   "Errored",
			wl:     genWl(),
			exCode: http.StatusInternalServerError,
			err:    errors.New("Errored"),
		},
	}

	for _, testItem := range tests {
		mockRepo := new(repository.MockRepo)
		mockRepo.On("Save", testItem.wl).
			Return(testItem.err)
		svc := NewWorklogService(mockRepo)

		t.Run(testItem.name, func(t *testing.T) {
			returnedCode, returnedErr := svc.CreateWorklog(testItem.wl)

			if returnedErr != nil {
				assert.EqualError(t, testItem.err, returnedErr.Error())
			} else {
				assert.Nil(t, returnedErr)
			}
			assert.Equal(t, testItem.exCode, returnedCode)
			mockRepo.AssertExpectations(t)
			mockRepo.AssertCalled(t, "Save", testItem.wl)
		})
	}
}

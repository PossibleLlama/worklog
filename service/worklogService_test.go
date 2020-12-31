package service

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"
	"github.com/stretchr/testify/assert"
)

var (
	src     = rand.NewSource(time.Now().UnixNano())
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

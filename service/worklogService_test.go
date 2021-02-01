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

const (
	strLength = 30
	arrLength = 128
)

func genCfg() *model.Config {
	return model.NewConfig(
		helpers.RandString(strLength),
		formats[rand.Intn(len(formats))],
		int(src.Int63()))
}

func genWl() *model.Work {
	tags := make([]string, src.Intn(arrLength))
	for index := range tags {
		tags[index] = helpers.RandString(src.Intn(strLength))
	}

	return model.NewWork(
		helpers.RandString(strLength),
		helpers.RandString(strLength),
		helpers.RandString(strLength),
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
			err:  errors.New(helpers.RandString(strLength)),
		},
	}

	for _, testItem := range tests {
		mockRepo := new(repository.MockRepo)
		mockRepo.On("Configure", testItem.cfg).
			Return(testItem.err)
		svc := NewWorklogService(mockRepo)

		t.Run(testItem.name, func(t *testing.T) {
			returnedErr := svc.Configure(testItem.cfg)

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
			err:    errors.New(helpers.RandString(strLength)),
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

func TestGetWorklogsBetween(t *testing.T) {
	sTime := time.Now()
	eTime := sTime.Add(time.Hour)
	rev1Wl := genWl()
	rev2Wl := &model.Work{
		ID:          rev1Wl.ID,
		Revision:    rev1Wl.Revision + 1,
		Title:       rev1Wl.Title,
		Description: rev1Wl.Description,
		Author:      rev1Wl.Author,
		Duration:    rev1Wl.Duration,
		Tags:        rev1Wl.Tags,
		When:        rev1Wl.When,
		CreatedAt:   rev1Wl.CreatedAt,
	}
	rev3Wl := &model.Work{
		ID:          rev1Wl.ID,
		Revision:    rev1Wl.Revision + 2,
		Title:       rev1Wl.Title,
		Description: rev1Wl.Description,
		Author:      rev1Wl.Author,
		Duration:    rev1Wl.Duration,
		Tags:        rev1Wl.Tags,
		When:        rev1Wl.When,
		CreatedAt:   rev1Wl.CreatedAt,
	}

	wl2 := genWl()
	wl2.When = wl2.When.Add(time.Minute)
	wl3 := genWl()
	wl3.When = wl3.When.Add(time.Minute * 2)
	wl4 := genWl()
	wl4.When = wl4.When.Add(time.Minute * 3)

	var tests = []struct {
		name   string
		sTime  time.Time
		eTime  time.Time
		filter *model.Work
		retWl  []*model.Work
		expWl  []*model.Work
		exCode int
		err    error
	}{
		{
			name:   "Success found 1",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{rev1Wl},
			expWl:  []*model.Work{rev1Wl},
			exCode: http.StatusOK,
			err:    nil,
		},
		{
			name:   "Success found 0",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{},
			expWl:  []*model.Work{},
			exCode: http.StatusNotFound,
			err:    nil,
		},
		{
			name:   "Success deduplicates only has latest revision from ordered list",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{rev1Wl, rev2Wl},
			expWl:  []*model.Work{rev2Wl},
			exCode: http.StatusOK,
			err:    nil,
		},
		{
			name:   "Success deduplicates only has latest revision from list",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{rev3Wl, rev1Wl, rev2Wl},
			expWl:  []*model.Work{rev3Wl},
			exCode: http.StatusOK,
			err:    nil,
		},
		{
			name:   "Success sorts into order",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{wl3, rev1Wl, wl4, wl2},
			expWl:  []*model.Work{rev1Wl, wl2, wl3, wl4},
			exCode: http.StatusOK,
			err:    nil,
		},
		{
			name:   "Errored",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{},
			expWl:  []*model.Work{},
			exCode: http.StatusInternalServerError,
			err:    errors.New(helpers.RandString(strLength)),
		},
	}

	for _, testItem := range tests {
		mockRepo := new(repository.MockRepo)
		mockRepo.On(
			"GetAllBetweenDates",
			testItem.sTime,
			testItem.eTime,
			testItem.filter).
			Return(testItem.retWl, testItem.err)
		svc := NewWorklogService(mockRepo)

		t.Run(testItem.name, func(t *testing.T) {
			returnedWls, returnedCode, returnedErr := svc.GetWorklogsBetween(
				testItem.sTime,
				testItem.eTime,
				testItem.filter)

			if returnedErr != nil {
				assert.EqualError(t, testItem.err, returnedErr.Error())
			} else {
				assert.Nil(t, returnedErr)
			}
			assert.Equal(t, testItem.exCode, returnedCode)
			assert.Equal(t, len(testItem.expWl), len(returnedWls))
			assert.Equal(t, testItem.expWl, returnedWls)
			for index := range testItem.expWl {
				assert.Equal(t, testItem.expWl[index], returnedWls[index])
			}
			mockRepo.AssertExpectations(t)
			mockRepo.AssertCalled(t,
				"GetAllBetweenDates",
				testItem.sTime,
				testItem.eTime,
				testItem.filter)
		})
	}
}

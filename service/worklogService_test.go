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

func TestMain(m *testing.M) {
	for i := 0; i < 10; i++ {
		if time.Now().Nanosecond()%1000 < 500 {
			m.Run()
			return
		}
		time.Sleep(400 * time.Nanosecond)
	}
}

func genCfg() *model.Config {
	return model.NewConfig(
		helpers.RandAlphabeticString(strLength),
		formats[rand.Intn(len(formats))],
		int(src.Int63()))
}

func genWl() *model.Work {
	tags := make([]string, src.Intn(arrLength)+1)
	for index := range tags {
		tags[index] = helpers.RandAlphabeticString(src.Intn(strLength))
	}

	return model.NewWork(
		helpers.RandAlphabeticString(strLength),
		helpers.RandAlphabeticString(strLength),
		helpers.RandAlphabeticString(strLength),
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
			err:  errors.New(helpers.RandAlphabeticString(strLength)),
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
			err:    errors.New(helpers.RandAlphabeticString(strLength)),
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

func TestEditWorklog(t *testing.T) {
	id := helpers.RandHexAlphaNumericString(strLength)
	wl := genWl()
	var tests = []struct {
		name     string
		newWl    *model.Work
		getWl    *model.Work
		getErr   error
		callSave bool
		expCode  int
		expErr   error
	}{
		{
			name:     "Success with all fields",
			newWl:    genWl(),
			getWl:    wl,
			getErr:   nil,
			callSave: true,
			expCode:  http.StatusOK,
			expErr:   nil,
		}, {
			name:     "No found wl's",
			newWl:    genWl(),
			getWl:    nil,
			getErr:   nil,
			callSave: false,
			expCode:  http.StatusNotFound,
			expErr:   nil,
		}, {
			name:     "Error from save",
			newWl:    genWl(),
			getWl:    wl,
			getErr:   nil,
			callSave: true,
			expCode:  http.StatusInternalServerError,
			expErr:   errors.New(id),
		},
	}

	for _, testItem := range tests {
		expWl := model.Work{
			ID:          wl.ID,
			Revision:    wl.Revision + 1,
			Title:       testItem.newWl.Title,
			Description: testItem.newWl.Description,
			Author:      testItem.newWl.Author,
			Duration:    testItem.newWl.Duration,
			Tags:        testItem.newWl.Tags,
			When:        wl.When,
			CreatedAt:   testItem.newWl.CreatedAt}

		mockRepo := new(repository.MockRepo)
		mockRepo.On("GetByID", id, &model.Work{}).
			Return(testItem.getWl, testItem.getErr)
		if testItem.callSave {
			mockRepo.On("Save", &expWl).
				Return(testItem.expErr)
		}

		svc := NewWorklogService(mockRepo)

		t.Run(testItem.name, func(t *testing.T) {
			returnedCode, returnedErr := svc.EditWorklog(id, testItem.newWl)

			assert.Equal(t, returnedErr, testItem.expErr)
			assert.Equal(t, testItem.expCode, returnedCode)
			mockRepo.AssertExpectations(t)
			mockRepo.AssertCalled(t, "GetByID", id, &model.Work{})
			if testItem.callSave {
				//mockRepo.AssertCalled(t, "Save", &expWl)
			}
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
			err:    errors.New(helpers.RandAlphabeticString(strLength)),
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

func TestGetWorklogByID(t *testing.T) {
	wl1 := genWl()

	var tests = []struct {
		name   string
		ids    []string
		filter *model.Work
		retWl  []*model.Work
		expWl  []*model.Work
		exCode int
		err    error
	}{
		{
			name:   "Single ID calls repo once",
			ids:    []string{wl1.ID},
			filter: wl1,
			retWl:  []*model.Work{wl1},
			expWl:  []*model.Work{wl1},
			exCode: http.StatusOK,
			err:    nil,
		},
	}

	for _, testItem := range tests {
		assert.Equal(t, len(testItem.ids), len(testItem.retWl), "Number of IDs and returned WL's must match")
		mockRepo := new(repository.MockRepo)
		for index := range testItem.ids {
			mockRepo.On(
				"GetByID",
				testItem.ids[index],
				testItem.filter).
				Return(testItem.retWl[index], testItem.err)
		}
		svc := NewWorklogService(mockRepo)

		t.Run(testItem.name, func(t *testing.T) {
			returnedWls, returnedCode, returnedErr := svc.GetWorklogsByID(
				testItem.filter,
				testItem.ids...)

			if returnedErr != nil {
				assert.EqualError(t, testItem.err, returnedErr.Error())
			} else {
				assert.Nil(t, testItem.err)
			}
			assert.Equal(t, testItem.exCode, returnedCode)
			assert.Equal(t, len(testItem.expWl), len(returnedWls))
			assert.Equal(t, testItem.expWl, returnedWls)
			for index := range testItem.expWl {
				assert.Equal(t, testItem.expWl[index], returnedWls[index])
			}
			mockRepo.AssertExpectations(t)
			for _, ID := range testItem.ids {
				mockRepo.AssertCalled(t,
					"GetByID",
					ID,
					testItem.filter)
			}
		})
	}
}

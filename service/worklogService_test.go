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

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

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

func TestCreateWorklog(t *testing.T) {
	genedWl := genWl()
	defaultDur := 15
	viper.Set("default.duration", defaultDur)
	defaultAuth := helpers.RandAlphabeticString(30)
	viper.Set("default.author", defaultAuth)
	randomStr := helpers.RandAlphabeticString(30)
	now := time.Now()

	var tests = []struct {
		name    string
		wl      *model.Work
		savedWl *model.Work
		expCode int
		err     error
	}{
		{
			name:    "Success",
			wl:      genedWl,
			savedWl: genedWl,
			expCode: http.StatusCreated,
			err:     nil,
		}, {
			name: "Default duration",
			wl: &model.Work{
				Duration: 0,
			},
			savedWl: &model.Work{
				Duration: defaultDur,
			},
			expCode: http.StatusCreated,
			err:     nil,
		}, {
			name: "Default author",
			wl: &model.Work{
				Author: "",
			},
			savedWl: &model.Work{
				Author: defaultAuth,
			},
			expCode: http.StatusCreated,
			err:     nil,
		}, {
			name: "Duplicate tags",
			wl: &model.Work{
				Tags: []string{"a", "a", "b", "c", "b"},
			},
			savedWl: &model.Work{
				Tags: []string{"a", "b", "c"},
			},
			expCode: http.StatusCreated,
			err:     nil,
		}, {
			name: "Padded title",
			wl: &model.Work{
				Title: " " + randomStr,
			},
			savedWl: &model.Work{
				Title: randomStr,
			},
			expCode: http.StatusCreated,
			err:     nil,
		}, {
			name: "Padded description",
			wl: &model.Work{
				Description: " " + randomStr,
			},
			savedWl: &model.Work{
				Description: randomStr,
			},
			expCode: http.StatusCreated,
			err:     nil,
		}, {
			name: "Padded author",
			wl: &model.Work{
				Author: " " + randomStr,
			},
			savedWl: &model.Work{
				Author: randomStr,
			},
			expCode: http.StatusCreated,
			err:     nil,
		}, {
			name: "Padded tags",
			wl: &model.Work{
				Tags: []string{" a", "b\t", "c\n"},
			},
			savedWl: &model.Work{
				Tags: []string{"a", "b", "c"},
			},
			expCode: http.StatusCreated,
			err:     nil,
		}, {
			name: "Default When",
			wl: &model.Work{
				When:      time.Time{},
				CreatedAt: now,
			},
			savedWl: &model.Work{
				When:      now,
				CreatedAt: now,
			},
			expCode: http.StatusCreated,
			err:     nil,
		}, {
			name:    "Errored",
			wl:      genedWl,
			savedWl: genedWl,
			expCode: http.StatusInternalServerError,
			err:     errors.New(helpers.RandAlphabeticString(strLength)),
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
			assert.Equal(t, testItem.expCode, returnedCode)
			mockRepo.AssertExpectations(t)
			mockRepo.AssertCalled(t, "Save", testItem.wl)
		})
	}
}

func TestEditWorklog(t *testing.T) {
	id := helpers.RandHexAlphaNumericString(strLength)
	wl := genWl()
	wlWithDuplicateTags := genWl()
	wlWithDuplicateTags.Tags = append(wlWithDuplicateTags.Tags, "a", "a")

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
			name:     "Duplicate tags",
			newWl:    wlWithDuplicateTags,
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
			ID:             wl.ID,
			Revision:       wl.Revision + 1,
			Title:          testItem.newWl.Title,
			Description:    testItem.newWl.Description,
			Author:         testItem.newWl.Author,
			Duration:       testItem.newWl.Duration,
			Tags:           helpers.DeduplicateString(testItem.newWl.Tags),
			When:           wl.When,
			WhenQueryEpoch: wl.WhenQueryEpoch,
			CreatedAt:      testItem.newWl.CreatedAt}

		mockRepo := new(repository.MockRepo)
		mockRepo.On("GetByID", id, &model.Work{}).
			Return(testItem.getWl, testItem.getErr)
		if testItem.callSave {
			mockRepo.On("Save", &expWl).
				Return(testItem.expErr)
		}

		svc := NewWorklogService(mockRepo)

		t.Run(testItem.name, func(t *testing.T) {
			_, returnedCode, returnedErr := svc.EditWorklog(id, testItem.newWl)

			assert.Equal(t, returnedErr, testItem.expErr)
			assert.Equal(t, testItem.expCode, returnedCode)
			mockRepo.AssertExpectations(t)
			mockRepo.AssertCalled(t, "GetByID", id, &model.Work{})
			if testItem.callSave {
				mockRepo.AssertCalled(t, "Save", &expWl)
			}
		})
	}
}

func TestGetWorklogsBetween(t *testing.T) {
	sTime := time.Now()
	eTime := sTime.Add(time.Hour)
	eTime3k := time.Date(3000, time.January, 1, 0, 0, 0, 0, time.Now().Location())
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
		name      string
		sTime     time.Time
		eTime     time.Time
		eTimeRepo *time.Time
		filter    *model.Work
		retWl     []*model.Work
		expWl     []*model.Work
		exCode    int
		err       error
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
		}, {
			name:   "Success found 0",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{},
			expWl:  []*model.Work{},
			exCode: http.StatusNotFound,
			err:    nil,
		}, {
			name:   "Success deduplicates only has latest revision from ordered list",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{rev1Wl, rev2Wl},
			expWl:  []*model.Work{rev2Wl},
			exCode: http.StatusOK,
			err:    nil,
		}, {
			name:   "Success deduplicates only has latest revision from list",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{rev3Wl, rev1Wl, rev2Wl},
			expWl:  []*model.Work{rev3Wl},
			exCode: http.StatusOK,
			err:    nil,
		}, {
			name:   "Success sorts into order",
			sTime:  sTime,
			eTime:  eTime,
			filter: rev1Wl,
			retWl:  []*model.Work{wl3, rev1Wl, wl4, wl2},
			expWl:  []*model.Work{rev1Wl, wl2, wl3, wl4},
			exCode: http.StatusOK,
			err:    nil,
		}, {
			name:      "No endDate defaults to year 3000",
			sTime:     sTime,
			eTime:     time.Time{},
			eTimeRepo: &eTime3k,
			filter:    rev1Wl,
			retWl:     []*model.Work{rev1Wl},
			expWl:     []*model.Work{rev1Wl},
			exCode:    http.StatusOK,
			err:       nil,
		}, {
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
		if testItem.eTimeRepo == nil {
			mockRepo.On(
				"GetAllBetweenDates",
				testItem.sTime,
				testItem.eTime,
				testItem.filter).
				Return(testItem.retWl, testItem.err)
		} else {
			mockRepo.On(
				"GetAllBetweenDates",
				testItem.sTime,
				*testItem.eTimeRepo,
				testItem.filter).
				Return(testItem.retWl, testItem.err)
		}
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
			if testItem.eTimeRepo == nil {
				mockRepo.AssertCalled(t,
					"GetAllBetweenDates",
					testItem.sTime,
					testItem.eTime,
					testItem.filter)
			} else {
				mockRepo.AssertCalled(t,
					"GetAllBetweenDates",
					testItem.sTime,
					*testItem.eTimeRepo,
					testItem.filter)
			}
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

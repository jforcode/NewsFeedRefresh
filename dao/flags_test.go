package dao

import (
	"testing"

	"github.com/golang/glog"
	"github.com/jforcode/DbUtil"
)

func TestFlags(t *testing.T) {
	db, dbMain := dbTestInit()
	defer db.Close()

	dbUtil.ClearTables(db, "news_api_flags")

	t.Run("get non-existent flag", func(t *testing.T) {
		defer dbUtil.ClearTables(db, "news_api_flags")
		flag, err := dbMain.GetFlag("test", "int")

		if !(err == nil && flag == nil) {
			t.FailNow()
		}
	})

	t.Run("set and get non-existent flag", func(t *testing.T) {
		tests := []struct {
			name     string
			key      string
			value    string
			expected interface{}
			typeTo   FlagType
		}{
			{"string", "test", "asdf", "asdf", FlagTypeString},
			{"int", "test", "1234", 1234, FlagTypeInt},
			{"boolean", "test", "TRUE", true, FlagTypeBool},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				defer dbUtil.ClearTables(db, "news_api_flags")

				err := dbMain.SetFlag(test.key, test.value, test.typeTo)
				if !(err == nil) {
					t.FailNow()
				}

				flag, err := dbMain.GetFlag(test.key, test.typeTo)
				if !(err == nil && flag != nil && test.expected == flag.Value) {
					t.FailNow()
				}
			})
		}
	})

	t.Run("update flag", func(t *testing.T) {
		tests := []struct {
			name         string
			key          string
			value        string
			updatedValue string
			expected     interface{}
			typeTo       FlagType
		}{
			{"string", "test", "asdf", "def", "def", FlagTypeString},
			{"int", "test", "1234", "5678", 5678, FlagTypeInt},
			{"boolean", "test", "TRUE", "FALSE", false, FlagTypeBool},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				defer dbUtil.ClearTables(db, "news_api_flags")

				dbMain.SetFlag(test.key, test.value, test.typeTo)

				err := dbMain.SetFlag(test.key, test.updatedValue, test.typeTo)
				if !(err == nil) {
					glog.Errorln(err)
					t.FailNow()
				}

				flag, err := dbMain.GetFlag(test.key, test.typeTo)
				if !(err == nil && flag != nil && test.expected == flag.Value) {
					glog.Errorln(err, flag, test.expected, flag.Value)
					t.FailNow()
				}

				dbUtil.GetRowCount(db, "news_api_flags", "flag_key = ?", []interface{}{test.key})
			})
		}
	})

	// TODO: type mismatch in flags. conversion errors. setting in different type, fetching different type
}

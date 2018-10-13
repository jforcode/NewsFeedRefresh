package dao

import (
	"testing"

	"github.com/golang/glog"
	"github.com/jforcode/DbUtil"
)

func areSourcesEqual(source1 *Source, source2 *Source) bool {
	return source1.Name == source2.Name && source1.SourceId == source2.SourceId
}

func TestSources(t *testing.T) {
	db, dbMain := dbTestInit()
	defer db.Close()

	dbUtil.ClearTables(db, "sources")

	t.Run("save and get sources", func(t *testing.T) {
		sources := make([]*Source, 2)
		sources[0] = &Source{SourceId: "test", Name: "testSource"}
		sources[1] = &Source{SourceId: "test2", Name: "testSource2"}

		numSaved, err := dbMain.SaveSources(sources)
		if !(err == nil && numSaved == 2) {
			glog.Errorln(err)
			t.FailNow()
		}

		rowCount, err := dbUtil.GetRowCount(db, "sources", "", []interface{}{})
		if err != nil || rowCount != 2 {
			t.FailNow()
		}

		gotSources, err := dbMain.GetSources()
		if !(err == nil &&
			len(gotSources) == 2 &&
			areSourcesEqual(gotSources[0], sources[0]) &&
			areSourcesEqual(gotSources[1], sources[1])) {

			t.FailNow()
		}
	})
}

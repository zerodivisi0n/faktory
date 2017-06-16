package storage

import (
	"testing"
	"time"

	"github.com/mperham/worq/util"
	"github.com/stretchr/testify/assert"
)

func TestBasicTimedSet(t *testing.T) {
	t.Parallel()

	db, err := OpenStore("../test/timed.db")
	assert.NoError(t, err)
	j1 := []byte(fakeJob())

	past := time.Now()

	assert.Equal(t, 0, db.Retries().Size())
	err = db.Retries().AddElement(past, "1239712983", j1)
	assert.NoError(t, err)
	assert.Equal(t, 1, db.Retries().Size())

	j2 := []byte(fakeJob())
	err = db.Retries().AddElement(past, "1239712984", j2)
	assert.NoError(t, err)
	assert.Equal(t, 2, db.Retries().Size())

	current := time.Now()
	err = db.Retries().AddElement(current.Add(10*time.Second), "1239712985", []byte(fakeJob()))
	assert.NoError(t, err)
	assert.Equal(t, 3, db.Retries().Size())

	results, err := db.Retries().RemoveBefore(current.Add(1 * time.Second))
	assert.NoError(t, err)
	assert.Equal(t, 1, db.Retries().Size())
	assert.Equal(t, 2, len(results))
	values := [][]byte{j1, j2}
	assert.Equal(t, values, results)
}

//func TestTimestampFormat(t *testing.T) {
//tstamp := time.Now()
//jid := "aksdfask"
//key := fmt.Sprintf("%.10d|%s", tstamp.Unix(), jid)

//fmt.Println(key)
//}

func fakeJob() string {
	return `{"jid":"` + util.RandomJid() + `","created_at":1234567890.123,"queue":"default","args":[1,2,3]}`
}
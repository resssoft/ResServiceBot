package testManager

import (
	"errors"
	"strings"
)

func (d *data) getTest(code string) (*testParent, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if item, ok := d.tests[code]; ok {
		return item, nil
	}
	return nil, errors.New("not found")
}

func (d *data) setTest(newItem *testParent) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if newItem == nil {
		return errors.New("empty test")
	}
	d.tests[newItem.Code] = newItem
	return nil
}

func (d *data) appendTestRows(data string, newItem *testParent) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if newItem == nil {
		return errors.New("empty test")
	}
	testData := strings.Split(data, "\n")
	if len(testData) < 2 {
		return errors.New("low row numbers")
	}
	newItem.Children = append(newItem.Children, &testChild{
		Question:   testData[0],
		Answers:    testData[1:],
		Variants:   testData[1:],
		Incorrect:  "wrong!",
		Correct:    "correct",
		ShowAnswer: true,
	})
	return nil
}

func (d *data) getActiveTest(uid int64) (*testParent, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if code, ok := d.active[uid]; ok {
		if item, ok := d.tests[code]; ok {
			return item, nil
		}
		return nil, errors.New("not found test")
	}
	return nil, errors.New("not found active test")
}

func (d *data) setActiveTest(uid int64, code string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.active[uid] = code
}

func (d *data) delActiveTest(uid int64) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.active, uid)
}

func (d *data) getUserTest(uid int64) (*testParent, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if item, ok := d.users[uid]; ok {
		return item, nil
	}
	return nil, errors.New("not found users test")
}

func (d *data) setAUserTest(uid int64, item *testParent) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.users[uid] = item
}

func (d *data) delUserTest(uid int64) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.users, uid)
}

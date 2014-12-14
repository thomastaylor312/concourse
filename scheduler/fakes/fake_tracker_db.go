// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/concourse/atc/db"
	"github.com/concourse/atc/scheduler"
)

type FakeTrackerDB struct {
	SaveBuildStatusStub        func(buildID int, status db.Status) error
	saveBuildStatusMutex       sync.RWMutex
	saveBuildStatusArgsForCall []struct {
		buildID int
		status  db.Status
	}
	saveBuildStatusReturns struct {
		result1 error
	}
}

func (fake *FakeTrackerDB) SaveBuildStatus(buildID int, status db.Status) error {
	fake.saveBuildStatusMutex.Lock()
	fake.saveBuildStatusArgsForCall = append(fake.saveBuildStatusArgsForCall, struct {
		buildID int
		status  db.Status
	}{buildID, status})
	fake.saveBuildStatusMutex.Unlock()
	if fake.SaveBuildStatusStub != nil {
		return fake.SaveBuildStatusStub(buildID, status)
	} else {
		return fake.saveBuildStatusReturns.result1
	}
}

func (fake *FakeTrackerDB) SaveBuildStatusCallCount() int {
	fake.saveBuildStatusMutex.RLock()
	defer fake.saveBuildStatusMutex.RUnlock()
	return len(fake.saveBuildStatusArgsForCall)
}

func (fake *FakeTrackerDB) SaveBuildStatusArgsForCall(i int) (int, db.Status) {
	fake.saveBuildStatusMutex.RLock()
	defer fake.saveBuildStatusMutex.RUnlock()
	return fake.saveBuildStatusArgsForCall[i].buildID, fake.saveBuildStatusArgsForCall[i].status
}

func (fake *FakeTrackerDB) SaveBuildStatusReturns(result1 error) {
	fake.SaveBuildStatusStub = nil
	fake.saveBuildStatusReturns = struct {
		result1 error
	}{result1}
}

var _ scheduler.TrackerDB = new(FakeTrackerDB)

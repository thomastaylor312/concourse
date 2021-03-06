// Code generated by counterfeiter. DO NOT EDIT.
package schedulerfakes

import (
	sync "sync"
	time "time"

	lager "code.cloudfoundry.org/lager"
	atc "github.com/concourse/concourse/atc"
	db "github.com/concourse/concourse/atc/db"
	algorithm "github.com/concourse/concourse/atc/db/algorithm"
	scheduler "github.com/concourse/concourse/atc/scheduler"
)

type FakeBuildScheduler struct {
	SaveNextInputMappingStub        func(lager.Logger, db.Job, db.Resources) error
	saveNextInputMappingMutex       sync.RWMutex
	saveNextInputMappingArgsForCall []struct {
		arg1 lager.Logger
		arg2 db.Job
		arg3 db.Resources
	}
	saveNextInputMappingReturns struct {
		result1 error
	}
	saveNextInputMappingReturnsOnCall map[int]struct {
		result1 error
	}
	ScheduleStub        func(lager.Logger, *algorithm.VersionsDB, []db.Job, db.Resources, atc.VersionedResourceTypes) (map[string]time.Duration, error)
	scheduleMutex       sync.RWMutex
	scheduleArgsForCall []struct {
		arg1 lager.Logger
		arg2 *algorithm.VersionsDB
		arg3 []db.Job
		arg4 db.Resources
		arg5 atc.VersionedResourceTypes
	}
	scheduleReturns struct {
		result1 map[string]time.Duration
		result2 error
	}
	scheduleReturnsOnCall map[int]struct {
		result1 map[string]time.Duration
		result2 error
	}
	TriggerImmediatelyStub        func(lager.Logger, db.Job, db.Resources, atc.VersionedResourceTypes) (db.Build, scheduler.Waiter, error)
	triggerImmediatelyMutex       sync.RWMutex
	triggerImmediatelyArgsForCall []struct {
		arg1 lager.Logger
		arg2 db.Job
		arg3 db.Resources
		arg4 atc.VersionedResourceTypes
	}
	triggerImmediatelyReturns struct {
		result1 db.Build
		result2 scheduler.Waiter
		result3 error
	}
	triggerImmediatelyReturnsOnCall map[int]struct {
		result1 db.Build
		result2 scheduler.Waiter
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBuildScheduler) SaveNextInputMapping(arg1 lager.Logger, arg2 db.Job, arg3 db.Resources) error {
	fake.saveNextInputMappingMutex.Lock()
	ret, specificReturn := fake.saveNextInputMappingReturnsOnCall[len(fake.saveNextInputMappingArgsForCall)]
	fake.saveNextInputMappingArgsForCall = append(fake.saveNextInputMappingArgsForCall, struct {
		arg1 lager.Logger
		arg2 db.Job
		arg3 db.Resources
	}{arg1, arg2, arg3})
	fake.recordInvocation("SaveNextInputMapping", []interface{}{arg1, arg2, arg3})
	fake.saveNextInputMappingMutex.Unlock()
	if fake.SaveNextInputMappingStub != nil {
		return fake.SaveNextInputMappingStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.saveNextInputMappingReturns
	return fakeReturns.result1
}

func (fake *FakeBuildScheduler) SaveNextInputMappingCallCount() int {
	fake.saveNextInputMappingMutex.RLock()
	defer fake.saveNextInputMappingMutex.RUnlock()
	return len(fake.saveNextInputMappingArgsForCall)
}

func (fake *FakeBuildScheduler) SaveNextInputMappingCalls(stub func(lager.Logger, db.Job, db.Resources) error) {
	fake.saveNextInputMappingMutex.Lock()
	defer fake.saveNextInputMappingMutex.Unlock()
	fake.SaveNextInputMappingStub = stub
}

func (fake *FakeBuildScheduler) SaveNextInputMappingArgsForCall(i int) (lager.Logger, db.Job, db.Resources) {
	fake.saveNextInputMappingMutex.RLock()
	defer fake.saveNextInputMappingMutex.RUnlock()
	argsForCall := fake.saveNextInputMappingArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeBuildScheduler) SaveNextInputMappingReturns(result1 error) {
	fake.saveNextInputMappingMutex.Lock()
	defer fake.saveNextInputMappingMutex.Unlock()
	fake.SaveNextInputMappingStub = nil
	fake.saveNextInputMappingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildScheduler) SaveNextInputMappingReturnsOnCall(i int, result1 error) {
	fake.saveNextInputMappingMutex.Lock()
	defer fake.saveNextInputMappingMutex.Unlock()
	fake.SaveNextInputMappingStub = nil
	if fake.saveNextInputMappingReturnsOnCall == nil {
		fake.saveNextInputMappingReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.saveNextInputMappingReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildScheduler) Schedule(arg1 lager.Logger, arg2 *algorithm.VersionsDB, arg3 []db.Job, arg4 db.Resources, arg5 atc.VersionedResourceTypes) (map[string]time.Duration, error) {
	var arg3Copy []db.Job
	if arg3 != nil {
		arg3Copy = make([]db.Job, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.scheduleMutex.Lock()
	ret, specificReturn := fake.scheduleReturnsOnCall[len(fake.scheduleArgsForCall)]
	fake.scheduleArgsForCall = append(fake.scheduleArgsForCall, struct {
		arg1 lager.Logger
		arg2 *algorithm.VersionsDB
		arg3 []db.Job
		arg4 db.Resources
		arg5 atc.VersionedResourceTypes
	}{arg1, arg2, arg3Copy, arg4, arg5})
	fake.recordInvocation("Schedule", []interface{}{arg1, arg2, arg3Copy, arg4, arg5})
	fake.scheduleMutex.Unlock()
	if fake.ScheduleStub != nil {
		return fake.ScheduleStub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.scheduleReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBuildScheduler) ScheduleCallCount() int {
	fake.scheduleMutex.RLock()
	defer fake.scheduleMutex.RUnlock()
	return len(fake.scheduleArgsForCall)
}

func (fake *FakeBuildScheduler) ScheduleCalls(stub func(lager.Logger, *algorithm.VersionsDB, []db.Job, db.Resources, atc.VersionedResourceTypes) (map[string]time.Duration, error)) {
	fake.scheduleMutex.Lock()
	defer fake.scheduleMutex.Unlock()
	fake.ScheduleStub = stub
}

func (fake *FakeBuildScheduler) ScheduleArgsForCall(i int) (lager.Logger, *algorithm.VersionsDB, []db.Job, db.Resources, atc.VersionedResourceTypes) {
	fake.scheduleMutex.RLock()
	defer fake.scheduleMutex.RUnlock()
	argsForCall := fake.scheduleArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeBuildScheduler) ScheduleReturns(result1 map[string]time.Duration, result2 error) {
	fake.scheduleMutex.Lock()
	defer fake.scheduleMutex.Unlock()
	fake.ScheduleStub = nil
	fake.scheduleReturns = struct {
		result1 map[string]time.Duration
		result2 error
	}{result1, result2}
}

func (fake *FakeBuildScheduler) ScheduleReturnsOnCall(i int, result1 map[string]time.Duration, result2 error) {
	fake.scheduleMutex.Lock()
	defer fake.scheduleMutex.Unlock()
	fake.ScheduleStub = nil
	if fake.scheduleReturnsOnCall == nil {
		fake.scheduleReturnsOnCall = make(map[int]struct {
			result1 map[string]time.Duration
			result2 error
		})
	}
	fake.scheduleReturnsOnCall[i] = struct {
		result1 map[string]time.Duration
		result2 error
	}{result1, result2}
}

func (fake *FakeBuildScheduler) TriggerImmediately(arg1 lager.Logger, arg2 db.Job, arg3 db.Resources, arg4 atc.VersionedResourceTypes) (db.Build, scheduler.Waiter, error) {
	fake.triggerImmediatelyMutex.Lock()
	ret, specificReturn := fake.triggerImmediatelyReturnsOnCall[len(fake.triggerImmediatelyArgsForCall)]
	fake.triggerImmediatelyArgsForCall = append(fake.triggerImmediatelyArgsForCall, struct {
		arg1 lager.Logger
		arg2 db.Job
		arg3 db.Resources
		arg4 atc.VersionedResourceTypes
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("TriggerImmediately", []interface{}{arg1, arg2, arg3, arg4})
	fake.triggerImmediatelyMutex.Unlock()
	if fake.TriggerImmediatelyStub != nil {
		return fake.TriggerImmediatelyStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.triggerImmediatelyReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeBuildScheduler) TriggerImmediatelyCallCount() int {
	fake.triggerImmediatelyMutex.RLock()
	defer fake.triggerImmediatelyMutex.RUnlock()
	return len(fake.triggerImmediatelyArgsForCall)
}

func (fake *FakeBuildScheduler) TriggerImmediatelyCalls(stub func(lager.Logger, db.Job, db.Resources, atc.VersionedResourceTypes) (db.Build, scheduler.Waiter, error)) {
	fake.triggerImmediatelyMutex.Lock()
	defer fake.triggerImmediatelyMutex.Unlock()
	fake.TriggerImmediatelyStub = stub
}

func (fake *FakeBuildScheduler) TriggerImmediatelyArgsForCall(i int) (lager.Logger, db.Job, db.Resources, atc.VersionedResourceTypes) {
	fake.triggerImmediatelyMutex.RLock()
	defer fake.triggerImmediatelyMutex.RUnlock()
	argsForCall := fake.triggerImmediatelyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeBuildScheduler) TriggerImmediatelyReturns(result1 db.Build, result2 scheduler.Waiter, result3 error) {
	fake.triggerImmediatelyMutex.Lock()
	defer fake.triggerImmediatelyMutex.Unlock()
	fake.TriggerImmediatelyStub = nil
	fake.triggerImmediatelyReturns = struct {
		result1 db.Build
		result2 scheduler.Waiter
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildScheduler) TriggerImmediatelyReturnsOnCall(i int, result1 db.Build, result2 scheduler.Waiter, result3 error) {
	fake.triggerImmediatelyMutex.Lock()
	defer fake.triggerImmediatelyMutex.Unlock()
	fake.TriggerImmediatelyStub = nil
	if fake.triggerImmediatelyReturnsOnCall == nil {
		fake.triggerImmediatelyReturnsOnCall = make(map[int]struct {
			result1 db.Build
			result2 scheduler.Waiter
			result3 error
		})
	}
	fake.triggerImmediatelyReturnsOnCall[i] = struct {
		result1 db.Build
		result2 scheduler.Waiter
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildScheduler) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.saveNextInputMappingMutex.RLock()
	defer fake.saveNextInputMappingMutex.RUnlock()
	fake.scheduleMutex.RLock()
	defer fake.scheduleMutex.RUnlock()
	fake.triggerImmediatelyMutex.RLock()
	defer fake.triggerImmediatelyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBuildScheduler) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ scheduler.BuildScheduler = new(FakeBuildScheduler)

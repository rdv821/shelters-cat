// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/catService/internal/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// SheltersCatRepository is an autogenerated mock type for the SheltersCatRepository type
type SheltersCatRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *SheltersCatRepository) Create(_a0 context.Context, _a1 *model.Cat) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Cat) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *SheltersCatRepository) Delete(_a0 context.Context, _a1 uuid.UUID) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *SheltersCatRepository) Get(_a0 context.Context, _a1 uuid.UUID) (*model.Cat, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.Cat
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.Cat); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Cat)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *SheltersCatRepository) Update(_a0 context.Context, _a1 *model.Cat) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Cat) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

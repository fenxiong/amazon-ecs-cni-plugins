// Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aws/amazon-ecs-cni-plugins/pkg/cniipwrapper (interfaces: IP)

// Package mock_cniipwrapper is a generated GoMock package.
package mock_cniipwrapper

import (
	net "net"
	reflect "reflect"

	ns "github.com/containernetworking/cni/pkg/ns"
	gomock "github.com/golang/mock/gomock"
)

// MockIP is a mock of IP interface
type MockIP struct {
	ctrl     *gomock.Controller
	recorder *MockIPMockRecorder
}

// MockIPMockRecorder is the mock recorder for MockIP
type MockIPMockRecorder struct {
	mock *MockIP
}

// NewMockIP creates a new mock instance
func NewMockIP(ctrl *gomock.Controller) *MockIP {
	mock := &MockIP{ctrl: ctrl}
	mock.recorder = &MockIPMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIP) EXPECT() *MockIPMockRecorder {
	return m.recorder
}

// DelLinkByNameAddr mocks base method
func (m *MockIP) DelLinkByNameAddr(arg0 string, arg1 int) (*net.IPNet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelLinkByNameAddr", arg0, arg1)
	ret0, _ := ret[0].(*net.IPNet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DelLinkByNameAddr indicates an expected call of DelLinkByNameAddr
func (mr *MockIPMockRecorder) DelLinkByNameAddr(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelLinkByNameAddr", reflect.TypeOf((*MockIP)(nil).DelLinkByNameAddr), arg0, arg1)
}

// SetHWAddrByIP mocks base method
func (m *MockIP) SetHWAddrByIP(arg0 string, arg1, arg2 net.IP) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHWAddrByIP", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHWAddrByIP indicates an expected call of SetHWAddrByIP
func (mr *MockIPMockRecorder) SetHWAddrByIP(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHWAddrByIP", reflect.TypeOf((*MockIP)(nil).SetHWAddrByIP), arg0, arg1, arg2)
}

// SetupVeth mocks base method
func (m *MockIP) SetupVeth(arg0 string, arg1 int, arg2 ns.NetNS) (net.Interface, net.Interface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetupVeth", arg0, arg1, arg2)
	ret0, _ := ret[0].(net.Interface)
	ret1, _ := ret[1].(net.Interface)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SetupVeth indicates an expected call of SetupVeth
func (mr *MockIPMockRecorder) SetupVeth(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetupVeth", reflect.TypeOf((*MockIP)(nil).SetupVeth), arg0, arg1, arg2)
}

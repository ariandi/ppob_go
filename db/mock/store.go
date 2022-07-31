// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ariandi/ppob_go/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/ariandi/ppob_go/db/sqlc"
	dto "github.com/ariandi/ppob_go/dto"
	token "github.com/ariandi/ppob_go/token"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateCategory mocks base method.
func (m *MockStore) CreateCategory(arg0 context.Context, arg1 db.CreateCategoryParams) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCategory", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCategory indicates an expected call of CreateCategory.
func (mr *MockStoreMockRecorder) CreateCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCategory", reflect.TypeOf((*MockStore)(nil).CreateCategory), arg0, arg1)
}

// CreatePartner mocks base method.
func (m *MockStore) CreatePartner(arg0 context.Context, arg1 db.CreatePartnerParams) (db.Partner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePartner", arg0, arg1)
	ret0, _ := ret[0].(db.Partner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePartner indicates an expected call of CreatePartner.
func (mr *MockStoreMockRecorder) CreatePartner(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePartner", reflect.TypeOf((*MockStore)(nil).CreatePartner), arg0, arg1)
}

// CreateProduct mocks base method.
func (m *MockStore) CreateProduct(arg0 context.Context, arg1 db.CreateProductParams) (db.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", arg0, arg1)
	ret0, _ := ret[0].(db.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockStoreMockRecorder) CreateProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockStore)(nil).CreateProduct), arg0, arg1)
}

// CreateProvider mocks base method.
func (m *MockStore) CreateProvider(arg0 context.Context, arg1 db.CreateProviderParams) (db.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProvider", arg0, arg1)
	ret0, _ := ret[0].(db.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProvider indicates an expected call of CreateProvider.
func (mr *MockStoreMockRecorder) CreateProvider(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProvider", reflect.TypeOf((*MockStore)(nil).CreateProvider), arg0, arg1)
}

// CreateRole mocks base method.
func (m *MockStore) CreateRole(arg0 context.Context, arg1 db.CreateRoleParams) (db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRole", arg0, arg1)
	ret0, _ := ret[0].(db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRole indicates an expected call of CreateRole.
func (mr *MockStoreMockRecorder) CreateRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockStore)(nil).CreateRole), arg0, arg1)
}

// CreateRoleUser mocks base method.
func (m *MockStore) CreateRoleUser(arg0 context.Context, arg1 db.CreateRoleUserParams) (db.RoleUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRoleUser", arg0, arg1)
	ret0, _ := ret[0].(db.RoleUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRoleUser indicates an expected call of CreateRoleUser.
func (mr *MockStoreMockRecorder) CreateRoleUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRoleUser", reflect.TypeOf((*MockStore)(nil).CreateRoleUser), arg0, arg1)
}

// CreateTransaction mocks base method.
func (m *MockStore) CreateTransaction(arg0 context.Context, arg1 db.CreateTransactionParams) (db.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", arg0, arg1)
	ret0, _ := ret[0].(db.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockStoreMockRecorder) CreateTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockStore)(nil).CreateTransaction), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// CreateUserTx mocks base method.
func (m *MockStore) CreateUserTx(arg0 context.Context, arg1 db.CreateUserParams, arg2 *token.Payload, arg3 int64) (dto.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserTx", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(dto.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserTx indicates an expected call of CreateUserTx.
func (mr *MockStoreMockRecorder) CreateUserTx(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserTx", reflect.TypeOf((*MockStore)(nil).CreateUserTx), arg0, arg1, arg2, arg3)
}

// DeleteCategories mocks base method.
func (m *MockStore) DeleteCategories(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategories", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategories indicates an expected call of DeleteCategories.
func (mr *MockStoreMockRecorder) DeleteCategories(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategories", reflect.TypeOf((*MockStore)(nil).DeleteCategories), arg0, arg1)
}

// DeletePartner mocks base method.
func (m *MockStore) DeletePartner(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePartner", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePartner indicates an expected call of DeletePartner.
func (mr *MockStoreMockRecorder) DeletePartner(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePartner", reflect.TypeOf((*MockStore)(nil).DeletePartner), arg0, arg1)
}

// DeleteProduct mocks base method.
func (m *MockStore) DeleteProduct(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct.
func (mr *MockStoreMockRecorder) DeleteProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockStore)(nil).DeleteProduct), arg0, arg1)
}

// DeleteProvider mocks base method.
func (m *MockStore) DeleteProvider(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProvider", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProvider indicates an expected call of DeleteProvider.
func (mr *MockStoreMockRecorder) DeleteProvider(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProvider", reflect.TypeOf((*MockStore)(nil).DeleteProvider), arg0, arg1)
}

// DeleteRole mocks base method.
func (m *MockStore) DeleteRole(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRole", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRole indicates an expected call of DeleteRole.
func (mr *MockStoreMockRecorder) DeleteRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRole", reflect.TypeOf((*MockStore)(nil).DeleteRole), arg0, arg1)
}

// DeleteRoleUser mocks base method.
func (m *MockStore) DeleteRoleUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRoleUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRoleUser indicates an expected call of DeleteRoleUser.
func (mr *MockStoreMockRecorder) DeleteRoleUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRoleUser", reflect.TypeOf((*MockStore)(nil).DeleteRoleUser), arg0, arg1)
}

// DeleteTransaction mocks base method.
func (m *MockStore) DeleteTransaction(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTransaction indicates an expected call of DeleteTransaction.
func (mr *MockStoreMockRecorder) DeleteTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransaction", reflect.TypeOf((*MockStore)(nil).DeleteTransaction), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1)
}

// GetCategory mocks base method.
func (m *MockStore) GetCategory(arg0 context.Context, arg1 int64) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategory", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategory indicates an expected call of GetCategory.
func (mr *MockStoreMockRecorder) GetCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategory", reflect.TypeOf((*MockStore)(nil).GetCategory), arg0, arg1)
}

// GetPartner mocks base method.
func (m *MockStore) GetPartner(arg0 context.Context, arg1 int64) (db.Partner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPartner", arg0, arg1)
	ret0, _ := ret[0].(db.Partner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPartner indicates an expected call of GetPartner.
func (mr *MockStoreMockRecorder) GetPartner(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPartner", reflect.TypeOf((*MockStore)(nil).GetPartner), arg0, arg1)
}

// GetProduct mocks base method.
func (m *MockStore) GetProduct(arg0 context.Context, arg1 int64) (db.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", arg0, arg1)
	ret0, _ := ret[0].(db.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockStoreMockRecorder) GetProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockStore)(nil).GetProduct), arg0, arg1)
}

// GetProvider mocks base method.
func (m *MockStore) GetProvider(arg0 context.Context, arg1 int64) (db.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProvider", arg0, arg1)
	ret0, _ := ret[0].(db.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProvider indicates an expected call of GetProvider.
func (mr *MockStoreMockRecorder) GetProvider(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProvider", reflect.TypeOf((*MockStore)(nil).GetProvider), arg0, arg1)
}

// GetRole mocks base method.
func (m *MockStore) GetRole(arg0 context.Context, arg1 int64) (db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", arg0, arg1)
	ret0, _ := ret[0].(db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole.
func (mr *MockStoreMockRecorder) GetRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockStore)(nil).GetRole), arg0, arg1)
}

// GetRoleUserByID mocks base method.
func (m *MockStore) GetRoleUserByID(arg0 context.Context, arg1 int64) (db.RoleUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoleUserByID", arg0, arg1)
	ret0, _ := ret[0].(db.RoleUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoleUserByID indicates an expected call of GetRoleUserByID.
func (mr *MockStoreMockRecorder) GetRoleUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoleUserByID", reflect.TypeOf((*MockStore)(nil).GetRoleUserByID), arg0, arg1)
}

// GetRoleUserByRoleID mocks base method.
func (m *MockStore) GetRoleUserByRoleID(arg0 context.Context, arg1 db.GetRoleUserByRoleIDParams) ([]db.RoleUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoleUserByRoleID", arg0, arg1)
	ret0, _ := ret[0].([]db.RoleUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoleUserByRoleID indicates an expected call of GetRoleUserByRoleID.
func (mr *MockStoreMockRecorder) GetRoleUserByRoleID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoleUserByRoleID", reflect.TypeOf((*MockStore)(nil).GetRoleUserByRoleID), arg0, arg1)
}

// GetRoleUserByUserID mocks base method.
func (m *MockStore) GetRoleUserByUserID(arg0 context.Context, arg1 db.GetRoleUserByUserIDParams) ([]db.RoleUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoleUserByUserID", arg0, arg1)
	ret0, _ := ret[0].([]db.RoleUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoleUserByUserID indicates an expected call of GetRoleUserByUserID.
func (mr *MockStoreMockRecorder) GetRoleUserByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoleUserByUserID", reflect.TypeOf((*MockStore)(nil).GetRoleUserByUserID), arg0, arg1)
}

// GetTransaction mocks base method.
func (m *MockStore) GetTransaction(arg0 context.Context, arg1 int64) (db.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", arg0, arg1)
	ret0, _ := ret[0].(db.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockStoreMockRecorder) GetTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockStore)(nil).GetTransaction), arg0, arg1)
}

// GetTransactionByTxID mocks base method.
func (m *MockStore) GetTransactionByTxID(arg0 context.Context, arg1 string) (db.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionByTxID", arg0, arg1)
	ret0, _ := ret[0].(db.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionByTxID indicates an expected call of GetTransactionByTxID.
func (mr *MockStoreMockRecorder) GetTransactionByTxID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionByTxID", reflect.TypeOf((*MockStore)(nil).GetTransactionByTxID), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 int64) (db.GetUserRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.GetUserRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// GetUserByUsername mocks base method.
func (m *MockStore) GetUserByUsername(arg0 context.Context, arg1 string) (db.GetUserByUsernameRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", arg0, arg1)
	ret0, _ := ret[0].(db.GetUserByUsernameRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockStoreMockRecorder) GetUserByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockStore)(nil).GetUserByUsername), arg0, arg1)
}

// ListCategory mocks base method.
func (m *MockStore) ListCategory(arg0 context.Context, arg1 db.ListCategoryParams) ([]db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCategory", arg0, arg1)
	ret0, _ := ret[0].([]db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCategory indicates an expected call of ListCategory.
func (mr *MockStoreMockRecorder) ListCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCategory", reflect.TypeOf((*MockStore)(nil).ListCategory), arg0, arg1)
}

// ListPartner mocks base method.
func (m *MockStore) ListPartner(arg0 context.Context, arg1 db.ListPartnerParams) ([]db.Partner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPartner", arg0, arg1)
	ret0, _ := ret[0].([]db.Partner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPartner indicates an expected call of ListPartner.
func (mr *MockStoreMockRecorder) ListPartner(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPartner", reflect.TypeOf((*MockStore)(nil).ListPartner), arg0, arg1)
}

// ListProduct mocks base method.
func (m *MockStore) ListProduct(arg0 context.Context, arg1 db.ListProductParams) ([]db.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProduct", arg0, arg1)
	ret0, _ := ret[0].([]db.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProduct indicates an expected call of ListProduct.
func (mr *MockStoreMockRecorder) ListProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProduct", reflect.TypeOf((*MockStore)(nil).ListProduct), arg0, arg1)
}

// ListProductByCatID mocks base method.
func (m *MockStore) ListProductByCatID(arg0 context.Context, arg1 db.ListProductByCatIDParams) ([]db.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProductByCatID", arg0, arg1)
	ret0, _ := ret[0].([]db.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProductByCatID indicates an expected call of ListProductByCatID.
func (mr *MockStoreMockRecorder) ListProductByCatID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProductByCatID", reflect.TypeOf((*MockStore)(nil).ListProductByCatID), arg0, arg1)
}

// ListProvider mocks base method.
func (m *MockStore) ListProvider(arg0 context.Context, arg1 db.ListProviderParams) ([]db.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProvider", arg0, arg1)
	ret0, _ := ret[0].([]db.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProvider indicates an expected call of ListProvider.
func (mr *MockStoreMockRecorder) ListProvider(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProvider", reflect.TypeOf((*MockStore)(nil).ListProvider), arg0, arg1)
}

// ListRole mocks base method.
func (m *MockStore) ListRole(arg0 context.Context, arg1 db.ListRoleParams) ([]db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRole", arg0, arg1)
	ret0, _ := ret[0].([]db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRole indicates an expected call of ListRole.
func (mr *MockStoreMockRecorder) ListRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRole", reflect.TypeOf((*MockStore)(nil).ListRole), arg0, arg1)
}

// ListRoleUser mocks base method.
func (m *MockStore) ListRoleUser(arg0 context.Context, arg1 db.ListRoleUserParams) ([]db.RoleUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRoleUser", arg0, arg1)
	ret0, _ := ret[0].([]db.RoleUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRoleUser indicates an expected call of ListRoleUser.
func (mr *MockStoreMockRecorder) ListRoleUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRoleUser", reflect.TypeOf((*MockStore)(nil).ListRoleUser), arg0, arg1)
}

// ListTransaction mocks base method.
func (m *MockStore) ListTransaction(arg0 context.Context, arg1 db.ListTransactionParams) ([]db.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTransaction", arg0, arg1)
	ret0, _ := ret[0].([]db.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTransaction indicates an expected call of ListTransaction.
func (mr *MockStoreMockRecorder) ListTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTransaction", reflect.TypeOf((*MockStore)(nil).ListTransaction), arg0, arg1)
}

// ListUser mocks base method.
func (m *MockStore) ListUser(arg0 context.Context, arg1 db.ListUserParams) ([]db.ListUserRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUser", arg0, arg1)
	ret0, _ := ret[0].([]db.ListUserRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUser indicates an expected call of ListUser.
func (mr *MockStoreMockRecorder) ListUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUser", reflect.TypeOf((*MockStore)(nil).ListUser), arg0, arg1)
}

// UpdateCategory mocks base method.
func (m *MockStore) UpdateCategory(arg0 context.Context, arg1 db.UpdateCategoryParams) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCategory", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCategory indicates an expected call of UpdateCategory.
func (mr *MockStoreMockRecorder) UpdateCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCategory", reflect.TypeOf((*MockStore)(nil).UpdateCategory), arg0, arg1)
}

// UpdateInactiveCategory mocks base method.
func (m *MockStore) UpdateInactiveCategory(arg0 context.Context, arg1 db.UpdateInactiveCategoryParams) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInactiveCategory", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInactiveCategory indicates an expected call of UpdateInactiveCategory.
func (mr *MockStoreMockRecorder) UpdateInactiveCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInactiveCategory", reflect.TypeOf((*MockStore)(nil).UpdateInactiveCategory), arg0, arg1)
}

// UpdateInactivePartner mocks base method.
func (m *MockStore) UpdateInactivePartner(arg0 context.Context, arg1 db.UpdateInactivePartnerParams) (db.Partner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInactivePartner", arg0, arg1)
	ret0, _ := ret[0].(db.Partner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInactivePartner indicates an expected call of UpdateInactivePartner.
func (mr *MockStoreMockRecorder) UpdateInactivePartner(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInactivePartner", reflect.TypeOf((*MockStore)(nil).UpdateInactivePartner), arg0, arg1)
}

// UpdateInactiveProduct mocks base method.
func (m *MockStore) UpdateInactiveProduct(arg0 context.Context, arg1 db.UpdateInactiveProductParams) (db.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInactiveProduct", arg0, arg1)
	ret0, _ := ret[0].(db.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInactiveProduct indicates an expected call of UpdateInactiveProduct.
func (mr *MockStoreMockRecorder) UpdateInactiveProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInactiveProduct", reflect.TypeOf((*MockStore)(nil).UpdateInactiveProduct), arg0, arg1)
}

// UpdateInactiveProvider mocks base method.
func (m *MockStore) UpdateInactiveProvider(arg0 context.Context, arg1 db.UpdateInactiveProviderParams) (db.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInactiveProvider", arg0, arg1)
	ret0, _ := ret[0].(db.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInactiveProvider indicates an expected call of UpdateInactiveProvider.
func (mr *MockStoreMockRecorder) UpdateInactiveProvider(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInactiveProvider", reflect.TypeOf((*MockStore)(nil).UpdateInactiveProvider), arg0, arg1)
}

// UpdateInactiveRole mocks base method.
func (m *MockStore) UpdateInactiveRole(arg0 context.Context, arg1 db.UpdateInactiveRoleParams) (db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInactiveRole", arg0, arg1)
	ret0, _ := ret[0].(db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInactiveRole indicates an expected call of UpdateInactiveRole.
func (mr *MockStoreMockRecorder) UpdateInactiveRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInactiveRole", reflect.TypeOf((*MockStore)(nil).UpdateInactiveRole), arg0, arg1)
}

// UpdateInactiveRoleUser mocks base method.
func (m *MockStore) UpdateInactiveRoleUser(arg0 context.Context, arg1 db.UpdateInactiveRoleUserParams) (db.RoleUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInactiveRoleUser", arg0, arg1)
	ret0, _ := ret[0].(db.RoleUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInactiveRoleUser indicates an expected call of UpdateInactiveRoleUser.
func (mr *MockStoreMockRecorder) UpdateInactiveRoleUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInactiveRoleUser", reflect.TypeOf((*MockStore)(nil).UpdateInactiveRoleUser), arg0, arg1)
}

// UpdateInactiveTransaction mocks base method.
func (m *MockStore) UpdateInactiveTransaction(arg0 context.Context, arg1 db.UpdateInactiveTransactionParams) (db.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInactiveTransaction", arg0, arg1)
	ret0, _ := ret[0].(db.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInactiveTransaction indicates an expected call of UpdateInactiveTransaction.
func (mr *MockStoreMockRecorder) UpdateInactiveTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInactiveTransaction", reflect.TypeOf((*MockStore)(nil).UpdateInactiveTransaction), arg0, arg1)
}

// UpdateInactiveUser mocks base method.
func (m *MockStore) UpdateInactiveUser(arg0 context.Context, arg1 db.UpdateInactiveUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInactiveUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInactiveUser indicates an expected call of UpdateInactiveUser.
func (mr *MockStoreMockRecorder) UpdateInactiveUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInactiveUser", reflect.TypeOf((*MockStore)(nil).UpdateInactiveUser), arg0, arg1)
}

// UpdatePartner mocks base method.
func (m *MockStore) UpdatePartner(arg0 context.Context, arg1 db.UpdatePartnerParams) (db.Partner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePartner", arg0, arg1)
	ret0, _ := ret[0].(db.Partner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePartner indicates an expected call of UpdatePartner.
func (mr *MockStoreMockRecorder) UpdatePartner(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePartner", reflect.TypeOf((*MockStore)(nil).UpdatePartner), arg0, arg1)
}

// UpdateProduct mocks base method.
func (m *MockStore) UpdateProduct(arg0 context.Context, arg1 db.UpdateProductParams) (db.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", arg0, arg1)
	ret0, _ := ret[0].(db.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockStoreMockRecorder) UpdateProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockStore)(nil).UpdateProduct), arg0, arg1)
}

// UpdateProvider mocks base method.
func (m *MockStore) UpdateProvider(arg0 context.Context, arg1 db.UpdateProviderParams) (db.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProvider", arg0, arg1)
	ret0, _ := ret[0].(db.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProvider indicates an expected call of UpdateProvider.
func (mr *MockStoreMockRecorder) UpdateProvider(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProvider", reflect.TypeOf((*MockStore)(nil).UpdateProvider), arg0, arg1)
}

// UpdateRole mocks base method.
func (m *MockStore) UpdateRole(arg0 context.Context, arg1 db.UpdateRoleParams) (db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRole", arg0, arg1)
	ret0, _ := ret[0].(db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRole indicates an expected call of UpdateRole.
func (mr *MockStoreMockRecorder) UpdateRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRole", reflect.TypeOf((*MockStore)(nil).UpdateRole), arg0, arg1)
}

// UpdateRoleUser mocks base method.
func (m *MockStore) UpdateRoleUser(arg0 context.Context, arg1 db.UpdateRoleUserParams) (db.RoleUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRoleUser", arg0, arg1)
	ret0, _ := ret[0].(db.RoleUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRoleUser indicates an expected call of UpdateRoleUser.
func (mr *MockStoreMockRecorder) UpdateRoleUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRoleUser", reflect.TypeOf((*MockStore)(nil).UpdateRoleUser), arg0, arg1)
}

// UpdateTransaction mocks base method.
func (m *MockStore) UpdateTransaction(arg0 context.Context, arg1 db.UpdateTransactionParams) (db.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", arg0, arg1)
	ret0, _ := ret[0].(db.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockStoreMockRecorder) UpdateTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockStore)(nil).UpdateTransaction), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockStore) UpdateUser(arg0 context.Context, arg1 db.UpdateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStoreMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStore)(nil).UpdateUser), arg0, arg1)
}

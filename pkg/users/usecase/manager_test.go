package usecase

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang-stepik-2022q1/reditclone/pkg/errors"
	"golang-stepik-2022q1/reditclone/pkg/users"
	"golang-stepik-2022q1/reditclone/pkg/users/repo"
	"testing"
)

func TestManager_GetByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	name := "John"
	user := &users.User{1, name, "password"}

	for _, tt := range [...]struct {
		name    string
		setup   func(st *repo.MockRepo)
		want    *users.User
		wantErr error
	}{
		{
			name: "Ok",
			setup: func(st *repo.MockRepo) {
				st.EXPECT().GetByName(name).Return(user, nil)
			},
			want: user,
		},
		{
			name: "Repo error",
			setup: func(st *repo.MockRepo) {
				st.EXPECT().GetByName(name).Return(nil, fmt.Errorf("Unexpected error"))
			},
			want:    nil,
			wantErr: errors.InternalError{"Unexpected error"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			st := repo.NewMockRepo(ctrl)
			manager := NewManager(st)
			tt.setup(st)

			item, err := manager.GetByName(context.Background(), name)

			assert.Equal(t, tt.want, item)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestManager_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userData := &users.UserIn{
		Name:     "John",
		Password: "my-pass",
	}
	userId := 123
	expectedUser := &users.User{
		Id:   userId,
		Name: userData.Name,
	}

	for _, tt := range [...]struct {
		name    string
		setup   func(st *repo.MockRepo)
		want    *users.User
		wantErr error
	}{
		{
			name: "Ok",
			setup: func(st *repo.MockRepo) {
				gomock.InOrder(
					st.EXPECT().GetByName(userData.Name).Return(nil, nil),
					st.EXPECT().Add(gomock.AssignableToTypeOf(&users.User{})).Return(int64(userId), nil),
				)
			},
			want: expectedUser,
		},
		{
			name: "Repo add error",
			setup: func(st *repo.MockRepo) {
				gomock.InOrder(
					st.EXPECT().GetByName(userData.Name).Return(nil, nil),
					st.EXPECT().Add(gomock.AssignableToTypeOf(&users.User{})).Return(int64(0), fmt.Errorf("Unexpected error")),
				)
			},
			want:    nil,
			wantErr: errors.InternalError{"Unexpected error"},
		},
		{
			name: "Repo find user error",
			setup: func(st *repo.MockRepo) {
				st.EXPECT().GetByName(userData.Name).Return(nil, fmt.Errorf("Unexpected error"))
			},
			want:    nil,
			wantErr: errors.InternalError{"Unexpected error"},
		},
		{
			name: "UserId exists",
			setup: func(st *repo.MockRepo) {
				st.EXPECT().GetByName(userData.Name).Return(&users.User{}, nil)
			},
			want:    nil,
			wantErr: UserExistsError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			st := repo.NewMockRepo(ctrl)
			manager := NewManager(st)
			tt.setup(st)

			item, err := manager.Create(context.Background(), userData)
			assert.Equal(t, tt.wantErr, err)
			if tt.want != nil {
				assert.Equal(t, tt.want.Id, item.Id)
				assert.Equal(t, tt.want.Name, item.Name)
				if !CheckPassword(item.PassHash, userData.Password) {
					t.Errorf("Wrong pass hash")
				}
			} else {
				assert.Equal(t, tt.want, item)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	pass := "SuperSecret"
	hashPass, _ := HashPass(pass)

	for _, tt := range [...]struct {
		name     string
		hashPass string
		pass     string
		expected bool
	}{
		{
			name:     "Ok",
			hashPass: hashPass,
			pass:     pass,
			expected: true,
		},
		{
			name:     "Wrong pass",
			hashPass: hashPass,
			pass:     "OtherPassword",
			expected: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CheckPassword(tt.hashPass, tt.pass))
		})
	}
}

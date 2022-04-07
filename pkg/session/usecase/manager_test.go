package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang-stepik-2022q1/reditclone/pkg/session"
	"golang-stepik-2022q1/reditclone/pkg/session/repo"
	"golang-stepik-2022q1/reditclone/pkg/users"
	"testing"
)

func TestManager_IssueToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := repo.NewMockRepo(ctrl)
	manager := NewManager(st)

	ctx := context.Background()
	user := &users.User{Id: 123, Name: "John"}

	st.EXPECT().Set(ctx, gomock.AssignableToTypeOf(session.SessionId(""))).Return(nil)

	token, err := manager.IssueToken(ctx, user)
	assert.Equal(t, nil, err)
	sess, _ := loadSession(token)
	assert.Equal(t, sess.User, session.UserClaims{Id: user.Id, Username: user.Name})

	//for _, tt := range [...]struct{
	//	name string
	//	setup func()
	//	expected string
	//	expectedErr error
	//} {
	//	{
	//		name: "Ok",
	//		setup: func() {
	//			st.EXPECT().
	//		},
	//	},
	//}
}

//func (m *Manager) IssueToken(ctx context.Context, u *users.User) (string, error) {
//	sess := &session.Session{
//		Id:   session.SessionId(uuid.New().String()),
//		User: session.UserClaims{u.Name, u.Id},
//		Iat:  time.Now().Unix(),
//		Exp:  expDate().Unix(),
//	}
//	err := m.repo.Set(ctx, sess.Id)
//	if err != nil {
//		log.Clog(ctx).Error("Error during session creation", log.Fields{"error": err.Error()})
//		return "", errors.InternalError{"Error during session creation"}
//	}
//	tokenString, err := generateToken(sess)
//	if err != nil {
//		detail := "Error during jwt token generation"
//		log.Clog(ctx).Error(detail, log.Fields{"error": err.Error()})
//		return "", errors.InternalError{detail}
//	}
//	return tokenString, nil
//}

//func TestManager_GetByName(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	name := "John"
//	user := &users.User{1, name, "password"}
//
//	for _, tt := range [...]struct {
//		name    string
//		setup   func(st *repo.MockRepo)
//		want    *users.User
//		wantErr error
//	}{
//		{
//			name: "Ok",
//			setup: func(st *repo.MockRepo) {
//				st.EXPECT().GetByName(name).Return(user, nil)
//			},
//			want: user,
//		},
//		{
//			name: "Repo error",
//			setup: func(st *repo.MockRepo) {
//				st.EXPECT().GetByName(name).Return(nil, fmt.Errorf("Unexpected error"))
//			},
//			want:    nil,
//			wantErr: errors.InternalError{"Unexpected error"},
//		},
//	} {
//		t.Run(tt.name, func(t *testing.T) {
//			st := repo.NewMockRepo(ctrl)
//			manager := NewManager(st)
//			tt.setup(st)
//
//			item, err := manager.GetByName(context.Background(), name)
//
//			assert.Equal(t, tt.want, item)
//			assert.Equal(t, tt.wantErr, err)
//		})
//	}
//}

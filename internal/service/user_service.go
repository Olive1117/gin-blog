package service

import (
	"context"
	"time"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/Olive1117/gin-blog/pkg/utils"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type userService struct {
	Repo             repository.UserRepo
	ArticleRepo      repository.ArticleRepo
	RefreshTokenRepo repository.RefreshTokenRepo
	jwt              model.JWTHandler
}

func NewUserService(repo repository.UserRepo, articleRepo repository.ArticleRepo, refreshtokenRepo repository.RefreshTokenRepo, jwt model.JWTHandler) UserService {
	return &userService{
		Repo:             repo,
		ArticleRepo:      articleRepo,
		RefreshTokenRepo: refreshtokenRepo,
		jwt:              jwt,
	}
}

func (ts *userService) Create(c context.Context, user *model.User) error {
	if users, err := ts.Repo.FindByUniqueKeys(c, user.Username, user.Email); err != nil || len(users) > 0 {
		// 唯一性信息冲突
		for _, u := range users {
			if u.Username == user.Username {
				return errs.ErrExistUsername
			}
			if u.Email == user.Email {
				return errs.ErrExistEmail
			}
		}
		return err
	}
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		logger.ErrorContext(c, "密码加密失败", logger.Err(err))
		return errs.ErrRegisterFail
	}
	user.Password = hashPassword
	return ts.Repo.Create(c, user)
}
func (ts *userService) Delete(c context.Context, id int64) error {
	return ts.Repo.Delete(c, id)
}
func (ts *userService) Get(c context.Context, id int64) (model.User, error) {
	user, err := ts.Repo.FindById(c, id)
	if err != nil {
		return model.User{}, err
	}
	postCount, err := ts.ArticleRepo.CountArticleByUserID(c, id)
	if err != nil {
		return model.User{}, err
	}
	user.PostCount = cast.ToInt(postCount)
	return user, nil
}
func (ts *userService) List(c context.Context, que model.PageQuery, filter *model.User) (model.PageResult[model.User], error) {
	return ts.Repo.List(c, que, filter)
}
func (ts *userService) Update(c context.Context, user *model.User, id int64) error {
	user.Password = "" // 不允许更新密码
	users, err := ts.Repo.FindByUniqueKeys(c, user.Username, user.Email)
	if err != nil {
		return err
	}
	for _, u := range users {
		// 如果不是用户自己的信息，发生冲突
		if u.ID != id {
			if u.Username == user.Username {
				return errs.ErrExistUsername
			}
			if u.Email == user.Email {
				return errs.ErrExistEmail
			}
		}
	}
	return ts.Repo.Update(c, id, user)
}
func (ts *userService) Login(c context.Context, req model.LoginRequest, userInfo model.RefreshTokens) (model.AuthResponse, model.RefreshCookie, error) {
	logger.DebugContext(c, "登录业务代码")
	var Authres model.AuthResponse
	var Refres model.RefreshCookie
	user, err := ts.Repo.GetByUsername(c, req.Username)
	if err != nil {
		return Authres, Refres, err
	}
	if ok := utils.CheckPassword(req.Password, user.Password); !ok {
		return Authres, Refres, errs.ErrAuth
	}
	id := uuid.NewString()
	// TODO 等数据库改成Role列表后需要修改
	accessToken, accessExpiresAt, err := ts.jwt.GenerateAccessToken(cast.ToString(user.ID), []string{user.Role})
	refreshToken, refreshExpiresAt, err := ts.jwt.GenerateRefreshToken(cast.ToString(user.ID), []string{user.Role}, id)
	if err != nil {
		logger.WarnContext(c, errs.ErrAuthToken.Message, logger.Err(err))
		return Authres, Refres, errs.ErrAuthToken
	}
	userInfo.UserID = user.ID
	userInfo.Jti = id
	userInfo.Expires_at = &refreshExpiresAt
	if ts.RefreshTokenRepo.RevokeByUser(c, user.ID) != nil {
		return Authres, Refres, err
	}
	if ts.RefreshTokenRepo.Create(c, userInfo) != nil {
		return Authres, Refres, err
	}
	Authres.AccessToken = accessToken
	Authres.ExpiresAt = accessExpiresAt
	Authres.TokenType = "Bearer"
	Refres.RefreshToken = refreshToken
	Refres.ExpiresAt = refreshExpiresAt
	return Authres, Refres, nil
}
func (ts *userService) ChangePassword(c context.Context, username string, oldPassword string, newPassword string) error {
	user, err := ts.Repo.GetByUsername(c, username)
	user.Password = ""
	if err != nil {
		return err
	}
	if ok := utils.CheckPassword(oldPassword, user.Password); !ok {
		return errs.ErrPasswordIncorrect
	}
	hashPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		logger.ErrorContext(c, "密码加密失败", logger.Err(err))
		return errs.ErrRegisterFail
	}
	user.Password = hashPassword
	return ts.Repo.Update(c, user.ID, user)
}
func (ts *userService) RefreshToken(c context.Context, tokenString string) (model.AuthResponse, error) {
	var Authres model.AuthResponse
	queToken, err := ts.jwt.ParseToken(tokenString)
	if err != nil {
		return Authres, err
	}
	refreshToken, err := ts.RefreshTokenRepo.GetByJTI(c, queToken.ID)
	if err != nil {
		return Authres, err
	}
	if refreshToken.Revoked_at != nil || refreshToken.Expires_at.Before(time.Now()) {
		return Authres, errs.ErrAuth
	}
	accessToken, accessExpiresAt, err := ts.jwt.GenerateAccessToken(queToken.ID, queToken.Roles)
	Authres.AccessToken = accessToken
	Authres.ExpiresAt = accessExpiresAt
	Authres.TokenType = "Bearer"
	return Authres, nil
}

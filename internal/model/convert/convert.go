package convert

import (
	"time"

	"github.com/Olive1117/gin-blog/internal/model"
)

func MapSlice[T any, R any](src []T, fn func(*T) *R) []R {
	if len(src) == 0 {
		return nil
	}
	res := make([]R, 0, len(src))
	for i := range src {
		if r := fn(&src[i]); r != nil {
			res = append(res, *r)
		}
	}
	return res
}

func UserFromDTO(userVO *model.UserDTO) *model.User {
	if userVO == nil {
		return nil
	}
	res := &model.User{
		Email:    userVO.Email,
		Nickname: userVO.Nickname,
		Avatar:   userVO.Avatar,
		Banner:   userVO.Banner,
		Bio:      userVO.Bio,
		Location: userVO.Location,
		Website:  userVO.Website,
	}
	if b, err := time.ParseInLocation("2006-01-02", userVO.Birthdate, time.Local); err == nil {
		res.Birthdate = b
	}
	return res
}
func UserToVO(user *model.User) *model.UserVO {
	if user == nil {
		return nil
	}
	return &model.UserVO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Banner:      user.Banner,
		Bio:         user.Bio,
		Location:    user.Location,
		Website:     user.Website,
		Birthdate:   user.Birthdate,
		PostCount:   user.PostCount,
		FriendCount: user.FriendCount,
		Role:        user.Role,
		State:       user.State,
		CreatedAt:   user.CreatedAt,
	}
}

func ArticleFromDTO(articleVO *model.ArticleDTO) *model.Article {
	if articleVO == nil {
		return nil
	}
	res := &model.Article{
		Title:   articleVO.Title,
		Desc:    articleVO.Desc,
		Content: articleVO.Content,
		State:   articleVO.State,
		Slug:    articleVO.Slug,
		Category: model.Category{
			Name: articleVO.CategoryName,
		},
		ImageCount: articleVO.ImageCount,
	}
	var taglist = make([]model.Tag, 0, len(articleVO.TagNames))
	for _, tag := range articleVO.TagNames {
		taglist = append(taglist, model.Tag{Name: tag})
	}
	res.Tags = taglist
	return res
}

func ArticleToVO(article *model.Article) *model.ArticleVO {
	if article == nil {
		return nil
	}
	res := &model.ArticleVO{
		ID:           article.ID,
		Title:        article.Title,
		Desc:         article.Desc,
		Content:      article.Content,
		State:        article.State,
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
		ShortID:      article.ShortID,
		Slug:         article.Slug,
		CategoryName: article.Category.Name,
		WordCount:    article.WordCount,
		ImageCount:   article.ImageCount,
	}
	var taglist = make([]string, 0, len(article.Tags))
	for _, tag := range article.Tags {
		taglist = append(taglist, tag.Name)
	}
	res.TagNames = taglist
	return res
}

func ArticleFromQuery(articleQuery *model.ArticleQuery) *model.Article {
	if articleQuery == nil {
		return nil
	}
	res := &model.Article{
		Category: model.Category{Name: articleQuery.CategoryName},
		Title:    articleQuery.Title,
		State:    articleQuery.State,
	}
	var taglist = make([]model.Tag, 0, len(articleQuery.TagNames))
	for _, tag := range articleQuery.TagNames {
		taglist = append(taglist, model.Tag{Name: tag})
	}
	res.Tags = taglist
	return res
}

func CategoryToVO(category *model.Category) *model.CategoryVO {
	if category == nil {
		return nil
	}
	return &model.CategoryVO{
		ID:    category.ID,
		Name:  category.Name,
		State: category.State,
	}
}

func TagToVO(tag *model.Tag) *model.TagVO {
	if tag == nil {
		return nil
	}
	return &model.TagVO{
		ID:    tag.ID,
		Name:  tag.Name,
		State: tag.State,
	}
}

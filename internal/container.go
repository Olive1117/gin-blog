package internal

import (
	"github.com/Olive1117/gin-blog/internal/handler"
	"github.com/Olive1117/gin-blog/internal/middleware"
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/internal/router"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitContainer(engine *gin.Engine, jwt model.JWTHandler, db *gorm.DB, tx model.TransactionManager) {
	middlewareContainer := &middleware.MiddlewareContainer{
		Jwt:         middleware.JwtAuth(jwt),
		Logger:      middleware.GinLogger(),
		GinRecovery: middleware.GinRecovery(true),
	}

	articleRepo := repository.NewArticleRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	tagRepo := repository.NewTagRepo(db)
	userRepo := repository.NewUserRepo(db)
	friendLinkRepo := repository.NewFriendLinkRepo(db)
	refreshtokenRepo := repository.NewreFreshTokenRepo(db)

	articleService := service.NewArticleService(articleRepo, tx, tagRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	usrService := service.NewUserService(userRepo, articleRepo, refreshtokenRepo, jwt)
	friendLinkService := service.NewFriendLinkService(friendLinkRepo)

	articleHandler := handler.NewArticleHandler(articleService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	userHandler := handler.NewUserHandler(usrService)
	friendLinkHandler := handler.NewFriendlinkHandler(friendLinkService)

	handlerContainer := &handler.HandlerContainer{
		Article:    articleHandler,
		Category:   categoryHandler,
		Tag:        tagHandler,
		User:       userHandler,
		FriendLink: friendLinkHandler,
	}
	router.InitRouter(engine, handlerContainer, middlewareContainer)
}

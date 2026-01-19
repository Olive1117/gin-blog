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
	loginRepo := repository.NewLoginRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	tagRepo := repository.NewTagRepo(db)

	articleService := service.NewArticleService(articleRepo, tx, tagRepo, categoryRepo)
	loginService := service.NewLoginService(loginRepo, jwt)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)

	articleHandler := handler.NewArticleHandler(articleService)
	loginHandler := handler.NewLoginHandler(loginService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)

	handlerContainer := &handler.HandlerContainer{
		Login:    loginHandler,
		Article:  articleHandler,
		Category: categoryHandler,
		Tag:      tagHandler,
	}
	router.InitRouter(engine, handlerContainer, middlewareContainer)
}

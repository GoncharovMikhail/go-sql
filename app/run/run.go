package run

import (
	builtInSql "database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"log"
	"sql/delivery/http/user"
	"sql/internal/db"
	"sql/pkg/db/user/impl/sql"
	userService "sql/pkg/service/user"
	"sql/pkg/service/user/impl"
)

func Run() error {
	pg, err := db.GetDb(&pgx.ConnConfig{})
	if err != nil {
		log.Panicf("")
	}
	defer closeDb(pg)
	ur := sql.PostgresUserRepository{
		Db: pg,
	}
	var us userService.UserService = &impl.UserServiceImpl{
		Ur: &ur,
	}
	engine := gin.Default()
	uc := user.UserController{
		Engine: engine,
		Us:     &us,
	}
	fmt.Println(uc)
	errEngineRun := engine.Run(":8080")
	if errEngineRun != nil {
		return errEngineRun
	}
	return nil
}

func closeDb(pg *builtInSql.DB) {
	err := pg.Close()
	if err != nil {
		panic(err)
	}
}

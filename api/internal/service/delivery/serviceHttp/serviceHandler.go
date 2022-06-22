package serviceHttp

import (
	"github.com/labstack/echo/v4"
	serviceUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/service/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/tools"
	"net/http"
)

type ServiceHandler interface {
	Clear(ctx echo.Context) error
	Status(ctx echo.Context) error
}

type serviceHandler struct {
	ServiceUseCase serviceUseCase.ServiceUseCase
}

func NewServiceHandler(ServiceUseCase serviceUseCase.ServiceUseCase) ServiceHandler {
	return serviceHandler{
		ServiceUseCase: ServiceUseCase,
	}
}

// Clear godoc
// @Summary      Очистка всех данных в базе
// @Description  Безвозвратное удаленией всей пользовательской информации из базы данных.
//@Tags         service
// @Accept          application/json
// @Produce      application/json
// @Success      200    {string}  string "Очистка базы успешно завершена"
// @Router       /api/service/clear [post]
func (h serviceHandler) Clear(ctx echo.Context) error {
	if err := h.ServiceUseCase.Clear(); err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusServiceUnavailable)
	}
	return ctx.NoContent(http.StatusOK)
}

// Status godoc
// @Summary      Получение информации о базе данных
// @Description  Получение информации о базе данных.
//@Tags         service
// @Accept          application/json
// @Produce      application/json
// @Success      200    {object}  domain.Status "Кол-во записей в базе данных, включая помеченные как "удалённые"."
// @Router       /api/service/status [get]
func (h serviceHandler) Status(ctx echo.Context) error {
	stat, err := h.ServiceUseCase.Status()
	if err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusServiceUnavailable)
	}
	return ctx.JSON(http.StatusOK, *stat)
}

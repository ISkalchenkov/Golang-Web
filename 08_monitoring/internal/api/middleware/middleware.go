package middleware

import (
	"errors"
	"fmt"
	"server/internal/pkg/domain"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const RequestIDKey = "requestID"

func AuthEchoMiddleware(service domain.SessionService, logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			_, err := service.CheckSession(ctx.Request().Header)
			if err != nil {
				switch {
				case errors.Is(err, domain.ErrNoSession):
					return ctx.NoContent(401)
				default:
					logger.Errorw(fmt.Sprintf("service.CheckSession failed: %v", err),
						"request_id", ctx.Get(RequestIDKey),
						"remote_addr", ctx.Request().RemoteAddr,
						"method", ctx.Request().Method,
						"url", ctx.Request().URL,
					)
					return err
				}
			}
			return next(ctx)
		}
	}
}

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			reqID := uuid.New()
			ctx.Set(RequestIDKey, reqID)
			return next(ctx)
		}
	}
}

func AccessLog(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			start := time.Now()
			if err = next(ctx); err != nil {
				ctx.Error(err)
			}
			stop := time.Now()

			errStr := ""
			if err != nil {
				errStr = err.Error()
			}

			logger.Infow("New request",
				"request_id", ctx.Get(RequestIDKey),
				"remote_addr", ctx.Request().RemoteAddr,
				"method", ctx.Request().Method,
				"url", ctx.Request().URL,
				"error", errStr,
				"status", ctx.Response().Status,
				"bytes_in", ctx.Request().ContentLength,
				"bytes_out", ctx.Response().Size,
				"latency", stop.Sub(start).String(),
			)
			return
		}
	}
}

func Recover(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					errStr := fmt.Sprint(err)
					logger.Errorw(errStr,
						"request_id", ctx.Get(RequestIDKey),
						"remote_addr", ctx.Request().RemoteAddr,
						"method", ctx.Request().Method,
						"url", ctx.Request().URL,
					)
					ctx.Error(err)
				}
			}()
			return next(ctx)
		}
	}
}

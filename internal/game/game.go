package game

import (
	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/serialize/protobuf"
	"net/http"
	"strings"
	"tetris/config"
	"tetris/pkg/log"
)

func StartUp() {
	sc := config.ServerConfig
	services := &component.Components{}

	// todo: 可以按照room类型做分服，暂时没有按照gate分服

	// gate
	{
		service := newGateService()
		opts := []component.Option{
			component.WithName("g"),
			component.WithNameFunc(func(s string) string {
				return strings.ToLower(s)
			}),
		}
		services.Register(service, opts...)
	}

	// rooms
	{
		for _, v := range sc.Rooms {
			service := newRoomService(v)
			opts := []component.Option{
				component.WithName(service.serviceName),
				component.WithNameFunc(func(s string) string {
					return strings.ToLower(s)
				}),
			}
			services.Register(service, opts...)
		}
	}

	opts := []nano.Option{
		// websocket
		//nano.WithWSPath("/nano"),
		//nano.WithIsWebsocket(true),
		// 心跳
		//nano.WithHeartbeatInterval(5 * time.Second),
		//nano.WithDebugMode(),
		nano.WithLogger(log.GetLogger()),
		nano.WithComponents(services),
		nano.WithSerializer(protobuf.NewSerializer()),
		nano.WithCheckOriginFunc(func(request *http.Request) bool {
			return true
		}),
	}
	nano.Listen(sc.Addr, opts...)
}

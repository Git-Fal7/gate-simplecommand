package examples

import (
	"fmt"
	"time"

	simplecommand "github.com/git-fal7/gate-simplecommand"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type PingCommand struct {
	proxy *proxy.Proxy
}

func (cmd PingCommand) Execute(invocation simplecommand.SimpleCommandInvocation) error {
	player, ok := invocation.Source.(proxy.Player)
	if !ok {
		return nil
	}
	var ping int64 = int64(player.Ping() / time.Millisecond)
	var result string = fmt.Sprintf("Your ping is %dms", ping)
	if len(invocation.Arguments) >= 1 {
		args := invocation.Arguments
		target := cmd.proxy.PlayerByName(args[0])
		if target != nil && target.ID() != player.ID() {
			ping = int64(target.Ping() / time.Millisecond)
			result = fmt.Sprintf("%s's Ping is %dms", target.Username(), ping)
		}
	}
	player.SendMessage(&component.Text{
		Content: result,
	})
	return nil
}

func (cmd PingCommand) Suggest(invocation simplecommand.SimpleCommandInvocation) []string {
	player, ok := invocation.Source.(proxy.Player)
	if !ok {
		return nil
	}
	var list []string
	if len(invocation.Arguments) == 1 {
		for _, target := range cmd.proxy.Players() {
			if target.ID() != player.ID() {
				list = append(list, target.Username())
			}
		}
	}
	return list
}

func (cmd PingCommand) HasPermission(source command.Source) bool {
	return source.HasPermission("plugin.hasprong")
}

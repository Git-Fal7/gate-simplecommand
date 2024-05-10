package simplecommand

import (
	"strings"

	"go.minekube.com/brigodier"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type SimpleCommand interface {
	Execute(invocation SimpleCommandInvocation) error
	Suggest(invocation SimpleCommandInvocation) []string
	HasPermission(source command.Source) bool
}

type SimpleCommandInvocation struct {
	Arguments []string
	Source    command.Source
}

func RegisterCommand(proxy *proxy.Proxy, commandName string, simpleCommand SimpleCommand, aliases ...string) {
	proxy.Command().Register(simpleCommandToBrigoder(commandName, simpleCommand))
	for _, alias := range aliases {
		proxy.Command().Register(simpleCommandToBrigoder(alias, simpleCommand))
	}
}

func simpleCommandToBrigoder(commandName string, simpleCommand SimpleCommand) brigodier.LiteralNodeBuilder {
	return brigodier.Literal(commandName).
		Requires(command.Requires(func(c *command.RequiresContext) bool {
			return simpleCommand.HasPermission(c.Source)
		})).
		Executes(command.Command(func(c *command.Context) error {
			invocation := SimpleCommandInvocation{
				Arguments: nil,
				Source:    c.Source,
			}
			return simpleCommand.Execute(invocation)
		})).
		Then(brigodier.Argument("arguments", brigodier.StringPhrase).
			Suggests(command.SuggestFunc(func(c *command.Context, b *brigodier.SuggestionsBuilder) *brigodier.Suggestions {
				invocation := SimpleCommandInvocation{
					Arguments: strings.Split(c.String("arguments"), " "),
					Source:    c.Source,
				}
				for _, suggestion := range simpleCommand.Suggest(invocation) {
					b.Suggest(suggestion)
				}
				return b.Build()
			})).
			Executes(command.Command(func(c *command.Context) error {
				invocation := SimpleCommandInvocation{
					Arguments: strings.Split(c.String("arguments"), " "),
					Source:    c.Source,
				}
				return simpleCommand.Execute(invocation)
			})))
}

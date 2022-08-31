package cli

import (
	"fmt"

	"github.com/alexeyco/simpletable"
)

type providerListCmd struct {
}

func (c *providerListCmd) Run(ctx *context) error {
	ps, err := ctx.snapshoter.ListProviders("")
	if err != nil {
		return err
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	simplifyId := ctx.snapshoter.SimplifyId
	if ctx.globals.Verbose > 0 {
		simplifyId = func(id string) string {
			return id
		}
	}

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "Name"},
			{Text: "Version"},
			{Text: "Type"},
		},
	}

	for _, p := range ps {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{Text: simplifyId(p.ID)},
			{Text: p.Name},
			{Text: p.Version},
			{Text: p.Type},
		})
	}

	fmt.Println(table.String())

	return nil
}

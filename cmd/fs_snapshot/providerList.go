package main

import (
	"github.com/alexeyco/simpletable"
)

type providerListCmd struct {
	ServerArgs serverArgs `embed:""`
}

func (c *providerListCmd) Run(ctx *context) error {
	ps, err := ctx.snapshoter.ListProviders("")
	if err != nil {
		return err
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	simplifyId := ctx.snapshoter.SimplifyID
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

	ctx.console.Print(table.String())
	return nil
}

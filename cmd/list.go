package cli

import (
	"sort"

	"github.com/alexeyco/simpletable"
)

type listCmd struct {
	ServerArgs serverArgs `embed:""`
}

func (c *listCmd) Run(ctx *context) error {
	ps, err := ctx.snapshoter.ListSnapshots("")
	if err != nil {
		return err
	}

	if len(ps) == 0 {
		ctx.console.Print("No snapshots exist.")
		return nil
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	if ctx.globals.Verbose == 0 {
		table.Header.Cells = append(table.Header.Cells, []*simpletable.Cell{
			{Text: "ID"},
			{Text: "Set ID"},
			{Text: "Original path"},
			{Text: "Snapshot path"},
			{Text: "Creation"},
			{Text: "State"},
		}...)
	} else {
		table.Header.Cells = append(table.Header.Cells, []*simpletable.Cell{
			{Text: "ID"},
			{Text: "Set ID"},
			{Text: "Original path"},
			{Text: "Snapshot path"},
			{Text: "Creation"},
			{Text: "Provider"},
			{Text: "State"},
			{Text: "Attributes"},
		}...)
	}

	sort.Slice(ps, func(a, b int) bool {
		return ps[a].CreationTime.After(ps[b].CreationTime)
	})

	for _, p := range ps {
		if ctx.globals.Verbose == 0 {
			table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
				{Text: ctx.snapshoter.SimplifyID(p.ID)},
				{Text: ctx.snapshoter.SimplifyID(p.Set.ID)},
				{Text: p.OriginalPath},
				{Text: p.SnapshotPath},
				{Text: p.CreationTime.Local().Format("2006-01-02 15:04")},
				{Text: p.State},
			})

		} else {
			table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
				{Text: p.ID},
				{Text: p.Set.ID},
				{Text: p.OriginalPath},
				{Text: p.SnapshotPath},
				{Text: p.CreationTime.Local().Format("2006-01-02 15:04:05 -07")},
				{Text: p.Provider.Name},
				{Text: p.State},
				{Text: p.Attributes},
			})
		}
	}

	ctx.console.Print(table.String())

	return nil
}

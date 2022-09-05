package cli

import (
	"fmt"
	"sort"

	"github.com/alexeyco/simpletable"
)

type setListCmd struct {
	ServerArgs serverArgs `embed:""`
}

func (c *setListCmd) Run(ctx *context) error {
	ps, err := ctx.snapshoter.ListSets("")
	if err != nil {
		return err
	}

	if len(ps) == 0 {
		ctx.console.Print("No snapshot sets exist.")
		return nil
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	if ctx.globals.Verbose == 0 {
		table.Header.Cells = append(table.Header.Cells, []*simpletable.Cell{
			{Text: "ID"},
			{Text: "Creation"},
			{Text: "Snapshot count", Align: simpletable.AlignRight},
		}...)
	} else {
		table.Header.Cells = append(table.Header.Cells, []*simpletable.Cell{
			{Text: "ID"},
			{Text: "Creation"},
			{Text: "Snapshot count", Align: simpletable.AlignRight},
			{Text: "Snapshot count on creation", Align: simpletable.AlignRight},
		}...)
	}

	sort.Slice(ps, func(a, b int) bool {
		return ps[a].CreationTime.After(ps[b].CreationTime)
	})

	for _, p := range ps {
		if ctx.globals.Verbose == 0 {
			table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
				{Text: ctx.snapshoter.SimplifyID(p.ID)},
				{Text: p.CreationTime.Local().Format("2006-01-02 15:04")},
				{Text: fmt.Sprintf("%v", len(p.Snapshots)), Align: simpletable.AlignRight},
			})

		} else {
			table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
				{Text: p.ID},
				{Text: p.CreationTime.Local().Format("2006-01-02 15:04:05 -07")},
				{Text: fmt.Sprintf("%v", len(p.Snapshots)), Align: simpletable.AlignRight},
				{Text: fmt.Sprintf("%v", p.SnapshotCountOnCreation), Align: simpletable.AlignRight},
			})
		}
	}

	ctx.console.Print(table.String())
	return nil
}

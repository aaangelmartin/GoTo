package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/i18n"
)

func newLsCmd() *cobra.Command {
	var tagFilter string
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   i18n.T("ls_short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, st, err := loadState()
			if err != nil {
				return err
			}
			aliases := st.List()
			if flags.useJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(aliases)
			}
			if len(aliases) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), i18n.T("ls_empty"))
				return nil
			}

			header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF79C6"))
			name := lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD"))
			urlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))
			tagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
			dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))

			t := table.New().
				Border(lipgloss.RoundedBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#44475A"))).
				Headers(i18n.T("ls_col_name"), i18n.T("ls_col_url"), i18n.T("ls_col_tags"), i18n.T("ls_col_hits")).
				StyleFunc(func(row, col int) lipgloss.Style {
					if row == table.HeaderRow {
						return header.PaddingLeft(1).PaddingRight(1)
					}
					switch col {
					case 0:
						return name.PaddingLeft(1).PaddingRight(1)
					case 1:
						return urlStyle.PaddingLeft(1).PaddingRight(1)
					case 2:
						return tagStyle.PaddingLeft(1).PaddingRight(1)
					default:
						return dim.PaddingLeft(1).PaddingRight(1)
					}
				})

			for _, a := range aliases {
				if tagFilter != "" && !a.HasTag(tagFilter) {
					continue
				}
				t.Row(a.Name, a.URL, strings.Join(a.Tags, ", "), fmt.Sprintf("%d", a.HitCount))
			}
			fmt.Fprintln(cmd.OutOrStdout(), t.Render())
			return nil
		},
	}
	cmd.Flags().StringVar(&tagFilter, "tag", "", i18n.T("ls_tag"))
	cmd.Flags().BoolVar(&flags.useJSON, "json", false, i18n.T("ls_json"))
	return cmd
}

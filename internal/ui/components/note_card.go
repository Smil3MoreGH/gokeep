// internal/ui/components/note_card.go
package components

import (
	"github.com/Smil3MoreGH/gokeep/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/russross/blackfriday/v2"
)

// NoteCard represents a single note card component
type NoteCard struct {
	app.Compo

	Note      models.Note
	IsEditing bool
	OnEdit    func(noteID int64)
	OnDelete  func(noteID int64)
	OnSave    func(note models.Note)
	OnCancel  func()

	editTitle   string
	editContent string
}

func (c *NoteCard) OnMount(ctx app.Context) {
	c.editTitle = c.Note.Title
	c.editContent = c.Note.Content
}

func (c *NoteCard) Render() app.UI {
	if c.IsEditing {
		return c.renderEditMode()
	}
	return c.renderViewMode()
}

// renderViewMode renders the note in view mode
func (c *NoteCard) renderViewMode() app.UI {
	return app.Div().
		Class("note-card").
		Style("background-color", c.Note.Color).
		Body(
			// Title
			app.If(
				c.Note.Title != "",
				func() app.UI {
					return app.H3().Class("note-title").Text(c.Note.Title)
				},
			),

			// Content (with Markdown rendering)
			app.Div().
				Class("note-content").
				Body(
					app.Raw(c.renderMarkdown(c.Note.Content)),
				),

			// Actions
			app.Div().Class("note-actions").Body(
				app.Button().
					Class("btn-icon").
					Title("Edit note").
					OnClick(c.onEditClick).
					Text("✏️"),
				app.Button().
					Class("btn-icon").
					Title("Delete note").
					OnClick(c.onDeleteClick).
					Text("🗑️"),
			),

			// Timestamp
			app.Div().
				Class("note-timestamp").
				Text(c.Note.UpdatedAt.Format("Jan 2, 15:04")),
		)
}

// renderEditMode renders the note in edit mode
func (c *NoteCard) renderEditMode() app.UI {
	return app.Div().
		Class("note-card editing").
		Style("background-color", c.Note.Color).
		Body(
			// Title input
			app.Input().
				Type("text").
				Class("note-title-input").
				Value(c.editTitle).
				Placeholder("Title").
				OnInput(c.onTitleInput).
				AutoFocus(true),

			// Content textarea
			app.Textarea().
				Class("note-content-input").
				Placeholder("Take a note...").
				Rows(5).
				Value(c.editContent).
				OnInput(c.onContentInput),

			// Color picker
			c.renderColorPicker(),

			// Actions
			app.Div().Class("note-actions").Body(
				app.Button().
					Class("btn btn-primary").
					Text("Save").
					OnClick(c.onSaveClick),
				app.Button().
					Class("btn btn-secondary").
					Text("Cancel").
					OnClick(c.onCancelClick),
			),
		)
}

// renderColorPicker renders the color selection UI
func (c *NoteCard) renderColorPicker() app.UI {
	colors := []models.NoteColor{
		models.ColorWhite,
		models.ColorYellow,
		models.ColorOrange,
		models.ColorPink,
		models.ColorPurple,
		models.ColorBlue,
		models.ColorGreen,
		models.ColorGray,
	}

	return app.Div().Class("color-picker").Body(
		app.Range(colors).Slice(func(i int) app.UI {
			color := string(colors[i])
			return app.Button().
				Class("color-option").
				Style("background-color", color).
				OnClick(func(ctx app.Context, e app.Event) {
					c.Note.Color = color
					c.Update()
				}).
				Body(
					app.If(
						c.Note.Color == color,
						func() app.UI {
							return app.Span().Text("✓")
						},
					),
				)
		}),
	)
}

// renderMarkdown converts markdown content to HTML
func (c *NoteCard) renderMarkdown(content string) string {
	if content == "" {
		return ""
	}

	// Configure blackfriday options
	extensions := blackfriday.CommonExtensions | blackfriday.Autolink

	// Render markdown to HTML
	html := blackfriday.Run([]byte(content), blackfriday.WithExtensions(extensions))

	return string(html)
}

// Event handlers

func (c *NoteCard) onEditClick(ctx app.Context, e app.Event) {
	if c.OnEdit != nil {
		c.OnEdit(c.Note.ID)
	}
}

func (c *NoteCard) onDeleteClick(ctx app.Context, e app.Event) {
	if c.OnDelete != nil {
		// Simple confirmation
		if app.Window().Call("confirm", "Delete this note?").Bool() {
			c.OnDelete(c.Note.ID)
		}
	}
}

func (c *NoteCard) onTitleInput(ctx app.Context, e app.Event) {
	value := ctx.JSSrc().Get("value").String()
	c.Dispatch(func(ctx app.Context) {
		c.editTitle = value
	})
}

func (c *NoteCard) onContentInput(ctx app.Context, e app.Event) {
	value := ctx.JSSrc().Get("value").String()
	c.Dispatch(func(ctx app.Context) {
		c.editContent = value
	})
}

func (c *NoteCard) onSaveClick(ctx app.Context, e app.Event) {
	if c.OnSave != nil {
		c.Note.Title = c.editTitle
		c.Note.Content = c.editContent
		c.OnSave(c.Note)
	}
}

func (c *NoteCard) onCancelClick(ctx app.Context, e app.Event) {
	if c.OnCancel != nil {
		// Reset to original values
		c.editTitle = c.Note.Title
		c.editContent = c.Note.Content
		c.OnCancel()
	}
}

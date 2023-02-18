package main

import (
  "errors"
  "fmt"
  "github.com/adrianescat/lets-go/internal/models"
  "github.com/adrianescat/lets-go/internal/validator"
  "github.com/julienschmidt/httprouter" // New import
  "net/http"
  "strconv"
)

// Remove the explicit FieldErrors struct field and instead embed the Validator
// type. Embedding this means that our snippetCreateForm "inherits" all the
// fields and methods of our Validator type (including the FieldErrors field).
type snippetCreateForm struct {
  Title   string
  Content string
  Expires int
  validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
  snippets, err := app.snippets.Latest()
  if err != nil {
	 app.serverError(w, err)
	 return
  }

  data := app.newTemplateData(r)
  data.Snippets = snippets

  app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
  params := httprouter.ParamsFromContext(r.Context())

  id, err := strconv.Atoi(params.ByName("id"))
  if err != nil || id < 1 {
	 app.notFound(w)
	 return
  }

  snippet, err := app.snippets.Get(id)
  if err != nil {
	 if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
	 } else {
		app.serverError(w, err)
	 }
	 return
  }

  data := app.newTemplateData(r)
  data.Snippet = snippet

  app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)

  // Initialize a new createSnippetForm instance and pass it to the template.
  // Notice how this is also a great opportunity to set any default or
  // 'initial' values for the form --- here we set the initial value for the
  // snippet expiry to 365 days.
  data.Form = snippetCreateForm{
	 Expires: 365,
  }

  app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()

  if err != nil {
	 app.clientError(w, http.StatusBadRequest)
	 return
  }

  expires, err := strconv.Atoi(r.PostForm.Get("expires"))

  if err != nil {
	 app.clientError(w, http.StatusBadRequest)
	 return
  }

  form := snippetCreateForm{
	 Title:   r.PostForm.Get("title"),
	 Content: r.PostForm.Get("content"),
	 Expires: expires,
	 // Remove the FieldErrors assignment from here.
  }

  // Because the Validator type is embedded by the snippetCreateForm struct,
  // we can call CheckField() directly on it to execute our validation checks.
  // CheckField() will add the provided key and error message to the
  // FieldErrors map if the check does not evaluate to true. For example, in
  // the first line here we "check that the form.Title field is not blank". In
  // the second, we "check that the form.Title field has a maximum character
  // length of 100" and so on.
  form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
  form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
  form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
  form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

  // Use the Valid() method to see if any of the checks failed. If they did,
  // then re-render the template passing in the form in the same way as
  // before.
  if !form.Valid() {
	 data := app.newTemplateData(r)
	 data.Form = form
	 app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
	 return
  }

  id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

  if err != nil {
	 app.serverError(w, err)
	 return
  }

  http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

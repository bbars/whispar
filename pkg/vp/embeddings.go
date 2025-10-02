package vp

type EmModel struct {
	// TODO: guess
	ModelViews []ModelView `vp:"_modelViews,omitempty"`

	// Edit restriction flag.
	ModelEditable bool `vp:"_modelEditable,omitempty"`
}

type EmInfo struct {
	// Author use name.
	PmAuthor string `vp:"pmAuthor,omitempty"`

	// Creation UNIX-timestamp in millis.
	PmCreateDateTime TimeString `vp:"pmCreateDateTime,omitempty"`

	// Last modification UNIX-timestamp in millis.
	PmLastModified TimeString `vp:"pmLastModified,omitempty"`

	// HTML-formatted documentation.
	Documentation string `vp:"documentation,omitempty"`

	// Plaintext documentation.
	DocumentationPlain string `vp:"documentation_plain,omitempty"`
}

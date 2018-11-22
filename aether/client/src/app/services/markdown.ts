// Services > Markdown

// This service initialises a global singleton markdown renderer object that we will be sharing whenever we need to convert markdown to HTML.

var Marked = require('../../../node_modules/marked')
var DOMPurify = require('../../../node_modules/dompurify')
// ^ Instead of relying on Marked's sanitiser, this additional sanitiser protects the app from insecure user input, since Markdown otherwise can pass through HTML tags like <script> untouched.

export { }

// We create a new renderer and inject it to marked, because we do not want user-provided image links to auto-load for user privacy reasons.
var renderer = new Marked.Renderer()
renderer.image = function(href: string, title: string, text: string) {
  return `![` + text + `](` + href + ` ` + title + `)`
}
Marked.setOptions({
  breaks: true,
  renderer: renderer
})

function MarkedRenderer(input: any) {
  return Marked(DOMPurify.sanitize(input, { ALLOWED_TAGS: [] }))
  // No HTML tags allowed - but you can use markdown to link stuff, etc.
  // This might be interesting to add in the future: {SAFE_FOR_TEMPLATES: true}
  // We don't really need it because we are binding to v-html so that is not evaluated for templates anyway. The side effect is that if somebody attempts to use {{ }} as part of the normal conversation (i.e. copy/pasting code) that would be stripped out, which isn't great. Let's do this only if this becomes a problem.
}

console.log('Markdown renderer initialised.')

module.exports = MarkedRenderer
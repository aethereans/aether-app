// Services > Markdown

// This service initialises a global singleton markdown renderer object that we will be sharing whenever we need to convert markdown to HTML.

// var Marked = require('../../../node_modules/marked')
var DOMPurify = require('../../../node_modules/dompurify')
// ^ Instead of relying on Marked's sanitiser, this additional sanitiser protects the app from insecure user input, since Markdown otherwise can pass through HTML tags like <script> untouched.
var HighlightJS = require('highlight.js'); // https://highlightjs.org/

export { }

// // We create a new renderer and inject it to marked, because we do not want user-provided image links to auto-load for user privacy reasons.
// var renderer = new Marked.Renderer()
// renderer.image = function(href: string, title: string, text: string) {
//   return `![` + text + `](` + href + ` ` + title + `)`
// }
// Marked.setOptions({
//   breaks: true,
//   renderer: renderer,
// })

// function MarkedRenderer(input: any) {
//   return Marked(DOMPurify.sanitize(input, { ALLOWED_TAGS: [] }))
//   // No HTML tags allowed - but you can use markdown to link stuff, etc.
//   // This might be interesting to add in the future: {SAFE_FOR_TEMPLATES: true}
//   // We don't really need it because we are binding to v-html so that is not evaluated for templates anyway. The side effect is that if somebody attempts to use {{ }} as part of the normal conversation (i.e. copy/pasting code) that would be stripped out, which isn't great. Let's do this only if this becomes a problem.
// }

/*New renderer*/

var MarkdownIt = require('markdown-it')
var rnd = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: false,
  highlight: function(str: string, lang: string) {
    if (lang && HighlightJS.getLanguage(lang)) {
      try {
        return HighlightJS.highlight(lang, str).value
      } catch (__) { }
    }

    return ''; // use external default escaping
  }
})

// Adds hashtag support.
rnd.linkify.add('#', {
  validate: function(text: any, pos: any, self: any) {
    var tail = text.slice(pos)

    if (!self.re.hashtag) {
      self.re.hashtag = new RegExp(
        '^([a-zA-Z0-9_]){1,15}(?!_)(?=$|' + self.re.src_ZPCc + ')'
      );
    }
    if (self.re.hashtag.test(tail)) {
      // Linkifier allows punctuation chars before prefix,
      // but we additionally disable `@` ("@@mention" is invalid)
      if (pos >= 2 && tail[pos - 2] === '@') {
        return false;
      }
      return tail.match(self.re.hashtag)[0].length;
    }
    return 0;
  },
  normalize: function(match: any) {
    // console.log(match)
    match.url = '#/searchscope/content?searchHashtag=' + match.text.replace(/^#/, '')
  }
})

// Adds aether:// link support.

rnd.linkify.add('aether://', {
  validate: function(text: any, pos: any, self: any) {
    var tail = text.slice(pos);

    if (!self.re.aetherlink) {
      self.re.aetherlink = /^[\w.-][\w\-\._~:/?#[\]@!\$&'\(\)\*\+",;=.]+$/
      console.log(self.re.aetherlink)
    }
    if (self.re.aetherlink.test(tail)) {
      // Linkifier allows punctuation chars before prefix,
      // but we additionally disable `@` ("@@mention" is invalid)
      if (pos >= 2 && tail[pos - 2] === '@') {
        return false;
      }
      return tail.match(self.re.aetherlink)[0].length;
    }
    return 0;
  },
  normalize: function(match: any) {
    match.url = '#/' + match.url.replace('aether://', '');
  }
});

function MarkdownItRenderer(input: any) {
  return rnd.render(DOMPurify.sanitize(input, { ALLOWED_TAGS: [] }))
}

// console.log('Markdown renderer initialised.')

// module.exports = MarkedRenderer
module.exports = MarkdownItRenderer

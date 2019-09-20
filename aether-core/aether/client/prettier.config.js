module.exports = {
  // --------------------------------------------------------------------
  // printWidth
  // --------------------------------------------------------------------
  //
  // @param {int} "prettier_options.printWidth"
  // @default 80
  //
  // Fit code within this line limit.
  // --------------------------------------------------------------------

  printWidth: 80,

  // --------------------------------------------------------------------
  // tabWidth
  // --------------------------------------------------------------------
  //
  // @param {int} "prettier_options.tabWidth"
  // @default 2
  //
  // Specify the number of spaces per indentation-level.
  //
  // IMPORTANT: By default, "tabWidth" is automatically set using the
  // SublimeText configured value for "tab_size". To disable this
  // behavior, you must first change the "disable_tab_width_auto_detection"
  // setting to "true".
  // --------------------------------------------------------------------

  tabWidth: 2,

  // --------------------------------------------------------------------
  // singleQuote
  // --------------------------------------------------------------------
  //
  // @param {bool} "prettier_options.singleQuote"
  // @default false
  //
  // If true, will use single instead of double quotes.
  // --------------------------------------------------------------------

  singleQuote: true,

  // --------------------------------------------------------------------
  // trailingComma
  // --------------------------------------------------------------------
  //
  // @param {string} "prettier_options.trailingComma"
  // @default "none"
  //
  // Controls the printing of trailing commas wherever possible.
  //
  // Valid options:
  //
  // "none" - No trailing commas
  // "es5"  - Trailing commas where valid in ES5 (objects, arrays, etc)
  // "all"  - Trailing commas wherever possible (function arguments)
  // --------------------------------------------------------------------

  trailingComma: 'es5',

  // --------------------------------------------------------------------
  // bracketSpacing
  // --------------------------------------------------------------------
  //
  // @param {bool} "prettier_options.bracketSpacing"
  // @default true
  //
  // Controls the printing of spaces inside array and objects.
  // --------------------------------------------------------------------

  bracketSpacing: true,

  // --------------------------------------------------------------------
  // jsxBracketSameLine
  // --------------------------------------------------------------------
  //
  // @param {bool} "prettier_options.jsxBracketSameLine"
  // @default false
  //
  // If true, puts the `>` of a multi-line jsx element at the end of
  // the last line instead of being alone on the next line.
  // --------------------------------------------------------------------

  jsxBracketSameLine: false,

  // --------------------------------------------------------------------
  // parser
  // --------------------------------------------------------------------
  //
  // @param {string} "prettier_options.parser"
  // @default "babylon"
  //
  // Which parser to use. Valid options are "flow", "babylon",
  // "typescript", "css", "json", "graphql", "markdown" and "yaml".
  //
  // NOTE: The `parser` option is automatically set by the plug-in
  // (JsPrettier), based on the contents of current file or selection.
  // --------------------------------------------------------------------

  // parser: "babel",

  // --------------------------------------------------------------------
  // semi
  // --------------------------------------------------------------------
  //
  // @param {bool} "prettier_options.semi"
  // @default true
  //
  // Whether to add a semicolon at the end of every line (semi: true), or
  // only at the beginning of lines that may introduce ASI failures (semi: false)
  // --------------------------------------------------------------------

  semi: false,

  // --------------------------------------------------------------------
  // requirePragma
  // --------------------------------------------------------------------
  //
  // @param {bool} "prettier_options.requirePragma"
  // @default false
  //
  // Prettier can restrict itself to only format files that contain a
  // special comment, called a pragma, at the top of the file. This is
  // very useful when gradually transitioning large, unformatted codebases
  // to prettier.
  // --------------------------------------------------------------------

  requirePragma: false,

  // --------------------------------------------------------------------
  // proseWrap
  // --------------------------------------------------------------------
  //
  // @param {string} "prettier_options.proseWrap"
  // @default "preserve"
  //
  // (Markdown and YAML Only) By default, Prettier will wrap markdown text
  // as-is since some services use a linebreak-sensitive renderer, e.g.
  // GitHub comment and BitBucket. In some cases you may want to rely on
  // SublimeText soft wrapping instead, so this option allows you to opt
  // out with "never".
  //
  // Valid options:
  //
  // "always" - Wrap prose if it exceeds the print width.
  // "never" - Do not wrap prose.
  // "preserve" (default) - Wrap prose as-is. available in v1.9.0+
  // --------------------------------------------------------------------

  proseWrap: 'preserve',

  // --------------------------------------------------------------------
  // arrowParens
  // --------------------------------------------------------------------
  //
  // @param {string} "prettier_options.arrowParens"
  // @default "avoid"
  //
  // Include parentheses around a sole arrow function parameter.
  //
  // Valid Options:
  //
  // - "avoid" (default) - Omit parentheses when possible. Example: `x => x`
  // - "always" - Always include parentheses. Example: `(x) => x`
  // --------------------------------------------------------------------

  arrowParens: 'avoid',

  // --------------------------------------------------------------------
  // htmlWhitespaceSensitivity
  // --------------------------------------------------------------------
  //
  // @param {string} "prettier_options.htmlWhitespaceSensitivity"
  // @default "css"
  //
  // Specify the global whitespace sensitivity for HTML files.
  //
  // Valid Options:
  //
  // - "css" - Respect the default value of CSS display property.
  // - "strict" - Whitespaces are considered sensitive.
  // - "ignore" - Whitespaces are considered insensitive.
  // --------------------------------------------------------------------

  htmlWhitespaceSensitivity: 'css',
}

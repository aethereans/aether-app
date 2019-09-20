// Services > Placeholder picker
// Get a quote every time a placeholder for a text field is blank.

export {}

interface Placeholder {
  Placeholder: string
}
let phs: Placeholder[] = [
  {
    Placeholder: 'Cancel all my meetings, somebody is wrong on the Internet',
  },
  {
    Placeholder: '[citation needed]',
  },
  {
    Placeholder: '‘Use the force, Harry’ — Gandalf',
  },
  {
    Placeholder:
      'This is my textbox, there are many like it but this one is mine',
  },
  {
    Placeholder: "A textbox of one's own",
  },
  {
    Placeholder: "It's turtles all the way down",
  },
]

function GetPlaceholder(): Placeholder {
  let rnd = Math.floor(Math.random() * phs.length)
  return phs[rnd]
}

module.exports = GetPlaceholder

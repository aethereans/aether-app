"use strict";
// Services > Quote picker
// This service allows us to render a quote every time a no-content object is shown.
Object.defineProperty(exports, "__esModule", { value: true });
var quotes = [
    {
        Quote: 'If you want to build a ship, don’t drum up the men to gather wood, divide the work and give orders. Instead, teach them to yearn for the vast and endless sea.',
        Author: 'Antoine de Saint-Exupéry',
    },
    {
        Quote: 'If you speak the truth, have a foot in the stirrup.',
        Author: 'Turkish proverb',
    },
    {
        Quote: 'Real stupidity beats artificial intelligence every time.',
        Author: 'Terry Pratchett',
    },
    {
        Quote: 'True glory consists in doing what deserves to be written; in writing what deserves to be read.',
        Author: 'Pliny the Elder',
    },
    {
        Quote: 'All the world’s a stage, And all the men and women merely players. They have their exits and their entrances; And one man in his time plays many parts.',
        Author: 'Shakespeare',
    },
    {
        Quote: 'A designer knows he has achieved perfection not when there is nothing left to add, but when there is nothing left to take away.',
        Author: 'Antoine de Saint-Exupéry',
    },
    {
        Quote: 'One doesn’t discover new lands without losing sight of the shore.',
        Author: 'André Gide',
    },
    {
        Quote: 'The trouble with having an open mind, of course, is that people will insist on coming along and trying to put things in it.',
        Author: 'Terry Pratchett',
    },
    {
        Quote: 'Victorious warriors win first and then go to war, while defeated warriors go to war first and then seek to win.',
        Author: 'Sun Tzu',
    },
    {
        Quote: 'If the radiance of a thousand suns were to burst at once into the sky, that would be like the splendor of the Mighty One... I am become Death, the Shatterer of Worlds.',
        Author: 'Oppenheimer, quoting Bhagavad Gita',
    },
    {
        Quote: 'As to diseases make a habit of two things — to help, or at least, to do no harm.',
        Author: 'Hippocrates',
    },
    {
        Quote: 'Never trust a computer you can’t throw out a window.',
        Author: 'Steve Wozniak',
    },
    {
        Quote: 'The pen is mightier than the sword if the sword is very short, and the pen is very sharp.',
        Author: 'Terry Pratchett',
    },
    {
        Quote: 'For once you have tasted flight you will walk the earth with your eyes turned skywards, for there you have been, and there you will long to return.',
        Author: 'Leonardo da Vinci',
    },
    {
        Quote: 'Any sufficiently advanced technology is indistinguishable from magic.',
        Author: 'Arthur C. Clarke',
    },
    {
        Quote: '“Then one fine mornin’ she puts on a New York station. You know her life was saved by Rock ‘n’ Roll.”',
        Author: 'The Velvet Underground, Rock and Roll',
    },
    {
        Quote: '“Beep... beep... beep... beep...”',
        Author: 'Sputnik I',
    },
    {
        Quote: 'Most of the great triumphs and tragedies of history are caused not by people being fundamentally good or fundamentally evil, but by people being fundamentally people.',
        Author: 'Terry Pratchett',
    },
    {
        Quote: 'If art interprets our dreams, the computer executes them in the guise of programs.',
        Author: 'Alan Perlis',
    },
    {
        Quote: 'There are two major products that come out of Berkeley: LSD and UNIX. We don’t believe this to be a coincidence.',
        Author: 'Jeremy S. Anderson',
    },
    {
        Quote: 'To iterate is human, to recurse divine.',
        Author: 'Peter Deutsch',
    },
    {
        Quote: 'There should be one — and preferably only one — obvious way to do it.',
        Author: 'Tim Peters, The Zen of Python',
    },
    {
        Quote: 'It is well known that a vital ingredient of success is not knowing that what you’re attempting can’t be done.',
        Author: 'Terry Pratchett',
    },
    {
        Quote: 'Any sufficiently complicated computer program contains an ad hoc, informally–specified, bug–ridden, slow implementation of half of Lisp.',
        Author: 'Greenspun’s Tenth Rule of Programming',
    },
    {
        Quote: 'A common mistake people make when trying to design something completely foolproof is to underestimate the ingenuity of complete fools.',
        Author: 'Hitchhiker’s Guide to the Galaxy, Douglas Adams',
    },
    {
        Quote: 'A little learning is a dangerous thing; Drink deep, or taste not the Pierian spring.',
        Author: 'Alexander Pope',
    },
    {
        Quote: 'It’s not worth doing something unless someone, somewhere, would much rather you weren’t doing it.',
        Author: 'Terry Pratchett',
    },
    {
        Quote: 'Of Manners gentle, of Affections mild; In Wit a man; Simplicity, a child.',
        Author: 'Alexander Pope',
    },
    {
        Quote: 'Beauty is the ultimate defense against complexity.',
        Author: 'David Gelernter',
    },
    {
        Quote: 'When the only tool you own is a hammer, every problem begins to resemble a nail.',
        Author: 'Abraham Maslow',
    },
    {
        Quote: 'When in doubt, use Caslon.',
        Author: 'William Caslon',
    },
    {
        Quote: 'You get what anybody gets — you get a lifetime.',
        Author: 'Sandman: Preludes and Nocturnes, Neil Gaiman',
    },
    {
        Quote: ' ‘All that is gold does not glitter, <br>Not all those who wander are lost; <br>The old that is strong does not wither, <br>Deep roots are not reached by the frost’',
        Author: 'Fellowship of the Ring, J.R.R. Tolkien',
    },
    {
        Quote: 'I’d rather be a rising ape than a falling angel.',
        Author: 'Terry Pratchett',
    },
    {
        Quote: 'In the beginning, the universe was created. This made a lot of people very angry and has been widely regarded as a bad move.',
        Author: 'Hitchhiker’s Guide to the Galaxy, Douglas Adams',
    },
    {
        Quote: 'None are more hopelessly enslaved than those who falsely believe they are free.',
        Author: 'Johann Wolfgang von Goethe',
    },
    {
        Quote: 'A change in perspective is worth 80 IQ points.',
        Author: 'Alan Kay',
    },
    {
        Quote: 'The presence of those seeking the truth is infinitely to be preferred to the presence of those who think they’ve found it.',
        Author: 'Terry Pratchett',
    },
    {
        Quote: 'Any society that would give up a little liberty to gain a little security will deserve neither and lose both.',
        Author: 'Benjamin Franklin',
    },
    {
        Quote: 'And on the pedestal these words appear: <br>‘My name is Ozymandias, king of kings: <br>Look on my works, ye Mighty, and despair!’ <br>Nothing beside remains.',
        Author: 'Ozymandias, Percy Bysshe Shelley',
    },
    {
        Quote: 'The birds have vanished down the sky. <br>Now the last cloud drains away. <br><br>We sit together, the mountain and me, <br>until only the mountain remains.',
        Author: 'Li Bai 李白',
    },
    {
        Quote: 'The fault, my dear Brutus, is not in our stars, <br>But in ourselves, that we are underlings.',
        Author: 'Julius Caesar, Shakespeare',
    },
    {
        Quote: 'We must plan for freedom, and not only for security, if for no other reason than that only freedom can make security secure.',
        Author: 'Karl Popper',
    },
    {
        Quote: 'Lies, damned lies and statistics.',
        Author: 'Benjamin Disraeli',
    },
    {
        Quote: 'The more pleasures a man captures, the more masters he will have to serve.',
        Author: 'Lucius Annaeus Seneca',
    },
    {
        Quote: 'Anyone who isn’t embarrassed of who they were last year probably isn’t learning enough. ',
        Author: 'Alain de Botton',
    },
];
var initd = false;
function fixWidows() {
    for (var i = 0; i < quotes.length; i++) {
        ;
        (function (i) {
            quotes[i].Quote = quotes[i].Quote.replace(/ (?=[^ ]*$)/i, '&nbsp;');
        })(i);
    }
    initd = true;
}
function GetQuote() {
    if (!initd) {
        fixWidows();
    }
    var rnd = Math.floor(Math.random() * quotes.length);
    return quotes[rnd];
}
module.exports = GetQuote;
//# sourceMappingURL=quotepicker.js.map
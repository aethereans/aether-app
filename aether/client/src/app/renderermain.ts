/*
This is the main entry point to the client app. See app.vue for the start logic, and globally-applicable css.
*/

/*----------  Electron and our non-GUI services  ----------*/
/*
  Main thing to realise is that there are two processes on the electron side, the main and the renderer. Main is basically a node process, and renderer is basically very similar to script that linked via the <script> tag of a frame that the node has surfaced, which means it has much fewer privileges (though still has access to Node APIs).

  As a result, we start the frontend binary from the main side, but we establish the client gRPC server on the renderer side when it is time to connect to that server, since the data frontend provides needs to be delivered to the renderer, not the main process.
*/

export { }

// Electron IPC setup before doing anything else
require('./services/eipc/eipc-renderer') // Register IPC events
var ipc = require('../../node_modules/electron-better-ipc') // Register IPC caller
// ^ Heads up, there are some IPC events registered in this renderermain, too.

// const unhandled = require('../../node_modules/electron-unhandled')
// unhandled()

const clapiserver = require('./services/clapiserver/clapiserver')
const feapiconsumer = require('./services/feapiconsumer/feapiconsumer')
const globalMethods = require('./services/globals/methods')
const clientAPIServerPort: number = clapiserver.StartClientAPIServer()

console.log('attempting to call get frontend ready')
ipc
  .callMain('GetFrontendReady')
  .then(function(resp: any) {
    console.log('frontend ready response received')
    console.log(resp)
  })
  .catch((err: any) => {
    console.log('fe ready promise rejected.')
    console.log(err)
  })

/*----------  Call mainmain to ask software update state  ----------*/

ipc.callMain('AskNewUpdateReady')

console.log('renderer client api server port: ', clientAPIServerPort)
ipc
  .callMain('SetClientAPIServerPort', clientAPIServerPort)
  .then(function(feDaemonStarted: boolean) {
    if (!feDaemonStarted) {
      // It's an Electron refresh, not a cold start.
      feapiconsumer.Initialise()
    }
  })
  .catch((err: any) => {
    console.log('this is the promise error:')
    console.log(err)
  })

/*----------  Vue + its plugins  ----------*/

const isDev = require('electron-is-dev')
if (isDev) {
  var Vue = require('../../node_modules/vue/dist/vue.js') // Production
} else {
  var Vue = require('../../node_modules/vue/dist/vue.min.js') // Production
}
// var Vue = require('../../node_modules/vue/dist/vue.js') // Development
var VueRouter = require('../../node_modules/vue-router').default

Vue.use(VueRouter)

// Register icons for our own use.
var Icon = require('../../node_modules/vue-awesome')
Vue.component('icon', Icon)

// Register the click-outside component
var ClickOutside = require('../../node_modules/v-click-outside')
Vue.use(ClickOutside)

/*----------  Third party dependencies  ----------*/

var Mousetrap = require('../../node_modules/mousetrap')
// var Spinner = require('../../node_modules/vue-simple-spinner')

/*----------  Components  ----------*/

// Global component declarations - do it here once.
Vue.component('a-app', require('./components/a-app.vue').default)
Vue.component('a-header', require('./components/a-header.vue').default)
Vue.component(
  'a-header-icon',
  require('./components/a-header-icon.vue').default
)
Vue.component('a-sidebar', require('./components/a-sidebar.vue').default)
Vue.component(
  'a-boardheader',
  require('./components/a-boardheader.vue').default
)
Vue.component('a-tabs', require('./components/a-tabs.vue').default)
Vue.component(
  'a-thread-listitem',
  require('./components/listitems/a-thread-listitem.vue').default
)
Vue.component(
  'a-vote-action',
  require('./components/a-vote-action.vue').default
)
Vue.component(
  'a-thread-header-entity',
  require('./components/a-thread-header-entity.vue').default
)
Vue.component('a-post', require('./components/a-post.vue').default)
Vue.component(
  'a-side-header',
  require('./components/a-side-header.vue').default
)
Vue.component(
  'a-breadcrumbs',
  require('./components/a-breadcrumbs.vue').default
)
Vue.component('a-username', require('./components/a-username.vue').default)
Vue.component('a-timestamp', require('./components/a-timestamp.vue').default)
Vue.component(
  'a-globalscopeheader',
  require('./components/a-globalscopeheader.vue').default
)
Vue.component(
  'a-board-listitem',
  require('./components/listitems/a-board-listitem.vue').default
)
Vue.component('a-hashimage', require('./components/a-hashimage.vue').default)
Vue.component('a-no-content', require('./components/a-no-content.vue').default)
Vue.component('a-markdown', require('./components/a-markdown.vue').default)
Vue.component(
  'a-avatar-block',
  require('./components/a-avatar-block.vue').default
)
Vue.component('a-composer', require('./components/a-composer.vue').default)
Vue.component('a-ballot', require('./components/a-ballot.vue').default)
Vue.component(
  'a-progress-bar',
  require('./components/a-progress-bar.vue').default
)
Vue.component(
  'a-inflight-info',
  require('./components/a-inflight-info.vue').default
)
Vue.component(
  'a-info-marker',
  require('./components/a-info-marker.vue').default
)
Vue.component('a-spinner', require('./components/a-spinner.vue').default)
Vue.component('a-notfound', require('./components/a-notfound.vue').default)
Vue.component('a-guidelight', require('./components/a-guidelight.vue').default)
Vue.component(
  'a-home-header',
  require('./components/a-home-header.vue').default
)
Vue.component(
  'a-popular-header',
  require('./components/a-popular-header.vue').default
)
Vue.component(
  'a-notifications-icon',
  require('./components/a-notifications-icon.vue').default
)
Vue.component(
  'a-notification-entity',
  require('./components/a-notification-entity.vue').default
)
Vue.component(
  'a-main-app-loader',
  require('./components/a-main-app-loader.vue').default
)
Vue.component(
  'a-global-header',
  require('./components/a-global-header.vue').default
)
Vue.component('a-fin-puck', require('./components/a-fin-puck.vue').default)
Vue.component(
  'a-bootstrapper',
  require('./components/a-bootstrapper.vue').default
)
Vue.component(
  'a-fingerprint',
  require('./components/a-fingerprint.vue').default
)
Vue.component(
  'a-settings-block',
  require('./components/a-settings-block.vue').default
)
Vue.component(
  'a-software-update-icon',
  require('./components/a-software-update-icon.vue').default
)
Vue.component(
  'a-patreon-button',
  require('./components/a-patreon-button.vue').default
)
Vue.component(
  'a-crypto-fund-button',
  require('./components/a-crypto-fund-button.vue').default
)
Vue.component('a-boardname', require('./components/a-boardname.vue').default)
Vue.component(
  'a-search-icon',
  require('./components/a-search-icon.vue').default
)
Vue.component(
  'a-search-header',
  require('./components/a-search-header.vue').default
)
Vue.component(
  'a-user-listitem',
  require('./components/listitems/a-user-listitem.vue').default
)
Vue.component(
  'a-post-listitem',
  require('./components/listitems/a-post-listitem.vue').default
)
Vue.component('a-threadname', require('./components/a-threadname.vue').default)
Vue.component('a-link', require('./components/a-link.vue').default)

/*----------  Third party components  ----------*/

Vue.component(
  'vue-simple-spinner',
  require('../../node_modules/vue-simple-spinner')
)

/*----------  Places  ----------*/

const Home = require('./components/locations/home.vue').default
const Popular = require('./components/locations/popular.vue').default

/*----------  Global scope (whole network, i.e. list of boards)  ----------*/
const GlobalScope = require('./components/locations/globalscope.vue').default
const GlobalRoot = require('./components/locations/globalscope/globalroot.vue')
  .default
const GlobalSubbed = require('./components/locations/globalscope/subbedroot.vue')
  .default

/*----------  Board scope (board entity + list of threads)  ----------*/
const NewBoard = require('./components/locations/globalscope/newboard.vue')
  .default
const BoardScope = require('./components/locations/boardscope.vue').default
const BoardRoot = require('./components/locations/boardscope/boardroot.vue')
  .default
const BoardInfo = require('./components/locations/boardscope/boardinfo.vue')
  .default
const ModActivity = require('./components/locations/boardscope/modactivity.vue')
  .default
const Elections = require('./components/locations/boardscope/elections.vue')
  .default
const Reports = require('./components/locations/boardscope/reports.vue').default

/*----------  Thread scope (thread entity + list of posts)  ----------*/
const NewThread = require('./components/locations/boardscope/newthread.vue')
  .default
const ThreadScope = require('./components/locations/threadscope.vue').default

/*----------  Settings scope  ----------*/
const SettingsScope = require('./components/locations/settingsscope.vue')
  .default
const SettingsRoot = require('./components/locations/settingsscope/settingsroot.vue')
  .default
const Defaults = require('./components/locations/settingsscope/defaults.vue')
  .default
const Shortcuts = require('./components/locations/settingsscope/shortcuts.vue')
  .default
const AdvancedSettings = require('./components/locations/settingsscope/advancedsettings.vue')
  .default
const About = require('./components/locations/settingsscope/about.vue').default
const Membership = require('./components/locations/settingsscope/membership.vue')
  .default
const Changelog = require('./components/locations/settingsscope/changelog.vue')
  .default
const AdminsQuickstart = require('./components/locations/settingsscope/adminsquickstart.vue')
  .default
const Intro = require('./components/locations/settingsscope/intro.vue').default
const NewUser = require('./components/locations/settingsscope/newuser.vue')
  .default
const SFWList = require('./components/locations/settingsscope/sfwlist.vue')
  .default
const Modship = require('./components/locations/settingsscope/modship.vue')
  .default
const Namemint = require('./components/locations/settingsscope/namemint.vue')
  .default

/*----------  User scope  ----------*/
const UserScope = require('./components/locations/userscope.vue').default
const UserRoot = require('./components/locations/userscope/userroot.vue')
  .default
const UserBoards = require('./components/locations/userscope/userboards.vue')
  .default
const UserThreads = require('./components/locations/userscope/userthreads.vue')
  .default
const UserPosts = require('./components/locations/userscope/userposts.vue')
  .default
const Notifications = require('./components/locations/userscope/notifications.vue')
  .default

/*----------  Status scope  ----------*/

const Status = require('./components/locations/status.vue').default

/*----------  Onboarding scope  ----------*/

const OnboardScope = require('./components/locations/onboardscope.vue').default
const OnboardRoot = require('./components/locations/onboardscope/onboardroot.vue')
  .default
const Onboard1 = require('./components/locations/onboardscope/onboard1.vue')
  .default
const Onboard2 = require('./components/locations/onboardscope/onboard2.vue')
  .default
const Onboard3 = require('./components/locations/onboardscope/onboard3.vue')
  .default
const Onboard4 = require('./components/locations/onboardscope/onboard4.vue')
  .default
const Onboard5 = require('./components/locations/onboardscope/onboard5.vue')
  .default
const Onboard6 = require('./components/locations/onboardscope/onboard6.vue')
  .default

/*----------  Search scope (search communities, content, users)  ----------*/
const SearchScope = require('./components/locations/searchscope.vue').default
const SearchCommunity = require('./components/locations/searchscope/communitysearch.vue')
  .default
const SearchContent = require('./components/locations/searchscope/contentsearch.vue')
  .default
const SearchUser = require('./components/locations/searchscope/usersearch.vue')
  .default
/*----------  Routes  ----------*/

const routes = [
  { path: '/', component: Home, name: 'Home' },
  { path: '/popular', component: Popular, name: 'Popular' },
  {
    path: '/globalscope',
    component: GlobalScope,
    children: [
      { path: '', component: GlobalRoot, name: 'Global' },
      {
        path: '/globalscope/subbed',
        component: GlobalSubbed,
        name: 'Global>Subbed',
      },
      {
        path: '/globalscope/newboard',
        component: NewBoard,
        name: 'Global>NewBoard',
      },
    ],
  },
  {
    path: '/searchscope',
    component: SearchScope,
    children: [
      { path: '', component: SearchCommunity, name: 'Search' },
      // ^ SearchCommunity is the root page
      {
        path: '/searchscope/content',
        component: SearchContent,
        name: 'Search>Content',
      },
      {
        path: '/searchscope/user',
        component: SearchUser,
        name: 'Search>People',
      },
    ],
  },
  {
    path: '/board/:boardfp',
    component: BoardScope,
    children: [
      { path: '', component: BoardRoot, name: 'Board' },
      {
        path: '/board/:boardfp/new',
        component: BoardRoot,
        name: 'Board>ThreadsNewList',
      },
      {
        path: '/board/:boardfp/info',
        component: BoardInfo,
        name: 'Board>BoardInfo',
      },
      {
        path: '/board/:boardfp/modactivity',
        component: ModActivity,
        name: 'Board>ModActivity',
      },
      {
        path: '/board/:boardfp/elections',
        component: Elections,
        name: 'Board>Elections',
      },
      {
        path: '/board/:boardfp/newthread',
        component: NewThread,
        name: 'Board>NewThread',
      },
      {
        path: '/board/:boardfp/reports',
        component: Reports,
        name: 'Board>Reports',
      },
    ],
  },
  {
    path: '/board/:boardfp/thread/:threadfp',
    component: ThreadScope,
    name: 'Thread',
    props: function(route: any) {
      let highlightSelectors: any = []
      if (
        !globalMethods.IsUndefined(route.query.highlightSelectors) &&
        route.query.highlightSelectors.length > 0
      ) {
        highlightSelectors = JSON.parse(route.query.highlightSelectors)
      }
      return {
        route_focusSelector: route.query.focusSelector,
        route_parentSelector: route.query.parentSelector,
        route_highlightSelectors: highlightSelectors,
      }
    },
  },
  {
    path: '/settings',
    component: SettingsScope,
    children: [
      { path: '', component: SettingsRoot, name: 'Settings' },
      {
        path: '/settings/defaults',
        component: Defaults,
        name: 'Settings>Defaults',
      },
      {
        path: '/settings/shortcuts',
        component: Shortcuts,
        name: 'Settings>Shortcuts',
      },
      {
        path: '/settings/advanced',
        component: AdvancedSettings,
        name: 'Settings>Advanced',
      },
      /*This is a little weird, these things are in settings scope but they're not in a settings path. That's because they exist in a router link that is in the settings structure. If you move this outside and try to use it, it uses the router link outside settings, which is the main main-block router link, which means the settings frame box won't be rendered. So this is not an oversight. */
      { path: '/intro', component: Intro, name: 'Intro' },
      { path: '/about', component: About, name: 'About' },
      { path: '/membership', component: Membership, name: 'Membership' },
      { path: '/changelog', component: Changelog, name: 'Changelog' },
      {
        path: '/adminsquickstart',
        component: AdminsQuickstart,
        name: 'AdminsQuickstart',
      },
      { path: '/newuser', component: NewUser, name: 'NewUser' },
      { path: '/sfwlist', component: SFWList, name: 'SFWList' },
      { path: '/modship', component: Modship, name: 'Modship' },
      { path: '/namemint', component: Namemint, name: 'Namemint' },
    ],
  },
  {
    path: '/user/:userfp',
    component: UserScope,
    children: [
      { path: '', component: UserRoot, name: 'User' },
      {
        path: '/user/:userfp/boards',
        component: UserBoards,
        name: 'User>Boards',
      },
      {
        path: '/user/:userfp/threads',
        component: UserThreads,
        name: 'User>Threads',
      },
      { path: '/user/:userfp/posts', component: UserPosts, name: 'User>Posts' },
      {
        path: '/user/:userfp/notifications',
        component: Notifications,
        name: 'User>Notifications',
      },
      { path: '*', redirect: '/user/:userfp' },
    ],
  },
  { path: '/status', component: Status, name: 'Status' },
  {
    path: '/onboard',
    components: { default: '', onboarding: OnboardScope },
    children: [
      { path: '', component: OnboardRoot, name: 'OnboardRoot' },
      { path: '/onboard/1', component: Onboard1, name: 'Onboard1' },
      { path: '/onboard/2', component: Onboard2, name: 'Onboard2' },
      { path: '/onboard/3', component: Onboard3, name: 'Onboard3' },
      { path: '/onboard/4', component: Onboard4, name: 'Onboard4' },
      { path: '/onboard/5', component: Onboard5, name: 'Onboard5' },
      { path: '/onboard/6', component: Onboard6, name: 'Onboard6' },
    ],
  },
  { path: '*', redirect: '/' },
]

// { path: '/user/:userfp/posts', component: UserPosts, name: 'User>Posts', },
// { path: '/user/:userfp/threads', component: UserThreads, name: 'User>Threads', },

/*----------  Plumbing  ----------*/

const router = new VueRouter({
  scrollBehavior() {
    // Always return to top while navigating.
    return { x: 0, y: 0 }
  },
  // ^ This does not work because we are using a fixed container and scroll inside it. Attempting to do it like this attempts to scroll the main container, which does not scroll. There is no way to specify which container needs to be scrolled, so we need to implement our own scroll behaviour.
  routes: routes,
  // mode: 'history'
})

// This keeps track of history, so we can appropriately disable back / forward buttons as needed.
// router.afterEach(HistoryWriter)

const Store = require('./store').default

new Vue({
  el: '#app',
  template: '<a-app></a-app>',
  router: router,
  store: Store,
  mounted(this: any) {
    ipc.callMain('SetRendererReady', true)
  },
})

let Sync = require('../../node_modules/vuex-router-sync').sync
Sync(Store, router)
/*
^ It adds a route module into the store, which contains the state representing the current route:
store.state.route.path   // current path (string)
store.state.route.params // current params (object)
store.state.route.query  // current query (object)
*/

// Disable events that are meaningless in this context.

// Drag start is being able to click and drag a link inside the app to outside of it. Since the app is a local one, that link will just be a local file, and it won't be useful to anybody.
document.addEventListener('dragstart', function(event: any) {
  event.preventDefault()
})

// Dragover is the event that gets fired when a dragged item is on a droppable target, every few hundred milliseconds. We have no drop targets.
document.addEventListener('dragover', function(event: any) {
  event.preventDefault()
})

// Cancelling drop prevents anything from being dropped into the container. This can be a mild security risk, if someone can convince you (or somehow automate dropping inside the app container), it can make the container ping a web address. This also assumes the container has the dropped remote address whitelisted, though, so it's a long shot. Still, defence in depth is preferable.
document.addEventListener('drop', function(event: any) {
  event.preventDefault()
})

/*----------  Some basic keyboard shortcuts  ----------*/

Mousetrap.bind('mod+,', function() {
  history.back()
  // if (event.target.nodeName.toLowerCase() !== 'textarea' && event.target.nodeName.toLowerCase() !== 'input' && event.target.contentEditable !== 'true') {
  //   history.back()
  // }
})

Mousetrap.bind('mod+.', function() {
  history.forward()
  // if (event.target.nodeName.toLowerCase() !== 'textarea' && event.target.nodeName.toLowerCase() !== 'input' && event.target.contentEditable !== 'true') {
  //   history.forward()
  // }
})

Mousetrap.bind('mod+/', function() {
  router.push('/user/' + Store.state.localUser.fingerprint + '/notifications')
})

/*----------  IPC maps  ----------*/

/*
These are here instead of eipc/eipc-renderer because they do require access to things that are instantiated here, such as router, and there is no way to get to them without importing the main. Importing main is not an option. So these should be here until I split the router into its own service file that is imported separately. That way, eipc import from there, and not from main.
*/
ipc.answerMain('RouteTo', function(route: string) {
  router.push(route)
  return
})

ipc.answerMain('FullscreenState', function(isFullscreen: boolean) {
  Store.state.appIsFullscreen = isFullscreen
})

ipc.answerMain('NewUpdateReady', function(newUpdateReady: boolean) {
  Store.state.newUpdateReady = newUpdateReady
})

/*----------  Exports  ----------*/

module.exports = { router: router }

/*========================================
=            Spell checker             =
========================================*/
const HunspellAsm = require('hunspell-asm')
const ElectronHunspell = require('electron-hunspell')
const path = require('path')

ElectronHunspell.enableLogger(console)

const init = async () => {
  const browserWindowProvider = new ElectronHunspell.SpellCheckerProvider()
    ; (window as any).browserWindowProvider = browserWindowProvider
  await browserWindowProvider.initialize({
    environment: HunspellAsm.ENVIRONMENT.NODE,
    locateBinary: function(file: any) {
      if (file.endsWith('.wasm')) {
        return 'node_modules/hunspell-asm/dist/cjs/lib/hunspell.wasm'
      }
      return file
    },
  })
  if (!isDev) {
    await browserWindowProvider.loadDictionary(
      'en',
      path.join(__dirname, '../app/dicts/en-US.dic'),
      path.join(__dirname, '../app/dicts/en-US.aff')
    )
  } else {
    await browserWindowProvider.loadDictionary(
      'en',
      'src/app/ext_dep/dicts/en-US.dic',
      'src/app/ext_dep/dicts/en-US.aff'
    )
  }

  browserWindowProvider.switchDictionary('en')
}
init()
/*=====  End of Spell checker   ======*/

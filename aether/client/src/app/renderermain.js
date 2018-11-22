"use strict";
/*
This is the main entry point to the client app. See app.vue for the start logic, and globally-applicable css.
*/
Object.defineProperty(exports, "__esModule", { value: true });
// Electron IPC setup before doing anything else
require('./services/eipc/eipc-renderer'); // Register IPC events
var ipc = require('../../node_modules/electron-better-ipc'); // Register IPC caller
// ^ Heads up, there are some IPC events registered in this renderermain, too.
// const unhandled = require('../../node_modules/electron-unhandled')
// unhandled()
var clapiserver = require('./services/clapiserver/clapiserver');
var feapiconsumer = require('./services/feapiconsumer/feapiconsumer');
var globalMethods = require('./services/globals/methods');
var clientAPIServerPort = clapiserver.StartClientAPIServer();
console.log('attempting to call get frontend ready');
ipc.callMain('GetFrontendReady').then(function (resp) {
    console.log('frontend ready response received');
    console.log(resp);
}).catch(function (err) {
    console.log('fe ready promise rejected.');
    console.log(err);
});
/*----------  Call mainmain to ask software update state  ----------*/
ipc.callMain('AskNewUpdateReady');
console.log('renderer client api server port: ', clientAPIServerPort);
ipc.callMain('SetClientAPIServerPort', clientAPIServerPort).then(function (feDaemonStarted) {
    if (!feDaemonStarted) {
        // It's an Electron refresh, not a cold start.
        feapiconsumer.Initialise();
    }
}).catch(function (err) {
    console.log('this is the promise error:');
    console.log(err);
});
/*----------  Vue + its plugins  ----------*/
var Vue = require('../../node_modules/vue/dist/vue.js');
var VueRouter = require('../../node_modules/vue-router').default;
Vue.use(VueRouter);
// Register icons for our own use.
var Icon = require('../../node_modules/vue-awesome');
Vue.component('icon', Icon);
// Register the click-outside component
var ClickOutside = require('../../node_modules/v-click-outside');
Vue.use(ClickOutside);
/*----------  Third party dependencies  ----------*/
var Mousetrap = require('../../node_modules/mousetrap');
// var Spinner = require('../../node_modules/vue-simple-spinner')
/*----------  Components  ----------*/
// Global component declarations - do it here once.
Vue.component('a-app', require('./components/a-app.vue').default);
Vue.component('a-header', require('./components/a-header.vue').default);
Vue.component('a-header-icon', require('./components/a-header-icon.vue').default);
Vue.component('a-sidebar', require('./components/a-sidebar.vue').default);
Vue.component('a-boardheader', require('./components/a-boardheader.vue').default);
Vue.component('a-tabs', require('./components/a-tabs.vue').default);
Vue.component('a-thread-entity', require('./components/a-thread-entity.vue').default);
Vue.component('a-vote-action', require('./components/a-vote-action.vue').default);
Vue.component('a-thread-header-entity', require('./components/a-thread-header-entity.vue').default);
Vue.component('a-post', require('./components/a-post.vue').default);
Vue.component('a-side-header', require('./components/a-side-header.vue').default);
Vue.component('a-breadcrumbs', require('./components/a-breadcrumbs.vue').default);
Vue.component('a-username', require('./components/a-username.vue').default);
Vue.component('a-timestamp', require('./components/a-timestamp.vue').default);
Vue.component('a-globalscopeheader', require('./components/a-globalscopeheader.vue').default);
Vue.component('a-board-entity', require('./components/a-board-entity.vue').default);
Vue.component('a-hashimage', require('./components/a-hashimage.vue').default);
Vue.component('a-no-content', require('./components/a-no-content.vue').default);
Vue.component('a-markdown', require('./components/a-markdown.vue').default);
Vue.component('a-avatar-block', require('./components/a-avatar-block.vue').default);
Vue.component('a-composer', require('./components/a-composer.vue').default);
Vue.component('a-ballot', require('./components/a-ballot.vue').default);
Vue.component('a-progress-bar', require('./components/a-progress-bar.vue').default);
Vue.component('a-inflight-info', require('./components/a-inflight-info.vue').default);
Vue.component('a-info-marker', require('./components/a-info-marker.vue').default);
Vue.component('a-spinner', require('./components/a-spinner.vue').default);
Vue.component('a-notfound', require('./components/a-notfound.vue').default);
Vue.component('a-guidelight', require('./components/a-guidelight.vue').default);
Vue.component('a-home-header', require('./components/a-home-header.vue').default);
Vue.component('a-popular-header', require('./components/a-popular-header.vue').default);
Vue.component('a-notifications-icon', require('./components/a-notifications-icon.vue').default);
Vue.component('a-notification-entity', require('./components/a-notification-entity.vue').default);
Vue.component('a-main-app-loader', require('./components/a-main-app-loader.vue').default);
Vue.component('a-global-header', require('./components/a-global-header.vue').default);
Vue.component('a-fin-puck', require('./components/a-fin-puck.vue').default);
Vue.component('a-bootstrapper', require('./components/a-bootstrapper.vue').default);
Vue.component('a-fingerprint', require('./components/a-fingerprint.vue').default);
Vue.component('a-settings-block', require('./components/a-settings-block.vue').default);
Vue.component('a-software-update-icon', require('./components/a-software-update-icon.vue').default);
Vue.component('a-patreon-button', require('./components/a-patreon-button.vue').default);
Vue.component('a-crypto-fund-button', require('./components/a-crypto-fund-button.vue').default);
Vue.component('a-boardname', require('./components/a-boardname.vue').default);
/*----------  Third party components  ----------*/
Vue.component('vue-simple-spinner', require('../../node_modules/vue-simple-spinner'));
/*----------  Places  ----------*/
var Home = require('./components/locations/home.vue').default;
var Popular = require('./components/locations/popular.vue').default;
/*----------  Global scope (whole network, i.e. list of boards)  ----------*/
var GlobalScope = require('./components/locations/globalscope.vue').default;
var GlobalRoot = require('./components/locations/globalscope/globalroot.vue').default;
var GlobalSubbed = require('./components/locations/globalscope/subbedroot.vue').default;
/*----------  Board scope (board entity + list of threads)  ----------*/
var NewBoard = require('./components/locations/globalscope/newboard.vue').default;
var BoardScope = require('./components/locations/boardscope.vue').default;
var BoardRoot = require('./components/locations/boardscope/boardroot.vue').default;
var BoardInfo = require('./components/locations/boardscope/boardinfo.vue').default;
var ModActivity = require('./components/locations/boardscope/modactivity.vue').default;
var Elections = require('./components/locations/boardscope/elections.vue').default;
var Reports = require('./components/locations/boardscope/reports.vue').default;
/*----------  Thread scope (thread entity + list of posts)  ----------*/
var NewThread = require('./components/locations/boardscope/newthread.vue').default;
var ThreadScope = require('./components/locations/threadscope.vue').default;
/*----------  Settings scope  ----------*/
var SettingsScope = require('./components/locations/settingsscope.vue').default;
var SettingsRoot = require('./components/locations/settingsscope/settingsroot.vue').default;
var AdvancedSettings = require('./components/locations/settingsscope/advancedsettings.vue').default;
var About = require('./components/locations/settingsscope/about.vue').default;
var Membership = require('./components/locations/settingsscope/membership.vue').default;
var Changelog = require('./components/locations/settingsscope/changelog.vue').default;
var AdminsQuickstart = require('./components/locations/settingsscope/adminsquickstart.vue').default;
var Intro = require('./components/locations/settingsscope/intro.vue').default;
var NewUser = require('./components/locations/settingsscope/newuser.vue').default;
var SFWList = require('./components/locations/settingsscope/sfwlist.vue').default;
var Modship = require('./components/locations/settingsscope/modship.vue').default;
var Namemint = require('./components/locations/settingsscope/namemint.vue').default;
/*----------  User scope  ----------*/
var UserScope = require('./components/locations/userscope.vue').default;
var UserRoot = require('./components/locations/userscope/userroot.vue').default;
var UserBoards = require('./components/locations/userscope/userboards.vue').default;
var UserThreads = require('./components/locations/userscope/userthreads.vue').default;
var UserPosts = require('./components/locations/userscope/userposts.vue').default;
var Notifications = require('./components/locations/userscope/notifications.vue').default;
/*----------  Status scope  ----------*/
var Status = require('./components/locations/status.vue').default;
/*----------  Onboarding scope  ----------*/
var OnboardScope = require('./components/locations/onboardscope.vue').default;
var OnboardRoot = require('./components/locations/onboardscope/onboardroot.vue').default;
var Onboard1 = require('./components/locations/onboardscope/onboard1.vue').default;
var Onboard2 = require('./components/locations/onboardscope/onboard2.vue').default;
var Onboard3 = require('./components/locations/onboardscope/onboard3.vue').default;
var Onboard4 = require('./components/locations/onboardscope/onboard4.vue').default;
var Onboard5 = require('./components/locations/onboardscope/onboard5.vue').default;
var Onboard6 = require('./components/locations/onboardscope/onboard6.vue').default;
/*----------  Routes  ----------*/
var routes = [
    { path: '/', component: Home, name: 'Home', },
    { path: '/popular', component: Popular, name: 'Popular', },
    {
        path: '/globalscope', component: GlobalScope,
        children: [
            { path: '', component: GlobalRoot, name: 'Global', },
            { path: '/globalscope/subbed', component: GlobalSubbed, name: 'Global>Subbed', },
            { path: '/globalscope/newboard', component: NewBoard, name: 'Global>NewBoard', },
        ]
    },
    {
        path: '/board/:boardfp', component: BoardScope,
        children: [
            { path: '', component: BoardRoot, name: 'Board', },
            { path: '/board/:boardfp/new', component: BoardRoot, name: 'Board>ThreadsNewList', },
            { path: '/board/:boardfp/info', component: BoardInfo, name: 'Board>BoardInfo', },
            { path: '/board/:boardfp/modactivity', component: ModActivity, name: 'Board>ModActivity', },
            { path: '/board/:boardfp/elections', component: Elections, name: 'Board>Elections', },
            { path: '/board/:boardfp/newthread', component: NewThread, name: 'Board>NewThread', },
            { path: '/board/:boardfp/reports', component: Reports, name: 'Board>Reports', },
        ]
    }, {
        path: '/board/:boardfp/thread/:threadfp', component: ThreadScope, name: 'Thread', props: function (route) {
            var highlightSelectors = [];
            if (!globalMethods.IsUndefined(route.query.highlightSelectors) && route.query.highlightSelectors.length > 0) {
                highlightSelectors = JSON.parse(route.query.highlightSelectors);
            }
            return {
                route_focusSelector: route.query.focusSelector,
                route_parentSelector: route.query.parentSelector,
                route_highlightSelectors: highlightSelectors,
            };
        }
    },
    {
        path: '/settings', component: SettingsScope,
        children: [
            { path: '', component: SettingsRoot, name: 'Settings', },
            { path: '/settings/advanced', component: AdvancedSettings, name: 'Settings>Advanced', },
            /*This is a little weird, these things are in settings scope but they're not in a settings path. That's because they exist in a router link that is in the settings structure. If you move this outside and try to use it, it uses the router link outside settings, which is the main main-block router link, which means the settings frame box won't be rendered. So this is not an oversight. */
            { path: '/intro', component: Intro, name: 'Intro', },
            { path: '/about', component: About, name: 'About', },
            { path: '/membership', component: Membership, name: 'Membership', },
            { path: '/changelog', component: Changelog, name: 'Changelog', },
            { path: '/adminsquickstart', component: AdminsQuickstart, name: 'AdminsQuickstart', },
            { path: '/newuser', component: NewUser, name: 'NewUser', },
            { path: '/sfwlist', component: SFWList, name: 'SFWList', },
            { path: '/modship', component: Modship, name: 'Modship', },
            { path: '/namemint', component: Namemint, name: 'Namemint', },
        ]
    },
    {
        path: '/user/:userfp', component: UserScope,
        children: [
            { path: '', component: UserRoot, name: 'User' },
            { path: '/user/:userfp/boards', component: UserBoards, name: 'User>Boards' },
            { path: '/user/:userfp/threads', component: UserThreads, name: 'User>Threads' },
            { path: '/user/:userfp/posts', component: UserPosts, name: 'User>Posts' },
            { path: '/user/:userfp/notifications', component: Notifications, name: 'User>Notifications' },
            { path: '*', redirect: '/user/:userfp' }
        ]
    },
    { path: '/status', component: Status, name: 'Status', },
    {
        path: '/onboard', components: { default: '', onboarding: OnboardScope },
        children: [
            { path: '', component: OnboardRoot, name: 'OnboardRoot', },
            { path: '/onboard/1', component: Onboard1, name: 'Onboard1', },
            { path: '/onboard/2', component: Onboard2, name: 'Onboard2', },
            { path: '/onboard/3', component: Onboard3, name: 'Onboard3', },
            { path: '/onboard/4', component: Onboard4, name: 'Onboard4', },
            { path: '/onboard/5', component: Onboard5, name: 'Onboard5', },
            { path: '/onboard/6', component: Onboard6, name: 'Onboard6', },
        ]
    },
    { path: '*', redirect: '/' }
];
// { path: '/user/:userfp/posts', component: UserPosts, name: 'User>Posts', },
// { path: '/user/:userfp/threads', component: UserThreads, name: 'User>Threads', },
/*----------  Plumbing  ----------*/
var router = new VueRouter({
    scrollBehavior: function () {
        return { x: 0, y: 0 };
    },
    // ^ This does not work because we are using a fixed container and scroll inside it. Attempting to do it like this attempts to scroll the main container, which does not scroll. There is no way to specify which container needs to be scrolled, so we need to implement our own scroll behaviour.
    routes: routes,
});
// This keeps track of history, so we can appropriately disable back / forward buttons as needed.
// router.afterEach(HistoryWriter)
var Store = require('./store').default;
new Vue({
    el: '#app',
    template: '<a-app></a-app>',
    router: router,
    store: Store,
    mounted: function () {
        ipc.callMain('SetRendererReady', true);
    }
});
var Sync = require('../../node_modules/vuex-router-sync').sync;
Sync(Store, router);
/*
^ It adds a route module into the store, which contains the state representing the current route:
store.state.route.path   // current path (string)
store.state.route.params // current params (object)
store.state.route.query  // current query (object)
*/
// Disable events that are meaningless in this context.
// Drag start is being able to click and drag a link inside the app to outside of it. Since the app is a local one, that link will just be a local file, and it won't be useful to anybody.
document.addEventListener('dragstart', function (event) { event.preventDefault(); });
// Dragover is the event that gets fired when a dragged item is on a droppable target, every few hundred milliseconds. We have no drop targets.
document.addEventListener('dragover', function (event) { event.preventDefault(); });
// Cancelling drop prevents anything from being dropped into the container. This can be a mild security risk, if someone can convince you (or somehow automate dropping inside the app container), it can make the container ping a web address. This also assumes the container has the dropped remote address whitelisted, though, so it's a long shot. Still, defence in depth is preferable.
document.addEventListener('drop', function (event) { event.preventDefault(); });
/*----------  Some basic keyboard shortcuts  ----------*/
Mousetrap.bind('mod+,', function () {
    history.back();
    // if (event.target.nodeName.toLowerCase() !== 'textarea' && event.target.nodeName.toLowerCase() !== 'input' && event.target.contentEditable !== 'true') {
    //   history.back()
    // }
});
Mousetrap.bind('mod+.', function () {
    history.forward();
    // if (event.target.nodeName.toLowerCase() !== 'textarea' && event.target.nodeName.toLowerCase() !== 'input' && event.target.contentEditable !== 'true') {
    //   history.forward()
    // }
});
/*----------  IPC maps  ----------*/
/*
These are here instead of eipc/eipc-renderer because they do require access to things that are instantiated here, such as router, and there is no way to get to them without importing the main. Importing main is not an option. So these should be here until I split the router into its own service file that is imported separately. That way, eipc import from there, and not from main.
*/
ipc.answerMain('RouteTo', function (route) {
    router.push(route);
    return;
});
ipc.answerMain('FullscreenState', function (isFullscreen) {
    Store.state.appIsFullscreen = isFullscreen;
});
ipc.answerMain('NewUpdateReady', function (newUpdateReady) {
    Store.state.newUpdateReady = newUpdateReady;
});
/*----------  Exports  ----------*/
module.exports = { router: router };
//# sourceMappingURL=renderermain.js.map
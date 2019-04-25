// Services > Tooltips

// This service provides a tooltip binder - for any view that has tooltips, this needs to run.

export {}
var Tippy: any = require('../../../../node_modules/tippy.js')

module.exports = {
  Mount() {
    Tippy('[hasTooltip]', {
      animation: 'fade',
      delay: 100,
      duration: [100, 500], // [show, hide]
      performance: true,
      // dynamicTitle: true,
      placement: 'top',
      theme: 'dark',
      offset: '0,-1',
      interactiveBorder: 10,
      // arrow: true,
      // arrowType: 'round',
      // size: 'small',
      // createPopperInstanceOnInit: true,
      // inertia: true
    })
  },
  MountInfomark() {
    Tippy('[hasInfomark]', {
      animation: 'fade',
      delay: 50,
      duration: [200, 500], // [show, hide]
      performance: true,
      // dynamicTitle: true,
      placement: 'bottom',
      theme: 'infomark',
      // size: 'small',
      // inertia: true
      arrow: true,
      arrowType: 'round',
      hideOnClick: true,
      interactive: true,
      allowTitleHTML: true,
      interactiveBorder: 10,
      offset: '0,5',
      // trigger: 'click',
    })
  },
  MountGuidelightTooltip() {
    Tippy('[hasGuidelightTooltip]', {
      animation: 'fade',
      delay: 50,
      duration: [200, 500], // [show, hide]
      performance: true,
      // dynamicTitle: true,
      placement: 'bottom',
      theme: 'infomark',
      // size: 'small',
      // inertia: true
      arrow: true,
      arrowType: 'round',
      hideOnClick: false,
      interactive: true,
      allowTitleHTML: true,
      interactiveBorder: 10,
      offset: '0,5',
      // trigger: 'click',
    })
  },
  MountMarker() {
    Tippy('[hasMarker]', {
      animation: 'fade',
      delay: 500, // These are easter eggs. Don't be obnoxious.
      duration: [100, 500], // [show, hide]
      performance: true,
      // dynamicTitle: true,
      // placement: 'top',
      theme: 'marker',
      offset: '0,0',
      interactiveBorder: 10,
      // arrow: true,
      // arrowType: 'round',
      // size: 'small',
      // createPopperInstanceOnInit: true,
      // inertia: true
    })
  },
  MountUsernameTooltip() {
    Tippy('[hasUsernameTooltip]', {
      animation: 'fade',
      delay: 50,
      duration: [200, 500], // [show, hide]
      performance: true,
      // dynamicTitle: true,
      placement: 'bottom',
      theme: 'usernametooltip',
      // size: 'small',
      // inertia: true
      arrow: true,
      arrowType: 'round',
      hideOnClick: true,
      interactive: true,
      allowTitleHTML: true,
      interactiveBorder: 10,
      offset: '0,5',
      // trigger: 'click',
    })
  },
}

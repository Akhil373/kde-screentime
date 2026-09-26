print("=== SCREEN TIME SCRIPT LOADED ===");

workspace.windowActivated.connect(function (window) {
  if (!window) return;

  var data = {
    pid: window.pid || 0,
    caption: window.caption || "",
    wm_class: window.resourceClass || "",
    virtual_desktop: window.desktops.map(function (desktop) {
      return desktop.x11DesktopNumber;
    }),
  };

  callDBus(
    "org.screentime",
    "/org/screentime/Tracker",
    "org.screentime.Tracker",
    "WindowActivated",
    JSON.stringify(data),
  );
});

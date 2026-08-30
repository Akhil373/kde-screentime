# kde plasma wayland screen time tracker


a simple screen time tracker for kde plasma wayland, written in go.


wayland restricts window metadata access, so a kwin script hooks into window activation events and sends active window metadata over dbus to the go daemon, which stores activity intervals in sqlite.


data analytics are still a work in progress.


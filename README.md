# kde plasma wayland screen time tracker


a simple screen time tracker for kde plasma wayland, written in go.


wayland restricts window metadata access, so a kwin script hooks into window activation events and sends active window metadata over dbus to the go daemon, which stores activity intervals in sqlite.


<img width="790" height="450" alt="image" src="https://github.com/user-attachments/assets/88f99eec-4aad-44c7-8f8d-030b2754b5da" />

data analytics are still a work in progress.


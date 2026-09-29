# kde plasma wayland screen time tracker


a simple screen time tracker for kde plasma wayland, written in go.


wayland restricts window metadata access, so a kwin script hooks into window activation events and sends active window metadata over dbus to the go daemon, which stores activity intervals in sqlite.

<img src="https://github.com/user-attachments/assets/115c53ae-2a66-4a17-9e4a-6bb620b5a6d3" alt="Weekly Screen Time" width="75%">
<img src="https://github.com/user-attachments/assets/d7ef1e45-7584-462a-9765-274de8dac376" alt="App Usage Breakdown" width="75%">


data analytics are still a work in progress.


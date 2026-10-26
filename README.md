# kde plasma wayland screen time tracker


a simple screen time tracker for kde plasma wayland, written in go.


wayland restricts window metadata access, so a kwin script hooks into window activation events and sends active window metadata over dbus to the go daemon, which stores activity intervals in sqlite.

<img src="https://github.com/user-attachments/assets/115c53ae-2a66-4a17-9e4a-6bb620b5a6d3" alt="Weekly Screen Time" width="75%">
<img src="https://github.com/user-attachments/assets/d7ef1e45-7584-462a-9765-274de8dac376" alt="App Usage Breakdown" width="75%">


data analytics are still a work in progress.

### Installation

Get the KWin script and `scrntime-linux-amd64` binary from the [Releases](../../releases) page, or build them yourself with the Makefile.

#### 1. Install the KWin script

Go to:

**System Settings → Window Management → KWin Scripts → Install from File**

Select `kde-screentime-tracker.kwinscript` and enable the script.

#### 2. Install `scrntime`

Rename the downloaded binary to `scrntime` and place it somewhere in your `PATH`, for example:

```bash
mkdir -p ~/.local/bin
mv scrntime-linux-amd64 ~/.local/bin/scrntime
chmod +x ~/.local/bin/scrntime
```

Make sure `~/.local/bin` is in your `PATH`.

#### 3. Start `scrntime` automatically

The `scrntime` binary needs to be running in the background for tracking to work.

You can add it to KDE's autostart:

**System Settings → Autostart → Add Application**

Select `scrntime`.

Alternatively, you can create a user systemd service if you prefer.

After that, the tracker should start automatically when you log in.

#### 4. Generate a report

Run:

```bash
scrntime -g
```

This generates an HTML report containing your screen-time data.

> This is a personal project and is not packaged as a polished desktop application yet, so installation currently involves a few manual steps.

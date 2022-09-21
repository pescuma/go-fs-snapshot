# go-fs-snapshot

Allows creation of file system snapshots.

Windows version based on [restic](https://github.com/restic/restic) source code: https://github.com/restic/restic/blob/master/internal/fs/vss_windows.go

## Platforms

### Windows

Well supported. Needs to run as Adminstrator, but also supports running a server as administrator and a client with normal privileges. To enable that, run as administrator:
```
fs_snapshot enable for current-user
```
After that you can run it normally and it should work. That command registers a task in Task Scheduler that runs as Administrator and allows communication throught RPC. The server runs while a backup is made and stops after 5 mins without a client connected.


### MacOS

Initial support implemented. Listing of volumes still not supported, so it assumes that there is only one APFS volume. To run, you need to give Full Disk Access permission to the executable or to `Terminal.app`. This is explained in `fs_snapshot enable for current-user`:
> MacOS does not allow to grant Full Disk Access permission from an application. You need to open 'System Preferences...', go to the 'Privacy' tab, select 'Full Disk Access' in the list on the left, click on the lock on the bottom, input your password and then add the correct application to the list on the right. If you intend to use this app inside terminal, you must select 'Terminal.app' in the list on the right (for some reason granting the permission to fs_snapshot does not work). In some other cases you may need to add and grant the permission to 'fs_snapshot'.


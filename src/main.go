package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

/*
#include <windows.h>

int GetSerivceProcessId(HANDLE handle) {
	SERVICE_STATUS_PROCESS ssStatus;
	DWORD dwBytesNeeded;

    if (!QueryServiceStatusEx(
		handle,                          // handle to service
		SC_STATUS_PROCESS_INFO,          // information level
		(LPBYTE)&ssStatus,               // address of structure
		sizeof(SERVICE_STATUS_PROCESS),  // size of structure
		&dwBytesNeeded))                 // size needed if buffer is too small
	{
		return -1;
	}

    return ssStatus.dwProcessId;
}
*/
import "C"

func checkErr(err error, msg string) {
	if err != nil {
		fmt.Println(msg, " Error: ", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: StopSrv.exe {service name}")
		os.Exit(0)
	}

	srvName := os.Args[1]
	srvMgr, err := mgr.Connect()
	checkErr(err, "Connect service manager")

	srv, err := srvMgr.OpenService(srvName)
	checkErr(err, "Open service")

	start := time.Now()
	for {
		status, err := srv.Query()
		checkErr(err, "Query service status")
		if status.State == svc.Running {

			_, err := srv.Control(svc.Stop)
			checkErr(err, "Stop service")

		} else if status.State == svc.StopPending {
			if time.Now().Sub(start) > 10*time.Second {
				checkErr(fmt.Errorf("Timeout"), "Wait Stopped")

			} else if time.Now().Sub(start) > 5*time.Second {
				processId := C.GetSerivceProcessId(C.HANDLE(srv.Handle))
				if processId <= 0 {
					checkErr(fmt.Errorf("Error pid: %d", processId),
						"Get service processId")
				}

				process, err := os.FindProcess(int(processId))
				checkErr(err, "Find process")

				err = process.Kill()
				checkErr(err, "Kill process")
			}
		} else {
			break
		}

		time.Sleep(time.Second)
	}

	os.Exit(0)
}

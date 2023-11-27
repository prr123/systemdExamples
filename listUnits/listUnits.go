// https://www.reddit.com/r/golang/comments/x5t31x/help_with_systemd_unit_status_using_dbus/

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"bytes"
//	"reflect"

	"github.com/coreos/go-systemd/v22/dbus"
)

func getUnitStatus(ctx context.Context, servList [][]byte) ([]dbus.UnitStatus, error) {
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	//changed
//	sliceOfUnits := []string{"systemd-networkd.service"}
	sliceOfUnits := make([]string, len(servList))
	for i:=0; i< len(servList); i++ {
		sliceOfUnits[i] = string(servList[i])
	}

// alt ByNamesContext(ctx context.Context, ["mariadb.service"])

// https://pkg.go.dev/github.com/coreos/go-systemd/v22/dbus#Conn.ListUnitsByNamesContext
	sliceOfUnitStatus, err := conn.ListUnitsByNamesContext(ctx, sliceOfUnits) 

	return sliceOfUnitStatus, err
}


func PrintUnitStatus (statusList []dbus.UnitStatus) {

	fmt.Printf("********** Units: %d *********\n", len(statusList))
	for i:=0; i< len(statusList); i++ {
		st := statusList[i]
		fmt.Printf("  ******** Unit: %d *********\n", i)
		fmt.Printf("  Name:        %s\n", st.Name)
		fmt.Printf("  Description: %s\n", st.Description)
		fmt.Printf("  LoadState:   %s\n", st.LoadState)
		fmt.Printf("  ActiveState: %s\n", st.ActiveState)
		fmt.Printf("  SubState:    %s\n", st.SubState)
		fmt.Printf("  Followed:    %s\n", st.Followed)
		fmt.Printf("  Path:        %s\n", st.Path)
		fmt.Printf("  JobId:       %d\n", st.JobId)
		fmt.Printf("  JobType:     %s\n", st.JobType)
		fmt.Printf("  JobPath:     %s\n", st.JobPath)
	}
	fmt.Printf("********** End Units *********\n")
}

func main() {

	numArgs := len(os.Args)
	useStr := "./listUnits [service file1] [service file 2] ... [service file n]"
	helpStr := "program that lists the status of systemd service units!"

	if numArgs == 1 {
		fmt.Printf("need a service file as argument!\n")
		fmt.Printf("usage: %s\n", useStr)
	}
	if numArgs == 2 {
		if os.Args[1] == "help" {
			fmt.Printf("help: %s\n", helpStr)
			fmt.Printf("usage: %s\n", useStr)
			os.Exit(0)
		}
	}

	servList := make([][]byte, numArgs-1)
	for i:=1; i< numArgs; i++ {
		servList[i-1] = []byte(os.Args[i])
	}

	for i, nam := range servList {
		namlen := len(nam)
//		fmt.Printf(" %q\n", nam[namlen-1])
		if nam[namlen-1] == ',' {
			nam = nam[:namlen-1]
			servList[i] = nam
		}
//		fmt.Printf(" %s\n", nam)

		idx := bytes.Index(nam, []byte(".service"))
		if idx == -1 {
			servList[i] = append(nam, []byte(".service")...)
		}
	}

	for i, nam := range servList {
		fmt.Printf("  %d: %s\n", i, nam)
	}

//	os.Exit(0)
	ctx := context.TODO()
	statusList, err := getUnitStatus(ctx, servList)
	if err != nil {log.Fatalf("error -- getUnitStatus: %v", err)}
	log.Printf("get Units: %d\n", len(statusList))

	PrintUnitStatus(statusList)

//	fmt.Println(reflect.TypeOf(status))
//	fmt.Println(status)
}

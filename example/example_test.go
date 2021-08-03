package example

import (
	"context"
	"fmt"
	"github.com/go-magic/mid-server/dispatcher"
	"testing"

	"github.com/go-magic/mid-server/register"
	"github.com/go-magic/mid-server/schedule"
)

func TestNewServer(t *testing.T) {
	center := register.NewRegisterCenter()
	center.Register(1, NewHttpServer)
	sch := schedule.NewSchedule(center, dispatcher.NewDispatcher(1))
	subResult, err := sch.Execute(context.Background(), initTasks())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(subResult)
}

func TestEngine(t *testing.T) {
	t.Fatal(startServer())
}

// internal/worker/sync_test.go
package worker_test

import (
	"testing"

	"bolao-copa/internal/worker"
)

func TestMapStatus_Scheduled(t *testing.T) {
	if worker.MapStatus("TIMED") != "scheduled" {
		t.Fatal("esperava scheduled")
	}
}

func TestMapStatus_Live(t *testing.T) {
	if worker.MapStatus("IN_PLAY") != "live" {
		t.Fatal("esperava live")
	}
}

func TestMapStatus_Finished(t *testing.T) {
	if worker.MapStatus("FINISHED") != "finished" {
		t.Fatal("esperava finished")
	}
}

func TestMapStatus_Awarded(t *testing.T) {
	if worker.MapStatus("AWARDED") != "finished" {
		t.Fatal("esperava finished para awarded")
	}
}
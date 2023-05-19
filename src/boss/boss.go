package boss

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	RUN_PATH         = "/run/"
	BOSS_STATUS_PATH = "/status"
	SCALING_PATH     = "/scaling/worker_count"
	SHUTDOWN_PATH    = "/shutdown"
)

type Boss struct {
	workerPool *WorkerPool
}

func (b *Boss) BossStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receive request to %s\n", r.URL.Path)

	output := struct{
		State	map[string]int	`json:"state"`
		Tasks	map[string]int	`json:"tasks"`
	}{
		b.workerPool.StatusCluster(),
		b.workerPool.StatusTasks(),
	}
	
	if b, err := json.MarshalIndent(output, "", "\t"); err != nil {
		panic(err)
	} else {
		w.Write(b)
	}
}

func (b *Boss) Close(w http.ResponseWriter, r *http.Request) {
	b.workerPool.Close()
	os.Exit(0)
}

func (b *Boss) ScalingWorker(w http.ResponseWriter, r *http.Request) {
	// STEP 1: get int (worker count) from POST body, or return an error
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("POST a count to /scaling/worker_count\n"))
		if err != nil {
			log.Printf("(1) could not write web response: %s\n", err.Error())
		}
		return
	}

	contents, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("could not read body of web request\n"))
		if err != nil {
			log.Printf("(2) could not write web response: %s\n", err.Error())
		}
		return
	}

	worker_count, err := strconv.Atoi(string(contents))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("body of post to /scaling/worker_count should be an int\n"))
		if err != nil {
			log.Printf("(3) could not write web response: %s\n", err.Error())
		}
		return
	}

	if worker_count > Conf.Worker_Cap {
		worker_count = Conf.Worker_Cap
		log.Printf("capping workers at %d to avoid big bills during debugging\n", worker_count)
	}
	log.Printf("Receive request to %s, worker_count of %d requested\n", r.URL.Path, worker_count)

	// STEP 2: adjust target worker count
	b.workerPool.SetTarget(worker_count)

	//respond with status
	b.BossStatus(w, r)
}

func (b *Boss) RunLambda(w http.ResponseWriter, r *http.Request) {
	b.workerPool.RunLambda(w, r)
}

func BossMain() (err error) {
	fmt.Printf("WARNING!  Boss incomplete (only use this as part of development process).")

	boss := Boss{
		workerPool: NewWorkerPool(),
	}

	// things shared by all servers
	http.HandleFunc(BOSS_STATUS_PATH, boss.BossStatus)
	http.HandleFunc(SCALING_PATH, boss.ScalingWorker)
	http.HandleFunc(RUN_PATH, boss.RunLambda)
	http.HandleFunc(SHUTDOWN_PATH, boss.Close)

	port := fmt.Sprintf(":%s", Conf.Boss_port)
	fmt.Printf("Listen on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
	panic("ListenAndServe should never return")
}

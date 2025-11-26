package main

import (
	"log"
	"net"
	"sync"
	"time"
)

//element in the queue
type Job struct {
	conn net.Conn

}

// represent the thread in the pool
type Worker struct {
	id int
	jobChan chan Job
	wg *sync.WaitGroup
}

//represent the thread pool
type Pool struct {
	jobQueue chan Job
	workers []*Worker
	wg sync.WaitGroup
}

func NewWorker(id int , jobChan chan Job , wg *sync.WaitGroup) *Worker {
	return &Worker{
		id: id,
		jobChan: jobChan,
		wg: wg,
	}
}

func (w *Worker) Start() {
	go func ()  {
		defer w.wg.Done()
		for job := range w.jobChan {
			log.Printf("Worker %d is handling job from %s" , w.id , job.conn.RemoteAddr())
			handleConnection(job.conn)
		}
	}()
}

func NewPool(numOfWorker int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers: make([]*Worker, numOfWorker),
	}
}

//push Job to queue
func (p *Pool) AddJob(conn net.Conn) {
	p.jobQueue <- Job{conn: conn}
}

func (p *Pool) Start() {
	for i := range len(p.workers) {
		p.wg.Add(1)
		worker := NewWorker(i , p.jobQueue , &p.wg)
		p.workers[i] = worker
		worker.Start()	
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() 
	var buf []byte = make([]byte , 1000)
	_ , err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, world\r\n"))
}

func main() {
	listener, err := net.Listen("tcp" , ":3000")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	//1 pool with 2 threads
	pool := NewPool(2)
	pool.Start()

	for {
		conn , err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		pool.AddJob(conn)
	}
}
package main

import "fmt"

func main() {
	pool := NewWorkerPool(3, 10)

	// Read results in a separate goroutine
	go func() {
		for result := range pool.Results() {
			fmt.Println("Result:", result)
		}
	}()

	// Submit jobs
	for i := 1; i <= 5; i++ {
		job := Job{ID: i, JobType: multiply, Data: i * 10} // Job attributes are private
		if err := pool.Submit(job); err != nil {
			fmt.Println("Submit error:", err)
		}
	}

	pool.ShutDown()
	fmt.Println("Pool shut down cleanly")
}

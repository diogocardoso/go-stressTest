package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type TestResult struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

type Report struct {
	TotalTime     time.Duration
	TotalRequests int
	StatusCodes   map[int]int
	MinDuration   time.Duration
	MaxDuration   time.Duration
	AvgDuration   time.Duration
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")

	flag.Parse()

	if *url == "" || *requests == 0 {
		log.Fatal("URL e número de requests são obrigatórios")
	}

	if *concurrency > *requests {
		*concurrency = *requests
	}

	results := make(chan TestResult, *requests)
	start := time.Now()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, *concurrency)

	// Distribuir requests
	for i := 0; i < *requests; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Adquire um slot
			defer func() { <-semaphore }() // Libera o slot

			result := makeRequest(*url)
			results <- result
		}(i)
	}

	// Esperar todos os requests terminarem
	go func() {
		wg.Wait()
		close(results)
	}()

	// Processar resultados
	report := Report{
		StatusCodes: make(map[int]int),
		MinDuration: time.Duration(1<<63 - 1), // Valor máximo para Duration
		MaxDuration: 0,
	}

	var totalDuration time.Duration
	for result := range results {
		report.TotalRequests++
		if result.Error == nil {
			report.StatusCodes[result.StatusCode]++
			// Atualizar estatísticas de tempo
			totalDuration += result.Duration
			if result.Duration < report.MinDuration {
				report.MinDuration = result.Duration
			}
			if result.Duration > report.MaxDuration {
				report.MaxDuration = result.Duration
			}
		} else {
			report.StatusCodes[-1]++ // Para erros
		}
	}

	if report.TotalRequests > 0 {
		report.AvgDuration = totalDuration / time.Duration(report.TotalRequests)
	}
	report.TotalTime = time.Since(start)

	// Imprimir relatório
	printReport(report)
}

func makeRequest(url string) TestResult {
	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return TestResult{
			StatusCode: -1,
			Duration:   duration,
			Error:      err,
		}
	}
	defer resp.Body.Close()

	return TestResult{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}

func printReport(report Report) {
	fmt.Printf("\n=== Relatório de Teste de Carga ===\n")
	fmt.Printf("Tempo total: %v\n", report.TotalTime)
	fmt.Printf("Total de requests: %d\n", report.TotalRequests)
	fmt.Printf("Tempo mínimo: %v\n", report.MinDuration)
	fmt.Printf("Tempo máximo: %v\n", report.MaxDuration)
	fmt.Printf("Tempo médio: %v\n\n", report.AvgDuration)

	fmt.Println("Distribuição de Status Code:")
	for code, count := range report.StatusCodes {
		if code == -1 {
			fmt.Printf("Erros: %d (%.2f%%)\n", count, float64(count)/float64(report.TotalRequests)*100)
		} else {
			fmt.Printf("Status %d: %d (%.2f%%)\n", code, count, float64(count)/float64(report.TotalRequests)*100)
		}
	}
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// ---------------- ALGORITHMS ----------------
func coinChangeRecursive(coins []int, target int, index int) int {
	recOps++
	if target == 0 {
		return 1
	} else if target < 0 || index == len(coins) {
		return 0
	} else {
		return coinChangeRecursive(coins, target-coins[index], index) +
			coinChangeRecursive(coins, target, index+1)
	}
}

func coinChangeIterative(coins []int, target int) int {
	dp := make([]int, target+1)
	dp[0] = 1

	for i := 0; i < len(coins); i++ {
		for j := coins[i]; j <= target; j++ {
			iterOps++
			dp[j] += dp[j-coins[i]]
		}
	}

	return dp[target]
}

// ---------------- DATA STORAGE ----------------

type Point struct {
	Input      int
	RecTime    int64
	IterTime   int64
	RecOps     int64
	IterOps    int64
	RecResult  int
	IterResult int
}

var history []Point
var recOps, iterOps int

// ---------------- TEMPLATE ----------------

var page = template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Time Complexity Dashboard</title>
	<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>

	<style>
		* {
			box-sizing: border-box;
			font-family: Inter, Arial, sans-serif;
		}

		body {
			margin: 0;
			padding: 40px;
			background: linear-gradient(180deg, #5B5FB2, #4A4E9E);
			min-height: 100vh;
		}

		/* Search Bar */
		.search-bar {
			display: flex;
			align-items: center;
			gap: 10px;
			max-width: 500px;
			background: white;
			border-radius: 999px;
			padding: 10px 16px;
			box-shadow: 0 4px 12px rgba(0,0,0,0.1);

			/* CENTER HORIZONTALLY */
			margin: 0 auto 30px auto;
		}


		.search-bar input {
			border: none;
			outline: none;
			font-size: 16px;
			flex: 1;
		}

		.search-bar button {
			border: none;
			background: #ff4d4f;
			color: white;
			padding: 6px 14px;
			border-radius: 20px;
			cursor: pointer;
		}

		/* Card */
		.card {
			background: white;
			border-radius: 16px;
			padding: 24px;
			box-shadow: 0 10px 30px rgba(0,0,0,0.15);
		}

		h2, h3 {
			margin-top: 0;
		}

		/* Charts */
		.chart-container {
			display: flex;
			gap: 30px;
			margin-top: 30px;
		}

		.chart {
			width: 50%;
		}

		/* Table */
		table {
			width: 100%;
			border-collapse: collapse;
			margin-top: 30px;
		}

		th, td {
			padding: 14px;
			text-align: left;
			border-bottom: 1px solid #eee;
		}

		th {
			background: #f8f9fb;
			color: #555;
			font-size: 14px;
		}

		tr:hover {
			background: #fafafa;
		}


	</style>
</head>

<body>

<!-- SEARCH -->
<form method="POST" class="search-bar">
	<input type="number" name="amount" placeholder="Masukkan angka..." required>
	<button type="submit">Submit</button>
</form>

<!-- MAIN CARD -->
<div class="card">

	<h2>Coin Change Time Complexity</h2>

	<div class="chart-container">
		<div class="chart">
			<h3>Running Time</h3>
			<canvas id="recursiveChart"></canvas>
		</div>
		<div class="chart">
			<h3>Base Operations</h3>
			<canvas id="iterativeChart"></canvas>
		</div>
	</div>

	<h3>Results</h3>

	<table>
		<tr>
			<th>Input</th>
			<th>Recursive Result</th>
			<th>Iterative Result</th>
			<th>Recursive Ops</th>
			<th>Iterative Ops</th>
		</tr>

		{{range .History}}
		<tr>
			<td>{{.Input}}</td>
			<td>{{.RecResult}}</td>
			<td>{{.IterResult}}</td>
			<td>{{.RecOps}}</td>
			<td>{{.IterOps}}</td>
		</tr>
		{{end}}
	</table>

</div>

<script>
const labels = [{{range .History}}{{.Input}},{{end}}];

const recursiveData = [{{range .History}}{{.RecTime}},{{end}}];
const iterativeData = [{{range .History}}{{.IterTime}},{{end}}];
const recOpsData = [{{range .History}}{{.RecOps}},{{end}}];
const iterOpsData = [{{range .History}}{{.IterOps}},{{end}}];

new Chart(document.getElementById('recursiveChart'), {
	type: 'line',
	data: {
		labels: labels,
		datasets: [
			{
				label: 'Recursive Time (ms)',
				data: recursiveData,
				borderColor: 'rgba(255, 99, 132, 1)',
				backgroundColor: 'rgba(255, 99, 132, 0.2)',
				yAxisID: 'y'
			},
			{
				label: 'Iterative Time (ms)',
				data: iterativeData,
				borderColor: 'rgba(54, 162, 235, 1)',
				backgroundColor: 'rgba(54, 162, 235, 0.2)',
				yAxisID: 'y'
			}
		]
	},
	options: {
		scales: {
			y: {
				type: 'linear',
				position: 'left',
				title: { display: true, text: 'Time (ms)' }
			}
		}
	}
});

new Chart(document.getElementById('iterativeChart'), {
	type: 'line',
	data: {
		labels: labels,
		datasets: [
			{
				label: 'Recursive Operations',
				data: recOpsData,
				borderColor: 'rgba(255, 99, 132, 1)',
				backgroundColor: 'rgba(255, 99, 132, 0.2)',
				yAxisID: 'y'
			},
			{
				label: 'Iterative Operations',
				data: iterOpsData,
				borderColor: 'rgba(54, 162, 235, 1)',
				backgroundColor: 'rgba(54, 162, 235, 0.2)',
				yAxisID: 'y'
			}
		]
	},
	options: {
		scales: {
			y: {
				type: 'linear',
				position: 'left',
				title: { display: true, text: 'Operations' }
			}
		}
	}
});
</script>


</body>
</html>
`))

// ---------------- HANDLER ----------------

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		amount, _ := strconv.Atoi(r.FormValue("amount"))

		coins := []int{1, 2, 5, 10, 20}

		// Iterative timing
		iterOps = 0
		startIter := time.Now()
		iterResult := coinChangeIterative(coins, amount)
		iterTime := time.Since(startIter).Milliseconds()

		// recursive timing
		recOps = 0
		startRec := time.Now()
		recResult := coinChangeRecursive(coins, amount, 0)
		recTime := time.Since(startRec).Milliseconds()

		history = append(history, Point{
			Input:      amount,
			RecTime:    recTime,
			IterTime:   iterTime,
			RecOps:     int64(recOps),
			IterOps:    int64(iterOps),
			RecResult:  recResult,
			IterResult: iterResult,
		})
		sort.Slice(history, func(i, j int) bool {
			return history[i].Input < history[j].Input
		})
	}

	page.Execute(w, struct {
		History []Point
	}{
		History: history,
	})
}

// ---------------- MAIN ----------------

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

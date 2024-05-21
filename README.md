# MonteCarloPointDistribution
Monte Carlo simulations often use the distribution of points in a circle because it is simple and effective for demonstrating randomness and uniform distribution. This concept is analogous to the distribution of user devices within a circular cellular base station cell. By simulating random placements of devices, engineers can analyze network performance, optimize resource allocation, and improve coverage. This method helps in understanding signal strength, interference, and data throughput, ensuring reliable service for users within the cell.

The current repository contains a Monte Carlo simulation implementation designed to explore the distribution of points in two distinct scenarios:
 * **Uniform Points Position in the Circle**: Points are distributed uniformly throughout the entire area of the circle.
 * **Uniform Distance from the Center**: Points are distributed such that their distances from the center are uniform.

The simulation generates points for both scenarios and compares the results against their respective analytical models. This comparison helps in understanding the differences between the two distribution methods and validates the accuracy of the simulation by checking how closely the simulated results match the theoretical expectations.

## Requirements

* Golang >= 1.22
* Python >= 3.8

## Setup

After instaling the requeriments, the Jupyter Notebook must be installed to run the notebook with the models comparation.

```
pip install -r requeriments.txt
```

## Running

The Monte Carlo Simulation was writed in Golang, to run the simulator:
```
go run cmd/monte_carlo/main.go -r 6 -k 40 -seed 42 -t 0 -radius 0.75
```
The code must be execute two times to generate both experiments scenarios. After that, run Jupyter Notebook to visualize the comparation notebook:
```
jupyter notebook
```
Then, select the file CirclesPlot.ipynb in the notebooks directory to visualize the models comparation.

### Examples
* Running 10¹ - 10⁶ experiments with 30 points each uniform located points and radius 0.5
```
go run cmd/monte_carlo/main.go -r 6 -k 30 -seed 42 -t 0 -radius 0.5
```
* Running the previous experiments with uniform distance points
```
go run cmd/monte_carlo/main.go -r 6 -k 30 -seed 42 -t 1 -radius 0.5
```
* Adding a buffer contraint of the number of generated points in memory during the simulation of 10% of the expected number
```
go run cmd/monte_carlo/main.go -r 6 -k 30 -seed 42 -t 1 -radius 0.5 -points-buffer 0.1
```
* Adding a buffer contraint of the number of generated metrics in memory during the simulation of 10% of the expected number
```
go run cmd/monte_carlo/main.go -r 6 -k 30 -seed 42 -t 1 -radius 0.5 -results-buffer 0.1
```
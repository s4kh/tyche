## Approach
1. Clear out misunderstandings in the challenge. Questions about the task.
2. Single process prototype. Spawn single worker and send get request
3. Make it concurrent
4. Check whether it satisfies the requirements
5. Optimize

## Solution
Spawns 16 workers with different `workerID` and `port` to get numbers from. Consumers are `goroutines`, called right after spawning the worker with same `port` then send requests and processes response. The ratio is 1:1. 
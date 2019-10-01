## Approach
1. Clear out misunderstandings in the challenge. Questions about the task.
2. Single process prototype. Spawn single worker and send get request
3. Make it concurrent
4. Check whether it satisfies the requirements
5. Optimize
6. __Create worker pool size of 16. Whenever there is an idle worker in the pool we can start task for our 150 data sample.__

## Solution
Spawns 16 workers with different `workerID` and `port` to get numbers from. Consumers are `goroutines`, called right after spawning the worker with same `port` then send requests and processes response. The ratio is 1:1.
1. `main` function spawns a worker(`spawnWorker`) and a consumer(`callEndpoint` as a goroutine) total of 16 of them.
2. Workers does not have to be goroutines because `cmd.Start` function does not wait for the process to finish.
3. Consumer calls the endpoint then reads from response body until it reaches `EOF` or error happens(in our case it means worker has crashed)
4. Writes total numbers read by current consumer to a channel
5. At the end ranging through the channel calculate the sum of all read numbers

Note: `sync.WaitGroup` is used for waiting all goroutines to finish

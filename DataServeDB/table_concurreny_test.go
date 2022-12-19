package main

/*
	Description: Tests concurrency related bugs and issues for table operations. Mainly race conditions and deadlocks.

	Notes:
		- Use REST API calls, if they are causing problems in deadlock detection then direct table calls.
		- Use mix of read and write operations.
		- Use mix of single and multiple rows operations.
		- Use go's own concurrency testing. See: https://forum.golangbridge.org/t/how-do-you-unit-test-a-concurrent-data-structure/26912
		- If there is better 3rd party library for concurrency testing, then get it approved.

*/

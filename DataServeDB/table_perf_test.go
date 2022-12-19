package main

/*

	Description: Tests performance of table operations.

	Notes:
		- Tests should only test performance and not the correctness of the code.
		- Use REST API as it will be mainly used in production.
		- Go has built in benchmarking api, use that.
		- Table create, change, and delete operations will not be used frequently in production, so they are not tested here.
		- Rows operations are tested here.
			- Use mix of insert, update, and delete operations.
			- Use mix of single and multiple rows operations.
			- Use 30% write operations and 70% read operations.
			- Use 50% write operations and 50% read operations.
			- Use 70% write operations and 30% read operations.
			- Above is to check performance of different write/read ratios. Usually, there will be more read operations than write operations.
			- Do tests with different number of rows in the table. 10, 100, 1,000, 100,000, and 1,000,000 rows.
			- For write operations also update indexed and non-indexed columns and check performance.

*/

package assertly

import "testing"

func Test_AssertValues(t *testing.T) {
	var actual = `[
	{
		"id": 1,
		"name": "user 1",
		"perf_rank": 100,
		"perf_score": "6.50",
		"quiz": "{\n\t\"1\": {\n\t\t\"id\": 1,\n\t\t\"score\": 10,\n\t\t\"taken\": \"2018-01-10 16:02:01 UTC\"\n\t},\n\t\"2\": {\n\t\t\"id\": 2,\n\t\t\"score\": 3,\n\t\t\"taken\": \"2018-01-15 08:02:23 UTC\"\n\t}\n}",
		"visited": "2018-01-15 08:02:23Z"
	},
	{
	"id": 2,
	"name": "user 2",
	"perf_rank": 101,
	"perf_score": "7.00",
	"quiz": "{\n\t\"1\": {\n\t\t\"id\": 1,\n\t\t\"score\": 10,\n\t\t\"taken\": \"2018-01-11 13:01:48 UTC\"\n\t},\n\t\"2\": {\n\t\t\"id\": 2,\n\t\t\"score\": 4,\n\t\t\"taken\": \"2018-01-12 09:00:26 UTC\"\n\t}\n}",
	"visited": "2018-01-12 09:00:26Z"
	},
	{
	"id": 3,
	"name": "user 3",
	"perf_rank": 99,
	"perf_score": "5.00",
	"quiz": "{\n\t\"1\": {\n\t\t\"id\": 1,\n\t\t\"score\": 5,\n\t\t\"taken\": \"2018-01-10 05:01:33 UTC\"\n\t},\n\t\"2\": {\n\t\t\"id\": 2,\n\t\t\"score\": 5,\n\t\t\"taken\": \"2018-01-12 07:30:52 UTC\"\n\t}\n}",
	"visited": "2018-01-12 07:30:52Z"
	}
]`

	var expected = `[
	{
		"@indexBy@": [
			"id"
		],
		"@timeFormat@": "yyyy-MM-dd HH:mm:ss"
	},
	{
		"id": 1,
		"name": "user 1",
		"perf_rank": 100,
		"perf_score": 6.5,
		"quiz": {
			"1": {
				"id": 1,
				"score": 10,
				"taken": "2018-01-10 16:02:01 UTC"
			},
			"2": {
				"id": 2,
				"score": 3,
				"taken": "2018-01-15 08:02:23 UTC"
			}
		},
		"visited": "2018-01-15 08:02:23 UTC"
	},
	{
		"id": 2,
		"name": "user 2",
		"perf_rank": 101,
		"perf_score": 7,
		"quiz": {
			"1": {
				"id": 1,
				"score": 10,
				"taken": "2018-01-11 13:01:48 UTC"
			},
			"2": {
				"id": 2,
				"score": 4,
				"taken": "2018-01-12 09:00:26 UTC"
			}
		},
		"visited": "2018-01-12 09:00:26 UTC"
	},
	{
		"id": 3,
		"name": "user 3",
		"perf_rank": 99,
		"perf_score": 5,
		"quiz": {
			"1": {
				"id": 1,
				"score": 5,
				"taken": "2018-01-10 05:01:33 UTC"
			},
			"2": {
				"id": 2,
				"score": 5,
				"taken": "2018-01-12 07:30:52 UTC"
			}
		},
		"visited": "2018-01-12 07:30:52 UTC"
	}
]
`
	AssertValues(t, expected, actual)
}

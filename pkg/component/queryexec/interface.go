package queryexec

type QueryExecutor interface {
	QueryExec(bool) ([]any, string, int, error)
}

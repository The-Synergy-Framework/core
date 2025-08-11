# core/health

Context-aware health checks with standard statuses and a simple registry.

## Status values
- healthy
- degraded
- unhealthy
- unknown (probe failure)

## API
```go
// Define a checker
c := health.FuncChecker(func(ctx context.Context) (*health.Result, error) {
    // e.g., ping DB with timeout from ctx
    return health.OK("db ok", nil), nil
})

reg := health.New()
reg.Register("db", c)

sum := reg.RunAll(context.Background())
_ = sum.Overall
for _, e := range sum.Entries {
    _ = e.Name; _ = e.Result; _ = e.Duration; _ = e.Error
}
``` 
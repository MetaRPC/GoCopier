# Lifecycle recipes

Start a demo copier with explicit settings, then pause it. **Start is active immediately:** pausing afterward is not atomic and cannot prevent a trade in between. Use demo accounts with no existing positions while learning. This fragment belongs inside your application function; supply the account objects and handle failures explicitly.

```go
// master and slave are *copier.Account values for disposable demo accounts.
service, err := copier.NewCopierService("copy.mrpc.pro:443", ownerKey)
if err != nil { return err }
defer service.Close()
started, err := service.Start(ctx, &copier.StartRequest{
    Master: master, Slave: slave, RiskType: "FixedLot", RiskValue: "0.01",
    CopySl: true, CopyTp: true, CopyPendingOrders: false, ReverseCopy: false,
})
if err != nil { return err }
if !started.Ok { return fmt.Errorf("start refused: %s", started.Error) }
paused, err := service.Pause(ctx, started.CopierId, true)
if err != nil { return err }
if !paused.Ok { return fmt.Errorf("pause refused: %s", paused.Error) }
```

## Resume deliberately

Use the same Pause method with `false` (`paused: false` in Ruby). Check the reply and re-read List. A billing suspension can refuse a resume. Do not auto-resume every paused copier on application startup.

## Remove after inspection

Call the Remove method with the saved copier ID only after checking test positions and pending orders. Check the reply, then confirm the relationship no longer appears in List. Removal is not a substitute for closing positions.

For REST, use `/Remove` with `id`, `userKey`, and the `ApiKey` header. Inspect the JSON response as well as HTTP status. `/UserCopiers` lists the owner’s current relationships. The published [OpenAPI contract](https://copy.mrpc.pro/swagger/v1/swagger.json) describes the current REST routes and parameters.

## Monitor before trusting replication

Start/List/Pause/Remove tests validate management behavior. For a full acceptance test, open and close one demo position and compare both tickets, actual lots, and [trade logs](logs.md). See the [test checklist](testing.md).

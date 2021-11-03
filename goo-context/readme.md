# goo-context

## Cancel

```
func main() {
    for {
        select {
        case <-goo_context.Cancel().Done():
            goo_log.Debug("exit")
            return

        default:
            goo_log.Debug(time.Now().Format("15:04:05"))
            time.Sleep(time.Second)
        }
    }
}
```

## 

```
func main() {
    ctx := goo_context.Timeout(5 * time.Second)
    	
    for {
        select {
        case <-ctx.Done():
            goo_log.Debug("exit")
            return

        default:
            goo_log.Debug(time.Now().Format("15:04:05"))
            time.Sleep(time.Second)
        }
    }
}
```
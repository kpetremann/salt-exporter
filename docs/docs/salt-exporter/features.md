# Features

## Supported events

| Event                                   | Salt tag                 |
|-----------------------------------------|--------------------------|
| Metrics for new job                     | `salt/job/<jid>/new`     |
| Metrics for response from a minion      | `salt/job/<jid>/ret/<*>` |
| Metrics for new runner job              | `salt/run/<jid>/new`     |
| Metrics for runner response             | `salt/run/<jid>/ret/<*>` |

## Customization

| Feature | Details |
|---------|---------|
| `<metric>`.enabled | All metrics can be either enabled or disabled |
| Add minion label for some metrics | It is not recommended on large environment as it could lead to cardinality issues |
| Filter out `test=true`/`mock=true` events| It can be useful to ignore tests |
| `salt_function_status`: dedicated metric for function/state per minion | The filter should be used to avoid cardinality issues |

See the [configuration page](./configuration.md) to use these features.

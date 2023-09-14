# IBCI wrapped ibc-go

IBCI is an integration interface for IBC. This is an example of what an IBCI wrapped version of ibc-go could look like.

There's a few high level goals here:

- Provide a minimal interface for integrations
- Standardize on IBC internal data structures
- Reuse existing relayer infrastructure
- Avoid costly re-factor of ibc-go

The approach is fairly simple:

- Wrap the IBC core module in an abci.Application
- Create a dedicated IBC mempool using cometBFT's mempool
  - The mempool impl is abci.Application compatible, can migrate to a diff mempool impl in the future
- Expose RPC endpoints in the format that existing relayers already expect (e.g. hermes)
- Impl the [IBCI](https://www.notion.so/polymer-labs/ADR-13-Enshrined-IBC-b2f4cb4bee254025ba3bdd15e650543d?pvs=4#ca1b8bab96944b63a3ec18383d71c634) for integrations

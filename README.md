# Blockchain IT Asset Management

Blockchain-based IT asset tracking system using Hyperledger Fabric for immutable audit trails of asset lifecycle events including procurement, assignment, transfers, and decommissioning.

Personal project, built to explore immutable asset-lifecycle audit trails on Hyperledger Fabric. It is not production software — see **Status** below for exactly what is and isn't implemented.

## Status

**Implemented**

- Go chaincode for asset create/transfer/query with history
- Express API exposing the chaincode
- Fabric connection profile

**Not implemented / known limitations**

- No React front-end (the earlier README claimed `App.jsx` and `AssetList.jsx`)
- Requires a running Fabric network; no test network included
- No chaincode unit tests

## Built with

- **Node** — express, fabric-client, fabric-ca-client

## Running it

```bash
npm install
npm start
```

## Layout

```
chaincode/
  asset_contract.go
config/
  connection-profile.json
package.json
src/
  routes/
    assets.js
  server.js
```


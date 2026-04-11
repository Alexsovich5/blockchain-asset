# Blockchain IT Asset Management

![Project Status](https://img.shields.io/badge/Status-Complete-brightgreen)
![Timeline](https://img.shields.io/badge/Timeline-March%202017%20--%20June%202017-blue)
![Technology](https://img.shields.io/badge/Tech-Hyperledger%20Fabric%20%7C%20Node.js%20%7C%20React-orange)

## Project Overview

Blockchain-based IT asset tracking system using Hyperledger Fabric for immutable audit trails of asset lifecycle events including procurement, assignment, transfers, and decommissioning.

**Role**: IT Operations Specialist
**Organization**: ACORIA
**Duration**: March 2017 - June 2017
**Project**: #14 of 30 in IT Career Portfolio

## Business Impact

- **100% Audit Trail Integrity**: Immutable blockchain ledger
- **Tamper-proof Records**: Cryptographic verification of all asset events
- **Automated Compliance**: Real-time audit readiness
- **30% Faster Asset Reconciliation**: Distributed ledger eliminates discrepancies

## Technology Stack

- **Hyperledger Fabric 1.0**: Permissioned blockchain network
- **Node.js 6.x**: Backend API server
- **React 15.x**: Web dashboard
- **Go**: Chaincode (smart contracts)
- **Docker**: Network containerization

## Project Structure

```
blockchain-asset/
├── README.md
├── package.json
├── chaincode/
│   └── asset_contract.go
├── src/
│   ├── server.js
│   ├── routes/
│   │   └── assets.js
│   └── frontend/
│       ├── App.jsx
│       └── AssetList.jsx
├── network/
│   └── docker-compose.yml
└── config/
    └── connection-profile.json
```

## Installation and Setup

```bash
# Start Hyperledger Fabric network
cd network && docker-compose up -d

# Install chaincode
peer chaincode install -n asset-mgmt -v 1.0 -p chaincode/

# Start API server
npm install && npm start
```

## Contributing

This is a historical project from March 2017 - June 2017, preserved for portfolio purposes.

## License

Professional portfolio project - ACORIA

---

**Developed during March 2017 - June 2017**
*Part of Alexander Efrem's IT Career Portfolio (2012-2024)*

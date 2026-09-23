/**
 * Blockchain Asset Management API Server
 */

'use strict';

const express = require('express');
const FabricClient = require('fabric-client');
const path = require('path');
const assetRoutes = require('./routes/assets');

const app = express();
app.use(express.json());
app.use(express.static('public'));

// Fabric client setup
const client = new FabricClient();
const channelName = 'asset-channel';
const chaincodeName = 'asset-mgmt';

async function initFabric() {
    const configPath = path.resolve(__dirname, '../config/connection-profile.json');
    await client.loadFromConfig(configPath);
    await client.initCredentialStores();
    console.log('Fabric client initialized');
}

// Mount routes
app.use('/api/assets', assetRoutes(client, channelName, chaincodeName));

app.get('/health', (req, res) => {
    res.json({ status: 'healthy', network: 'connected' });
});

const PORT = process.env.PORT || 3000;

initFabric()
    .then(() => {
        app.listen(PORT, () => {
            console.log(`Asset Management API running on port ${PORT}`);
        });
    })
    .catch(err => {
        console.error('Failed to initialize:', err);
        process.exit(1);
    });

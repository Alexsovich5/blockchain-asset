/**
 * Asset API Routes
 * IT Operations Specialist - ACORIA (2017)
 */

'use strict';

const express = require('express');

module.exports = function(fabricClient, channelName, chaincodeName) {
    const router = express.Router();

    async function invokeChaincode(fcn, args) {
        const channel = fabricClient.getChannel(channelName);
        const txId = fabricClient.newTransactionID();
        const request = {
            chaincodeId: chaincodeName,
            fcn: fcn,
            args: args,
            txId: txId
        };
        const results = await channel.sendTransactionProposal(request);
        await channel.sendTransaction({ proposalResponses: results[0], proposal: results[1] });
        return results[0][0].response.payload.toString();
    }

    async function queryChaincode(fcn, args) {
        const channel = fabricClient.getChannel(channelName);
        const request = {
            chaincodeId: chaincodeName,
            fcn: fcn,
            args: args
        };
        const result = await channel.queryByChaincode(request);
        return result[0].toString();
    }

    router.post('/register', async (req, res) => {
        try {
            const result = await invokeChaincode('registerAsset', [JSON.stringify(req.body)]);
            res.json(JSON.parse(result));
        } catch (err) {
            res.status(500).json({ error: err.message });
        }
    });

    router.get('/:id', async (req, res) => {
        try {
            const result = await queryChaincode('queryAsset', [req.params.id]);
            res.json(JSON.parse(result));
        } catch (err) {
            res.status(404).json({ error: 'Asset not found' });
        }
    });

    router.post('/:id/transfer', async (req, res) => {
        try {
            const result = await invokeChaincode('transferAsset', [
                req.params.id, req.body.newOwner, req.body.newDepartment
            ]);
            res.json(JSON.parse(result));
        } catch (err) {
            res.status(500).json({ error: err.message });
        }
    });

    router.get('/:id/history', async (req, res) => {
        try {
            const result = await queryChaincode('getHistory', [req.params.id]);
            res.json(JSON.parse(result));
        } catch (err) {
            res.status(500).json({ error: err.message });
        }
    });

    return router;
};

'use strict';

const util = require('./util.js');
const e2eUtils = require('./e2eUtils.js');
const impl_create = require('./create-channel.js');
const impl_join = require('./join-channel.js');
const impl_install = require('./install-chaincode.js');
const impl_instantiate = require('./instantiate-chaincode.js');
const BlockchainInterface = require('../../comm/blockchain-interface.js');
const commUtils = require('../../comm/util');
const commLogger = commUtils.getLogger('accelerator.js');
const TxStatus = require('../../comm/transaction');

class Accelerator extends BlockchainInterface {
    constructor(config_path) {
        super(config_path);
    }

    async init() {
        util.init(this.configPath);
        e2eUtils.init(this.configPath);
        try {
            await impl_create.run(this.configPath);
            await impl_join.run(this.configPath);
        } catch (err) {
            commLogger.error(`Fabric initialization failed: ${(err.stack ? err.stack : err)}`);
            throw err;
        }
    }

    async installSmartContract() {
        // todo: now all chaincodes are installed and instantiated in all peers, should extend this later
        try {
            await impl_install.run(this.configPath);
            await impl_instantiate.run(this.configPath);
        } catch (err) {
            commLogger.error(`Fabric chaincode install failed: ${(err.stack ? err.stack : err)}`);
            throw err;
        }
    }

    async getContext(name, args, clientIdx, txFile) {
        util.init(this.configPath);
        e2eUtils.init(this.configPath);
        this.txFile = txFile;
        if(this.txFile){
            this.txFile.name = name;
            commLogger.debug('getContext) name: ' + name +  ' clientIndex: ' + clientIdx + ' txFile: ' + JSON.stringify(this.txFile));
            if(this.txFile.readWrite === 'read') {
                if(this.txFile.roundCurrent === 0){
                    await e2eUtils.readFromFile(this.txFile.name);
                }
            }
        }
        let config  = require(this.configPath);
        let context = config.fabric.context;

        let channel;
        if(typeof context === 'undefined') {
            channel = util.getDefaultChannel();
        }
        else{
            channel = util.getChannel(context[name]);
        }

        if(!channel) {
            throw new Error('Could not find context information in the config file');
        }

        return await e2eUtils.getcontext(channel, clientIdx, txFile);
    }

    async releaseContext(context) {
        if(this.txFile && this.txFile.readWrite === 'write') {
            if(this.txFile.roundCurrent === (this.txFile.roundLength - 1)){
                await e2eUtils.writeToFile(this.txFile.name);
            }
        }
        await e2eUtils.releasecontext(context);
        await commUtils.sleep(1000);
    }

    async invokeSmartContract(context, contractID, contractVer, args, timeout) {
        let promises = [];
        args.forEach((item, index)=>{
            try {
                let simpleArgs = [];
                let func;
                for(let key in item) {
                    if(key === 'transaction_type') {
                        func = item[key].toString();
                    }
                    else {
                        simpleArgs.push(item[key].toString());
                    }
                }
                if(func) {
                    simpleArgs.splice(0, 0, func);
                }
                promises.push(e2eUtils.invokebycontext(context, contractID, contractVer, simpleArgs, timeout));
            }
            catch(err) {
                commLogger.error(err);
                let badResult = new TxStatus('artifact');
                badResult.SetStatusFail();
                promises.push(Promise.resolve(badResult));
            }
        });

        return await Promise.all(promises);
    }

    async queryState(context, contractID, contractVer, key, fcn = 'query') {
        // TODO: change string key to general object
        return await e2eUtils.querybycontext(context, contractID, contractVer, key.toString(), fcn);
    }
}



module.exports = Accelerator;

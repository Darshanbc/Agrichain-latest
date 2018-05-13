/**
 * Copyright 2017 IBM All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */
var path = require('path');
var fs = require('fs');
var util = require('util');
var hfc = require('fabric-client');
var helper = require('./helper.js');
var logger = helper.getLogger('Query');
var invoke = require('./invoke-transaction.js')
var py =require('python-shell')

var queryChaincode = async function(peer, channelName, chaincodeName, args, fcn, username, org_name) {
	try {
		// first setup the client for this org
		var client = await helper.getClientForOrg(org_name, username);
		logger.debug('Successfully got the fabric client for the organization "%s"', org_name);
		var channel = client.getChannel(channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', channelName);
			logger.error(message);
			throw new Error(message);
		}

		// send query
		var request = {
			targets : [peer], //queryByChaincode allows for multiple targets
			chaincodeId: chaincodeName,
			fcn: fcn,
			args: args
		};
		let response_payloads = await channel.queryByChaincode(request);
		if (response_payloads) {
			for (let i = 0; i < response_payloads.length; i++) {
				logger.info(args[0]+' now has ' + response_payloads[i].toString('utf8') +
					' after the move');
				return response_payloads[i].toString('utf8') ;
			}
		} else {
			logger.error('response_payloads is null');
			return 'response_payloads is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};
var getBlockByNumber = async function(peer, channelName, blockNumber, username, org_name) {
	try {
		// first setup the client for this org
		var client = await helper.getClientForOrg(org_name, username);
		logger.debug('Successfully got the fabric client for the organization "%s"', org_name);
		var channel = client.getChannel(channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', channelName);
			logger.error(message);
			throw new Error(message);
		}

		let response_payload = await channel.queryBlock(parseInt(blockNumber, peer));
		if (response_payload) {
			logger.debug(response_payload);
			return response_payload;
		} else {
			logger.error('response_payload is null');
			return 'response_payload is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};
var getTransactionByID = async function(peer, channelName, trxnID, username, org_name) {
	try {
		// first setup the client for this org
		var client = await helper.getClientForOrg(org_name, username);
		logger.debug('Successfully got the fabric client for the organization "%s"', org_name);
		var channel = client.getChannel(channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', channelName);
			logger.error(message);
			throw new Error(message);
		}

		let response_payload = await channel.queryTransaction(trxnID, peer);
		if (response_payload) {
			logger.debug(response_payload);
			return response_payload;
		} else {
			logger.error('response_payload is null');
			return 'response_payload is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};

var getWriteSetTransactionByID = async function(peer, channelName, trxnID, username, org_name) {
	try {
		// first setup the client for this org
		var client = await helper.getClientForOrg(org_name, username);
		logger.debug('Successfully got the fabric client for the organization "%s"', org_name);
		var channel = client.getChannel(channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', channelName);
			logger.error(message);
			throw new Error(message);
		}

		let response_payload = await channel.queryTransaction(trxnID, peer);
		if (response_payload) {
			logger.debug(response_payload);
			var writeArr=response_payload["transactionEnvelope"]["payload"]["data"]["actions"][0]["payload"]["action"]["proposal_response_payload"]["extension"]["results"]["ns_rwset"][0]["rwset"]["writes"]		
			// logger.debug(res)
			var fertdetails
			writeArr.forEach(function(element) {
				if(element.key==trxnID){
					console.log(element.value);
					fertdetails=element.value
				}
			})
					// });
			return fertdetails;
		} else {
			logger.error('response_payload is null');
			return 'response_payload is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};
var getBlockByHash = async function(peer, channelName, hash, username, org_name) {
	try {
		// first setup the client for this org
		var client = await helper.getClientForOrg(org_name, username);
		logger.debug('Successfully got the fabric client for the organization "%s"', org_name);
		var channel = client.getChannel(channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', channelName);
			logger.error(message);
			throw new Error(message);
		}

		let response_payload = await channel.queryBlockByHash(Buffer.from(hash), peer);
		if (response_payload) {
			logger.debug(response_payload);
			return response_payload;
		} else {
			logger.error('response_payload is null');
			return 'response_payload is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};
var getChainInfo = async function(peer, channelName, username, org_name) {
	try {
		// first setup the client for this org
		var client = await helper.getClientForOrg(org_name, username);
		logger.debug('Successfully got the fabric client for the organization "%s"', org_name);
		var channel = client.getChannel(channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', channelName);
			logger.error(message);
			throw new Error(message);
		}

		let response_payload = await channel.queryInfo(peer);
		if (response_payload) {
			logger.debug(response_payload);
			return response_payload;
		} else {
			logger.error('response_payload is null');
			return 'response_payload is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};
//getInstalledChaincodes
var getInstalledChaincodes = async function(peer, channelName, type, username, org_name) {
	try {
		// first setup the client for this org
		var client = await helper.getClientForOrg(org_name, username);
		logger.debug('Successfully got the fabric client for the organization "%s"', org_name);

		let response = null
		if (type === 'installed') {
			response = await client.queryInstalledChaincodes(peer, true); //use the admin identity
		} else {
			var channel = client.getChannel(channelName);
			if(!channel) {
				let message = util.format('Channel %s was not defined in the connection profile', channelName);
				logger.error(message);
				throw new Error(message);
			}
			response = await channel.queryInstantiatedChaincodes(peer, true); //use the admin identity
		}
		if (response) {
			if (type === 'installed') {
				logger.debug('<<< Installed Chaincodes >>>');
			} else {
				logger.debug('<<< Instantiated Chaincodes >>>');
			}
			var details = [];
			for (let i = 0; i < response.chaincodes.length; i++) {
				logger.debug('name: ' + response.chaincodes[i].name + ', version: ' +
					response.chaincodes[i].version + ', path: ' + response.chaincodes[i].path
				);
				details.push('name: ' + response.chaincodes[i].name + ', version: ' +
					response.chaincodes[i].version + ', path: ' + response.chaincodes[i].path
				);
			}
			return details;
		} else {
			logger.error('response is null');
			return 'response is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};
var getChannels = async function(peer, username, org_name) {
	try {
		// first setup the client for this org
		var client = await helper.getClientForOrg(org_name, username);
		logger.debug('Successfully got the fabric client for the organization "%s"', org_name);

		let response = await client.queryChannels(peer);
		if (response) {
			logger.debug('<<< channels >>>');
			var channelNames = [];
			for (let i = 0; i < response.channels.length; i++) {
				channelNames.push('channel id: ' + response.channels[i].channel_id);
			}
			logger.debug(channelNames);
			return response;
		} else {
			logger.error('response_payloads is null');
			return 'response_payloads is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};

var queryAndInvoke = async function(peer, trxnId, username, org_name,channelName, chaincodeName,  permission,fcn) {
	logger.debug("peer "+peer+"\ntxnId:"+trxnId)
	logger.debug("username "+username)
	logger.debug("org "+org_name)
	logger.debug("channelName"+channelName)
	logger.debug("chaincode name "+chaincodeName)
	logger.debug("fcn"+ fcn)
	// logger.debug( "targets "+peerTargets)
	logger.debug(" permission"+ permission)
	// =====================================

	try {
			// first setup the client for this org
			var client = await helper.getClientForOrg(org_name, username);
			logger.debug('Successfully got the fabric client for the organization "%s"', org_name);
			var channel = client.getChannel(channelName);
			if(!channel) {
				let message = util.format('Channel %s was not defined in the connection profile', channelName);
				logger.error(message);
				throw new Error(message);
			}
			
			let response_payloads = await channel.queryTransaction(trxnId,peer);
			logger.info("response payload:"+ response_payloads)
			if (response_payloads) {
				for (let i = 0; i < response_payloads.length; i++) {
					logger.info("reponse paloads"+response_payloads[i].toString('utf8'));

					logger.debug("validation code:"+response_payloads["validationCode"])
				}


				if(fcn==="ApproveOrDenyFertilizer"){
					var Id0=response_payloads["transactionEnvelope"]["payload"]["data"]["actions"][0]["payload"]["action"]["proposal_response_payload"]["extension"]["results"]["ns_rwset"][0]["rwset"]["reads"][0]["key"];
					var Id1=response_payloads["transactionEnvelope"]["payload"]["data"]["actions"][0]["payload"]["action"]["proposal_response_payload"]["extension"]["results"]["ns_rwset"][0]["rwset"]["reads"][1]["key"]			
					var Id2=response_payloads["transactionEnvelope"]["payload"]["data"]["actions"][0]["payload"]["action"]["proposal_response_payload"]["extension"]["results"]["ns_rwset"][0]["rwset"]["reads"][2]["key"]	
					var enrollId,cropId
					var args=[]
					var id=[]
					if (Id0.search("Criyagen")!=0){
						id.push(Id0)
					}
					if (Id1.search("Criyagen")!=0){
						id.push(Id1)
					}
					if(Id2.search("Criyagen")!=0){
						id.push(Id2)
					}
					var temp
					if (id[0].search("CROP")==0){
						cropId=id[0]
						enrollId=id[1]
					}else{
						cropId=id[1]
						enrollId=id[0]
					}
					logger.debug("enrollId: "+Id0)
					logger.debug("cropId: "+Id1)
					logger.debug("Criyagen: "+Id2)
					args.push(cropId)
					args.push(enrollId)
					args.push(trxnId)

					args.push(permission)
					logger.debug(args)
					
				}
			// 	var peerObj=JSON.parse(peerTargets)
			// 	logger.debug(peerObj)
			// //var msg
				var message= invoke.invokeChaincode(peer,channelName, chaincodeName, fcn, args, username, org_name)
				return message

			}else {
				logger.error('response_payloads is null');
				return 'response_payloads is null';
			}
		
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};
	// ==========================

exports.queryAndInvoke=queryAndInvoke;
exports.queryChaincode = queryChaincode;
exports.getBlockByNumber = getBlockByNumber;
exports.getTransactionByID = getTransactionByID;
exports.getBlockByHash = getBlockByHash;
exports.getChainInfo = getChainInfo;
exports.getInstalledChaincodes = getInstalledChaincodes;
exports.getChannels = getChannels;
exports.getWriteSetTransactionByID=getWriteSetTransactionByID
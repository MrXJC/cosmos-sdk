# Overview

## What is a Light Client?

A light client has all the security of a full node with minimal bandwidth requirements. The minimal bandwidth requirements allows developers to build fully secure, efficient and usable mobile apps, websites or any other application that does not want to download and follow the full blockchain.

A full node of ABCI is different from its light client in the following ways:

- Full Node
  - Node discovery 
  - Verify and broadcast valid transactions in mempool
  - Verify and store new blocks
  - If this node is a validator node, it could contribute to protect the safety of network and reach consensus.
  - Resource consuming: huge computing resources for transaction verification and huge storage resources for saving blocks
- Light Client 
  - Redirect requests to full nodes
  - Verify transaction according to its hash
  - Verify precommit info at specific height
  - Verify block at specific height
  - Verify the proof in abci query result
  - Only need limited computing and storage resources, available for mobiles and personal computers

## High-Level Architecture

An application developer that would like to build a third party integration can ship his application with the LCD for the Cosmos Hub (or any other zone) and only needs to initialise it. Afterwards his application can interact with the zone as if it was running against a full node.

![high-level](/Users/suyu/Documents/bianjie/exchange/high-level.png)

# Design Details

## Trusted validator set

The base design philosophy of lcd follows the two rules: 

1. **Doesn't trust any blockchin nodes, including validator nodes and other full nodes**
2. **Only trusts the whole validator set**

The original trusted validator set should be prepositioned into its trust store, usually this validator set comes from genesis file. During running time, if LCD detects different validator set, it will verify it and save new validated validator set to trust store.

![validator-set-change](/Users/suyu/Library/Containers/com.tencent.xinWeChat/Data/Library/Application%20Support/com.tencent.xinWeChat/2.0b4.0.9/52a79f1e0cc6e5dfa6e763bee3131e18/Message/MessageTemp/ca41b11bed09054ddbe2ab485a3b815e/File/trust-scope/picture/validatorSetChange.png)

If a new validatorset with more than 1/3 voting power different from current validatorset, then there must be one more other validatorset changes during this period. LCD will find out all changes and ensure each change doesn't affect more than 1/3 voting power. Detailed description about this process will be posted later.

![change-process](/Users/suyu/Library/Containers/com.tencent.xinWeChat/Data/Library/Application%20Support/com.tencent.xinWeChat/2.0b4.0.9/52a79f1e0cc6e5dfa6e763bee3131e18/Message/MessageTemp/ca41b11bed09054ddbe2ab485a3b815e/File/trust-scope/picture/changeProcess.png)

## Trust propagation

From the above section, we come to know how to get trusted validator set and how lcd keeps track of validator set evolution. Validator set is the foundation of trust, and the trust can propagate to other blockchain data, such as block and transaction. The propagate architecture is shown as follows: 

![change-process](/Users/suyu/Library/Containers/com.tencent.xinWeChat/Data/Library/Application%20Support/com.tencent.xinWeChat/2.0b4.0.9/52a79f1e0cc6e5dfa6e763bee3131e18/Message/MessageTemp/ca41b11bed09054ddbe2ab485a3b815e/File/trust-scope/picture/trustPropagate.png)

Taking transaction trust propagation for example:

1. **Query a transaction by its hash and require the fullnode to return related proof**

2. **Then we can get the transaction data like this**

   ```
   // Result of querying for a tx
   type ResultTx struct {
     Hash     cmn.HexBytes           `json:"hash"`
     Height   int64                  `json:"height"`
     Index    uint32                 `json:"index"`
     TxResult abci.ResponseDeliverTx `json:"tx_result"`
     Tx       types.Tx               `json:"tx"`
     Proof    types.TxProof          `json:"proof,omitempty"`
   }
   ```

3. **According to the height in transaction data, we can get the commit at this height. The commit contains all pre-commit data for the block during reaching consensus process**

4. **With the trusted validator set at this height, check if the pre-commit have more than 2/3 voting power**

5. **Extract DataHash from the validated commit. The DataHash is the merkle root of all transactions in this block**

6. **Rebuild a merkle tree with the transaction proof and other transaction data, and verify if the merkle root matches the DataHash**

7. **If all above steps passed, then definitely the transaction data is trusted**

Block verification and ABCI state verification are similar.
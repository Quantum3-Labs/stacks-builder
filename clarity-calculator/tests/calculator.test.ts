// calculator.test.ts
// Comprehensive unit tests for the calculator smart contract
// Using Clarinet testing framework

import { Clarinet, Tx, Chain, Account, types } from "https://deno.land/x/clarinet@v1.7.1/index.ts";

Clarinet.test({
  name: "Calculator contract: Basic arithmetic operations",
  async fn(chain: Chain, accounts: Map<string, Account>) {
    const deployer = accounts.get("deployer")!;
    const user1 = accounts.get("wallet_1")!;

    // Test addition
    console.log("Testing 'add' function...");
    let block = chain.mineBlock([
      Tx.contractCall("calculator", "add", [types.uint(5), types.uint(3)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(8);
    console.log("  'add(5, 3)' passed. Result: 8");

    // Check get-last-result after add
    let lastResult = chain.callReadOnlyFn("calculator", "get-last-result", [], deployer.address);
    lastResult.result.expectOk().expectInt(8);
    console.log("  'get-last-result' after 'add' passed. Result: 8");

    // Test subtraction
    console.log("\nTesting 'subtract' function...");
    block = chain.mineBlock([
      Tx.contractCall("calculator", "subtract", [types.uint(10), types.uint(4)], user1.address)
    ]);
    block.receipts[0].result.expectOk().expectInt(6);
    console.log("  'subtract(10, 4)' passed. Result: 6");

    // Check get-last-result after subtract
    lastResult = chain.callReadOnlyFn("calculator", "get-last-result", [], user1.address);
    lastResult.result.expectOk().expectInt(6);
    console.log("  'get-last-result' after 'subtract' passed. Result: 6");

    // Test multiplication
    console.log("\nTesting 'multiply' function...");
    block = chain.mineBlock([
      Tx.contractCall("calculator", "multiply", [types.uint(7), types.uint(2)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(14);
    console.log("  'multiply(7, 2)' passed. Result: 14");

    // Check get-last-result after multiply
    lastResult = chain.callReadOnlyFn("calculator", "get-last-result", [], deployer.address);
    lastResult.result.expectOk().expectInt(14);
    console.log("  'get-last-result' after 'multiply' passed. Result: 14");

    // Test division
    console.log("\nTesting 'divide' function...");
    block = chain.mineBlock([
      Tx.contractCall("calculator", "divide", [types.uint(15), types.uint(3)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(5);
    console.log("  'divide(15, 3)' passed. Result: 5");

    // Check get-last-result after divide
    lastResult = chain.callReadOnlyFn("calculator", "get-last-result", [], deployer.address);
    lastResult.result.expectOk().expectInt(5);
    console.log("  'get-last-result' after 'divide' passed. Result: 5");
  },
});

Clarinet.test({
  name: "Calculator contract: Division by zero error handling",
  async fn(chain: Chain, accounts: Map<string, Account>) {
    const deployer = accounts.get("deployer")!;

    console.log("\nTesting 'divide' function with division by zero...");

    // First, set a known value for last-result
    let block = chain.mineBlock([
      Tx.contractCall("calculator", "add", [types.uint(100), types.uint(0)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(100);

    let initialLastResult = chain.callReadOnlyFn("calculator", "get-last-result", [], deployer.address);
    initialLastResult.result.expectOk().expectInt(100);
    console.log("  Initial 'last-result' set to 100.");

    // Attempt division by zero
    block = chain.mineBlock([
      Tx.contractCall("calculator", "divide", [types.uint(10), types.uint(0)], deployer.address)
    ]);
    block.receipts[0].result.expectErr().expectUint(100); // ERR-DIVISION-BY-ZERO
    console.log("  'divide(10, 0)' returned expected error u100.");

    // Verify that get-last-result was NOT updated after the failed operation
    let lastResultAfterError = chain.callReadOnlyFn("calculator", "get-last-result", [], deployer.address);
    lastResultAfterError.result.expectOk().expectInt(100); // Should still be 100
    console.log("  'get-last-result' did not update after error, still 100.");
  },
});

Clarinet.test({
  name: "Calculator contract: get-last-result initial state and updates",
  async fn(chain: Chain, accounts: Map<string, Account>) {
    const deployer = accounts.get("deployer")!;
    console.log("\nTesting 'get-last-result' initial state and updates...");

    // Test initial state of get-last-result
    let initialLastResult = chain.callReadOnlyFn("calculator", "get-last-result", [], deployer.address);
    initialLastResult.result.expectOk().expectInt(0);
    console.log("  Initial 'get-last-result' is 0 as expected.");

    // Perform an operation and verify get-last-result updates
    let block = chain.mineBlock([
      Tx.contractCall("calculator", "add", [types.uint(25), types.uint(17)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(42);

    let updatedLastResult = chain.callReadOnlyFn("calculator", "get-last-result", [], deployer.address);
    updatedLastResult.result.expectOk().expectInt(42);
    console.log("  'get-last-result' updated to 42 after operation.");

    // Perform another operation and verify it updates again
    block = chain.mineBlock([
      Tx.contractCall("calculator", "multiply", [types.uint(6), types.uint(7)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(42);

    let finalLastResult = chain.callReadOnlyFn("calculator", "get-last-result", [], deployer.address);
    finalLastResult.result.expectOk().expectInt(42);
    console.log("  'get-last-result' updated to 42 after second operation.");
  },
});

Clarinet.test({
  name: "Calculator contract: Edge cases and boundary values",
  async fn(chain: Chain, accounts: Map<string, Account>) {
    const deployer = accounts.get("deployer")!;
    console.log("\nTesting edge cases and boundary values...");

    // Test with zero values
    let block = chain.mineBlock([
      Tx.contractCall("calculator", "add", [types.uint(0), types.uint(0)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(0);
    console.log("  'add(0, 0)' passed. Result: 0");

    // Test subtraction resulting in zero
    block = chain.mineBlock([
      Tx.contractCall("calculator", "subtract", [types.uint(5), types.uint(5)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectInt(0);
    console.log("  'subtract(5, 5)' passed. Result: 0");

    // Test multiplication by zero
    block = chain.mineBlock([
      Tx.contractCall("calculator", "multiply", [types.uint(100), types.uint(0)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(0);
    console.log("  'multiply(100, 0)' passed. Result: 0");

    // Test division with remainder (integer division)
    block = chain.mineBlock([
      Tx.contractCall("calculator", "divide", [types.uint(10), types.uint(3)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(3); // 10/3 = 3 (integer division)
    console.log("  'divide(10, 3)' passed. Result: 3 (integer division)");

    // Test large numbers
    block = chain.mineBlock([
      Tx.contractCall("calculator", "add", [types.uint(1000000), types.uint(2000000)], deployer.address)
    ]);
    block.receipts[0].result.expectOk().expectUint(3000000);
    console.log("  'add(1000000, 2000000)' passed. Result: 3000000");
  },
});

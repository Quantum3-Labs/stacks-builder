# Clarity Calculator Smart Contract

A simple calculator smart contract written in Clarity that demonstrates basic arithmetic operations on the Stacks blockchain.

## Features

- **Basic Arithmetic Operations**: Addition, subtraction, multiplication, and division
- **Error Handling**: Proper handling of division by zero with custom error codes
- **State Management**: Tracks the result of the last calculation performed
- **Read-Only Functions**: Efficient querying of the last result without transaction fees

## Contract Functions

### Public Functions (State-Changing)

- `add(a: uint, b: uint)` - Adds two unsigned integers
- `subtract(a: uint, b: uint)` - Subtracts two unsigned integers  
- `multiply(a: uint, b: uint)` - Multiplies two unsigned integers
- `divide(a: uint, b: uint)` - Divides two unsigned integers (with zero-division protection)

### Read-Only Functions

- `get-last-result()` - Returns the result of the last successful calculation

## Error Handling

The contract includes proper error handling for edge cases:

- **Division by Zero**: Returns error code `u100` when attempting to divide by zero
- **State Preservation**: Failed operations do not update the last calculation result

## Getting Started

### Prerequisites

- [Clarinet](https://github.com/hirosystems/clarinet) - Clarity development environment
- [Stacks CLI](https://docs.stacks.co/references/stacks-cli) - For deployment

### Installation

1. Clone or download this project
2. Navigate to the project directory:
   ```bash
   cd clarity-calculator
   ```

3. Install dependencies (if any):
   ```bash
   clarinet install
   ```

### Running Tests

Execute the comprehensive test suite:

```bash
clarinet test
```

The tests cover:
- All basic arithmetic operations
- Division by zero error handling
- State management verification
- Edge cases and boundary values

### Development

Start the local development environment:

```bash
clarinet console
```

This will start a local Stacks node and provide an interactive console for testing.

### Deployment

Deploy to testnet:

```bash
clarinet deploy --testnet
```

Deploy to mainnet:

```bash
clarinet deploy --mainnet
```

## Usage Examples

### Basic Operations

```clarity
;; Addition
(contract-call? .calculator add u5 u3) ;; Returns (ok u8)

;; Subtraction  
(contract-call? .calculator subtract u10 u4) ;; Returns (ok u6)

;; Multiplication
(contract-call? .calculator multiply u7 u2) ;; Returns (ok u14)

;; Division
(contract-call? .calculator divide u15 u3) ;; Returns (ok u5)
```

### Error Handling

```clarity
;; Division by zero
(contract-call? .calculator divide u10 u0) ;; Returns (err u100)
```

### Reading Last Result

```clarity
;; Get the result of the last successful calculation
(contract-call? .calculator get-last-result) ;; Returns (ok <last-result>)
```

## Contract Architecture

### Data Storage

- `last-calculation-result: int` - Stores the result of the last successful calculation

### Error Codes

- `ERR-DIVISION-BY-ZERO: u100` - Error code for division by zero attempts

### Type Handling

The contract uses both `uint` (unsigned integers) for inputs and `int` (signed integers) for storing results to handle negative results from subtraction operations.

## Security Considerations

- **Input Validation**: All inputs are validated as unsigned integers
- **Error Handling**: Division by zero is properly caught and handled
- **State Consistency**: Failed operations do not corrupt contract state
- **Gas Optimization**: Read-only functions are used where possible to minimize costs

## Testing

The project includes comprehensive unit tests covering:

1. **Basic Operations**: All arithmetic functions with various inputs
2. **Error Scenarios**: Division by zero and error propagation
3. **State Management**: Verification that last result is properly updated
4. **Edge Cases**: Zero values, large numbers, and boundary conditions

Run tests with:
```bash
clarinet test
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Resources

- [Clarity Language Documentation](https://docs.stacks.co/references/clarity-language)
- [Stacks Documentation](https://docs.stacks.co/)
- [Clarinet Documentation](https://github.com/hirosystems/clarinet)
- [Stacks Improvement Proposals](https://github.com/stacksgov/sips)

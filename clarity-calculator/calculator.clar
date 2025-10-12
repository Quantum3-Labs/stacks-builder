;; calculator.clar
;; A simple calculator smart contract demonstrating basic arithmetic operations
;; and error handling for division by zero.

;; Error codes
(define-constant ERR-DIVISION-BY-ZERO u100)

;; --- Data Storage ---
;; (define-data-var name type initial-value)
;; Stores the result of the last calculation.
;; We use 'int' (signed integer) for `last-calculation-result`
;; as subtraction can yield negative numbers, which 'uint' (unsigned integer)
;; cannot represent. [DOC 4] explains signed integers.
(define-data-var last-calculation-result int u0)

;; --- Public Functions ---
;; These functions modify the contract state (by updating `last-calculation-result`)
;; and can be called by any principal. They return a response type,
;; either (ok <value>) on success or (err <error-code>) on failure.

;; @desc Performs addition of two unsigned integers.
;; @param a (uint) - The first operand.
;; @param b (uint) - The second operand.
;; @returns (response int uint) - Ok with the sum, or an error.
(define-public (add (a uint) (b uint))
  (let
    ;; Perform addition. [DOC 1] shows (+ u2 u3) for addition.
    ((sum (+ a b)))
    (ok (begin
      ;; Update the last calculation result. We convert the uint sum to int.
      (var-set last-calculation-result (unwrap-panic (as-int sum)))
      sum
    ))
  )
)

;; @desc Performs subtraction of two unsigned integers.
;; @param a (uint) - The first operand.
;; @param b (uint) - The second operand.
;; @returns (response int uint) - Ok with the difference, or an error.
;; Note: The result can be negative. [DOC 1] shows (- 5 10) which implies
;; results can be negative, requiring an 'int' type.
(define-public (subtract (a uint) (b uint))
  (let
    ;; Convert uint inputs to int for subtraction to correctly handle potential negative results.
    ;; unwrap-panic is used here as uint to int conversion for positive numbers is safe.
    ((difference (- (unwrap-panic (as-int a)) (unwrap-panic (as-int b)))))
    (ok (begin
      ;; Update the last calculation result with the signed integer difference.
      (var-set last-calculation-result difference)
      difference
    ))
  )
)

;; @desc Performs multiplication of two unsigned integers.
;; @param a (uint) - The first operand.
;; @param b (uint) - The second operand.
;; @returns (response int uint) - Ok with the product, or an error.
(define-public (multiply (a uint) (b uint))
  (let
    ;; Perform multiplication. [DOC 1] shows (* u2 u16) for multiplication.
    ((product (* a b)))
    (ok (begin
      ;; Update the last calculation result. We convert the uint product to int.
      (var-set last-calculation-result (unwrap-panic (as-int product)))
      product
    ))
  )
)

;; @desc Performs division of two unsigned integers. Includes error handling for division by zero.
;; @param a (uint) - The dividend.
;; @param b (uint) - The divisor.
;; @returns (response int uint) - Ok with the quotient (integer division), or an error.
(define-public (divide (a uint) (b uint))
  ;; Check for division by zero before performing the operation.
  ;; If the divisor 'b' is u0, return an error.
  (if (is-eq b u0)
    (err ERR-DIVISION-BY-ZERO) ;; Return predefined error code.
    (let
      ;; Perform division. [DOC 1] explains that integer division drops decimals (e.g., (/ u10 u3) evaluates to u3).
      ((quotient (/ a b)))
      (ok (begin
        ;; Update the last calculation result. We convert the uint quotient to int.
        (var-set last-calculation-result (unwrap-panic (as-int quotient)))
        quotient
      ))
    )
  )
)

;; --- Read-Only Functions ---
;; These functions do not modify the contract state and can be called without a transaction.

;; @desc Retrieves the result of the last calculation performed by the contract.
;; @returns (response int uint) - Ok with the last stored integer result.
(define-read-only (get-last-result)
  ;; Return the value of `last-calculation-result`.
  (ok (var-get last-calculation-result))
)

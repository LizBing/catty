# Codex latest handoff

**Date:** 2026-07-12  
**Role:** Architecture and integration maintainer  
**Workstream:** R2-A contract draft
**Branch:** `main`  
**Base commit:** `0b82986`
**Head commit:** The integration commit containing this handoff

## Delivered

- Drafted R2-A as a bounded implementation contract for strict unresolved
  natives and an audited native inventory.
- Defined explicit native declaration/resolution states and a return-or-pending-
  exception invocation invariant.
- Scoped catchable `UnsatisfiedLinkError` across interpreter, IR, and tiered
  AOT/runtime paths without pulling broad AOT exception lowering into R2-A.
- Defined classification and generated-inventory rules for registry patches and
  synthetic native methods.
- Authorized no runtime implementation; owner remains unassigned pending review.

## Validation

- Inspected `rtda.Method`, classloader patching, interpreter native invocation,
  runtime/AOT bridge invocation, exception hierarchy, and all current native
  registrations/helper patterns.
- Confirmed the generic stub, unconditional return transfer, missing ULE class,
  function-only registry metadata, and duplicated native branch.
- No runtime code or test behavior changed.

## Assignment

- LizBing accepted R2-A and assigned Claude/DeepSeek as implementation owner;
  Codex is reviewer and integration maintainer.
- Implementation must decide the shared exception-construction helper location
  without introducing a package cycle; the contract constrains behavior rather
  than prescribing the package prematurely.

## Next action

Claude/DeepSeek implements R2-A A–G in an isolated worktree, runs all required
gates, and updates its handoff. Codex then performs an independent review.

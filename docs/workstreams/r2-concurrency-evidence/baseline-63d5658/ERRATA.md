# Baseline evidence erratum

`run-concurrency-results.txt` is the first complete 15-row run. Its summary and
Interpreter/IR/full AOT-runtime output are valid, but the harness retained only the last
five lines of stdout/stderr for AOT `NO-BUILD` cases.

The harness was corrected without rewriting that first run. The additive
`run-concurrency-results-v2.txt` rerun preserves complete AOT build stdout/stderr for the
same 15 fixtures. Review then found that the accepted “blocking points” category needed
separate sleep and join interruption coverage. The fixture set was expanded additively to
17 before Slice A closed. Planning review then required a direct one-slot producer-consumer
fixture so the Roadmap milestone would not be inferred from lower-level wait/notify tests.
That review also separated the opposite VM-liveness obligation: daemon Threads must not
prevent shutdown. `DaemonLiveness` raises the final denominator to 19.
`run-concurrency-results-v5.txt` is the authoritative complete research baseline; the
earlier runs remain preserved as superseded evidence.

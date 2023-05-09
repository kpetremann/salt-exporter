# execution module
salt foo test.true
salt foo test.exception

# state module
salt foo state.single test.succeed_with_changes name="succeed"
salt foo state.single test.fail_with_changes name="fails"

# state
salt foo state.sls test.succeed
salt foo state.sls test.fail

exit 0
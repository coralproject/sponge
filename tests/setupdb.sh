# It setups the database that we are going to uses for tests
# create mysql database coral_test that i can use for the tests (not for benchmarking)

mysql -uroot -e "create database coral_test"
mysql -uroot coral_test < coral_test.sql

# set the env variable for the strategy CONF
export STRATEGY_CONF="../../tests/strategy_test.json"

CGO_ENABLED=0 go build ./cmd/salt-live
docker compose -f ./e2e_test/docker-compose.yaml -f ./e2e_test/docker-compose.demo.yaml up -d

docker compose -f ./e2e_test/docker-compose.yaml -f ./e2e_test/docker-compose.demo.yaml exec -d salt_master sh -c 'sleep 3 && sh /test/exec_commands.sh'
~/go/bin/vhs tui_usage_demo.tape

docker compose -f ./e2e_test/docker-compose.yaml -f ./e2e_test/docker-compose.demo.yaml exec -d salt_master sh -c 'sleep 3 && sh /test/exec_commands.sh'
~/go/bin/vhs tui_overview_demo.tape

docker compose -f ./e2e_test/docker-compose.yaml -f ./e2e_test/docker-compose.demo.yaml down
rm ./salt-live

cmd /C "cd ../../../services/cmd/srv&&set GOOS=linux&&set GOARCH=386&&go build"
cmd /C "cd ..&&docker-compose build"
cmd /C "cd ..&&docker-compose up"
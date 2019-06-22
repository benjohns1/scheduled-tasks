cmd /C "cd ../../services/cmd/srv&&set GOOS=linux&&set GOARCH=386&&go build"
cmd /C "cd ../..&&docker-compose -f docker-compose.stage.yml build"
cmd /C "cd ../..&&docker-compose -f docker-compose.stage.yml up"
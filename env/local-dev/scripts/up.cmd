start /B /D "../../../services/cmd/srv" cmd /C "go build&&srv"
start /B /D "../../../app" cmd /C "npm run dev"
start /B /D "../../../app" cmd /C "npm run cy:open"
start /B /WAIT /D ".." cmd /C "docker-compose up"
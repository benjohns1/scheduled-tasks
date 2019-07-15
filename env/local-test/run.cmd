set ENV_FILEPATH=%CD%/.env
start /WAIT cmd /C "cd scripts&&up"
cmd /C "cd scripts&&down"